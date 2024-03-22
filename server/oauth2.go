package server

import (
	"auth-proxy/pkg/app"
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/counter"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type oauth2Provider struct {
	ClientID       string   `yaml:"client_id"`
	ClientSecret   string   `yaml:"client_secret"`
	AuthURI        string   `yaml:"auth_uri"`
	TokenURI       string   `yaml:"token_uri"`
	TokenRevokeURI string   `yaml:"token_revoke_uri"`
	UserInfoURI    string   `yaml:"user_info_uri"`
	EmailFieldName string   `yaml:"email_field_name"`
	NameFieldName  string   `yaml:"name_field_name"`
	RedirectURI    string   `yaml:"redirect_uri"`
	Scopes         []string `yaml:"scopes"`
}

type oauth2Params map[string]oauth2Provider

var (

	// Oauth2Params Oauth2 параметры для различных провайдеров.
	// [имя провайдера] -> { ... oauth2Provider ... }
	Oauth2Params oauth2Params

	// oauthStateLogin a random string
	oauthStateLogin = "login897098and89769087"

	// oauthStateLogout a random string
	oauthStateLogout = "logout543769807-981234"
)

// ReadOauth2Config reads YAML with Oauth2 params
func ReadOauth2Config(fileName string, env string) {
	yamlFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Printf("ERROR ReadOauth2Config(): file=%s, env=%s, error=%s\n", fileName, env, err.Error())
		return
	}

	envParams := make(map[string]oauth2Params)
	err = yaml.Unmarshal(yamlFile, &envParams)
	if err != nil {
		panic(err.Error())
	}
	Oauth2Params = envParams[env]
}

// *****************************************************

func buildOauthConfig(provider string) *oauth2.Config {
	params := Oauth2Params[provider]
	oauthConfig := &oauth2.Config{
		ClientID:     params.ClientID,
		ClientSecret: params.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  params.AuthURI,
			TokenURL: params.TokenURI,
		},
		Scopes:      params.Scopes,
		RedirectURL: app.Params.AdminAPI + params.RedirectURI,
	}
	return oauthConfig
}

// ListOauthProviders перечисляет Oauth2 login URIs для каждого провайдера сервиса аутентификации
func ListOauthProviders(c *gin.Context) {
	providerURLs := make(map[string][]string)
	for provider := range Oauth2Params {
		providerURLs[provider] = []string{"/oauthlogin/" + provider, "/oauthlogout/" + provider}
	}
	c.JSON(200, providerURLs)
}

// OauthLogin начинает процесс аутентификации для данного провайдера
func OauthLogin(c *gin.Context) {
	provider := c.Param("provider")
	OauthLoginProvider(c, provider)
}

