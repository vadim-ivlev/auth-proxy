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

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type oauth2Provider struct {
	ClientID       string   `yaml:"client_id"`
	ClientSecret   string   `yaml:"client_secret"`
	AuthURI        string   `yaml:"auth_uri"`
	TokenURI       string   `yaml:"token_uri"`
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

	// oauthStateString a random string
	oauthStateString = "slon897098and89769087moska"
)

// ReadOauth2Config reads YAML with Oauth2 params
func ReadOauth2Config(fileName string, env string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err.Error())
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
		RedirectURL: params.RedirectURI,
	}
	return oauthConfig
}

// ListOauthProviders перечисляет Oauth2 login URIs для каждого провайдера сервиса аутентификации
func ListOauthProviders(c *gin.Context) {
	providerURLs := make(map[string]string)
	for provider, _ := range Oauth2Params {
		providerURLs[provider] = "/oauthlogin/" + provider
	}
	c.JSON(200, providerURLs)
}

// OauthLogin начинает процесс аутентификации для данного провайдера
func OauthLogin(c *gin.Context) {
	provider := c.Param("provider")
	oauthConfig := buildOauthConfig(provider)
	url := oauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// OauthCallback заканчивает процесс аутентификации для данного провайдера.
// В случае успеха Oauth2 аутентификации
// аутентифицирует пользователя в auth-proxy,
// при условии что пользователь с Oauth2 email,
// зарегистрирован в auth-proxy.
// Перенаправляет браузер на auth-admin.now.sh в параметрах запроса передавая информацию о пользователе,
// или сообщение об ошибке.
func OauthCallback(c *gin.Context) {
	provider := c.Param("provider")

	// FIXME: use gin to get field values
	r := c.Request
	oauth2Email, oauth2Name, err := getOauth2UserInfo(provider, r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		log.Println(err)
		// c.Redirect(http.StatusTemporaryRedirect, app.Params.Redirects["/admin"]+"&oauth2error="+url.QueryEscape(err.Error()))
		redirectWithMessage(c, err.Error(), oauth2Email, oauth2Name)
		return
	}

	// fmt.Printf("EMAIL = %s \n", oauth2Email)

	// ищем пользователя с таким email в базе данных
	username := auth.GetUserNameByEmail(oauth2Email)

	// если нет пользователя возвращаем ошибку
	if username == "" {
		redirectWithMessage(c, "Извините, <b>"+oauth2Email+"</b> не зарегистрирован.", oauth2Email, oauth2Name)
		return
	}

	// если пользователь отключен возвращаем ошибку
	if !auth.IsUserEnabled(username) {
		redirectWithMessage(c, "Извините, "+username+" / "+oauth2Email+" заблокирован.", oauth2Email, oauth2Name)
		return
	}

	// сбрасываем счетчик неудачных попыток
	counter.ResetCounter(username)

	// SUCCESS. Аутентифицируем пользователя.
	err = SetSessionVariable(c, "user", username)
	if err != nil {
		redirectWithMessage(c, "Не удалось сохранить сессию "+username, oauth2Email, oauth2Name)
		return
	}

	redirectWithMessage(c, "", oauth2Email, oauth2Name)

}

// getOauth2UserInfo здесь выполняется вся грязная работа по извлечению email
// и имени пользователя из API конкретного провайдера.
// Принимает имя провайдера, строку состояния и код аторизации.
func getOauth2UserInfo(provider string, state string, code string) (email string, name string, err error) {
	if state != oauthStateString {
		return email, name, fmt.Errorf("invalid oauth state")
	}

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

	fmt.Println("****************************************************************")
	fmt.Println(userInfo)

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

// redirectWithMessage направляет браузер пользователя добавляя сообщение
// и информацию о пользователе если они не пустые.
func redirectWithMessage(c *gin.Context, oauth2error string, oauth2email string, oauth2name string) {
	u := app.Params.Redirects["/admin"]
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