func OauthLoginProvider(c *gin.Context, provider string) {
	oauthConfig := buildOauthConfig(provider)
	url := oauthConfig.AuthCodeURL(oauthStateLogin)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// OauthLogout начинает  Loging Out для данного провайдера
func OauthLogout(c *gin.Context) {
	provider := c.Param("provider")
	OauthLoginProvider(c, provider)
}

func OauthLogoutProvider(c *gin.Context, provider string) {
	oauthConfig := buildOauthConfig(provider)
	url := oauthConfig.AuthCodeURL(oauthStateLogout)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// OauthCallback заканчивает процесс входа/выхода для данного провайдера.
// Перенаправляет браузер на auth-admin в параметрах запроса передавая информацию
// о пользователе и ошибке
func OauthCallback(c *gin.Context) {
	provider := c.Param("provider")
	state := c.Query("state")
	code := c.Query("code")

	// если запрос был LogIn
	if state == oauthStateLogin {
		oauth2Email, oauth2Name, err := getOauth2UserInfo(provider, code)
		if err != nil {
			log.Println(err)
			redirectWithMessage(c, err.Error(), oauth2Email, oauth2Name)
			return
		}
		oauth2error := loginOauth2(c, provider, oauth2Email)
		redirectWithMessage(c, oauth2error, oauth2Email, oauth2Name)
		return
	}

	// если запрос был LogOut
	if state == oauthStateLogout {
		_ = logoutOauth2(provider, code)
		redirectWithMessage(c, "", "", "")
		return
	}

	// Непонятно чего хотел клиент
	redirectWithMessage(c, "OauthCallback: unknown state string", "", "")
}

// loginOauth2 пытается залогинить пользователя по паролю полученному из социальной сети
func loginOauth2(c *gin.Context, provider, oauth2Email string) (oauth2error string) {

	// ищем пользователя с таким email в базе данных
	username := auth.GetUserNameByEmail(oauth2Email)

	// если нет пользователя возвращаем ошибку
	if username == "" {
		return "Вход в <b>" + provider + "</b> осуществлен, но <b>" + oauth2Email + "</b> не зарегистрирован в auth-proxy."
	}

	// если пользователь отключен возвращаем ошибку
	if !auth.IsUserEnabled(username) {
		return "Вход в <b>" + provider + "</b> осуществлен, но пользователь <b>" + username + " / " + oauth2Email + "</b> заблокирован."
	}
	// сбрасываем счетчик неудачных попыток
	counter.ResetCounter(username)

	// SUCCESS. Аутентифицируем пользователя.
	err := SetSessionVariable(c, "user", username)
	if err != nil {
		return "Не удалось сохранить сессию " + username
	}

	return ""
}

// getOauth2UserInfo здесь выполняется вся грязная работа по извлечению email
// и имени пользователя из API конкретного провайдера.
// Принимает имя провайдера и код аторизации.
func getOauth2UserInfo(provider string, code string) (email string, name string, err error) {

	oauthConfig := buildOauthConfig(provider)

	// обмениваем код авторизации на токен доступа
	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return email, name, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	params := Oauth2Params[provider]

	// получаем информацию о пользователе из API провайдеоа
	response, err := http.Get(params.UserInfoURI + token.AccessToken)
	if err != nil {
		return email, name, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return email, name, fmt.Errorf("failed reading "+provider+" response body: %s", err.Error())
	}

	fmt.Printf("content=%s\n\n", contents)

	var userInfo map[string]interface{}

	// a fix for github
	if provider == "github" {
		userInfo = getGithubUserInfo(contents)
	} else {
		userInfo = jsonStringToMap(string(contents))
	}

	if userInfo == nil {
		return email, name, fmt.Errorf("can not get user info")
	}

	// fmt.Println("****************************************************************")
	// fmt.Println(userInfo)

	emailI, ok := userInfo[params.EmailFieldName]
	if ok {
		email, _ = emailI.(string)
	}
	nameI, ok := userInfo[params.NameFieldName]
	if ok {
		name, _ = nameI.(string)
	}

	if provider == "vk" {
		email, name = getVkUserEmailAndName(token, params, userInfo)
	}

	println("email, name = ", email, name)

	if email == "" {
		return email, name, fmt.Errorf("ошибка. email отсутствует ")
	}

	return email, name, nil
}

// getGithubUserInfo a dirty patch for github
func getGithubUserInfo(contents []byte) map[string]interface{} {
	a := make([]interface{}, 0)
	err := json.Unmarshal(contents, &a)
	if err != nil {
		return nil
	}
	if len(a) == 0 {
		return nil
	}
	ghUserInfo, ok := a[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return ghUserInfo
}

// getVkUserEmailAndName  extacts email and user name from VK.com API
func getVkUserEmailAndName(token *oauth2.Token, params oauth2Provider, userInfo map[string]interface{}) (email string, name string) {
	tokenEmail := token.Extra("email")
	if tokenEmail != nil {
		email = tokenEmail.(string)
	}

	respI, ok := userInfo[params.NameFieldName]
	if ok {
		resp, ok := respI.([]interface{})
		if ok && len(resp) > 0 {
			resp0 := resp[0]
			r, ok := resp0.(map[string]interface{})
			if ok {
				fname, _ := r["first_name"].(string)
				lname, _ := r["last_name"].(string)
				name = fname + " " + lname
			}
		}
	}
	return email, name
}

// logoutOauth2 revokes access token
func logoutOauth2(provider string, code string) error {

	oauthConfig := buildOauthConfig(provider)

	// обмениваем код авторизации на токен доступа
	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return fmt.Errorf("logoutOauth2: code exchange failed: %s", err.Error())
	}

	params := Oauth2Params[provider]

	if provider == "yandex" {
		revokeYandexToken(params.TokenRevokeURI, params.ClientID, params.ClientSecret, token.AccessToken)
		return nil
	}

	if provider == "facebook" {
		revokeFacebookToken(params.TokenRevokeURI, params.ClientID, params.ClientSecret, token.AccessToken)
		return nil
	}

	// revoke the token
	response, err := http.Get(params.TokenRevokeURI + token.AccessToken)
	if err != nil {
		return fmt.Errorf("logoutOauth2: failed revoking token: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("logoutOauth2: failed reading "+provider+" response body: %s", err.Error())
	}

	fmt.Printf("logoutOauth2: content=%s\n\n", contents)

	return nil
}

// revokeYandexToken FIXME:
func revokeYandexToken(urlStr string, clientID string, clientSecret string, accessToken string) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("access_token", accessToken)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	fmt.Println(resp, err)
}

func revokeFacebookToken(urlStr string, clientID string, clientSecret string, accessToken string) {

	client := &http.Client{}
	r, _ := http.NewRequest("DELETE", urlStr+accessToken, nil) // URL-encoded payload

	resp, err := client.Do(r)
	fmt.Println(resp, err)
}

// redirectWithMessage направляет браузер пользователя добавляя сообщение
// и информацию о пользователе если они не пустые.
func redirectWithMessage(c *gin.Context, oauth2error string, oauth2email string, oauth2name string) {
	// u := app.Params.Redirects["/admin"]
	u := app.Params.AdminUrl

	if oauth2error != "" {
		u += "&oauth2error=" + url.QueryEscape(oauth2error)
	}
	if oauth2email != "" {
		u += "&oauth2email=" + url.QueryEscape(oauth2email)
	}
	if oauth2name != "" {
		u += "&oauth2name=" + url.QueryEscape(oauth2name)
	}
	c.Redirect(http.StatusTemporaryRedirect, u)
	c.Abort()
}
