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

func getOauthConfig(provider string) *oauth2.Config {
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
	oauthConfig := getOauthConfig(provider)
	url := oauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
	// c.Abort()
}

// OauthCallback заканчивает процесс аутентификации для данного провайдера.
// В случае успеха Oauth2 аутентификации
// аутентифицирует пользователя в auth-proxy,
// при условии что пользователь с Oauth2 email,
// зарегистрирован в auth-proxy.
func OauthCallback(c *gin.Context) {
	provider := c.Param("provider")

	// FIXME: use gin to get field values
	r := c.Request
	oauth2Email, err := getOauth2UserEmail(provider, r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, app.Params.Redirects["/admin"]+"&oauth2error="+url.QueryEscape(err.Error()))
		return
	}

	fmt.Printf("EMAIL = %s \n", oauth2Email)

	// ищем пользователя с таким email в базе данных
	username := auth.GetUserNameByEmail(oauth2Email)

	// если нет пользователя возвращаем ошибку
	if username == "" {
		redirectWithMessage(c, "Извините, <b>"+oauth2Email+"</b> не зарегистрирован.<br>Зарегистрируйтесь, или обратитесь к администратору.")
		return
	}

	// если пользователь отключен возвращаем ошибку
	if !auth.IsUserEnabled(username) {
		redirectWithMessage(c, "Извините, "+username+" / "+oauth2Email+" заблокирован.")
		return
	}

	// сбрасываем счетчик неудачных попыток
	counter.ResetCounter(username)

	// Устанавливаем куки
	err = SetSessionVariable(c, "user", username)
	if err != nil {
		redirectWithMessage(c, "Не удалось сохранить сессию")
		return
	}
	// Успешно завершаем
	// c.JSON(200, "ok")
	// redirectWithMessage(c, "Success. "+username+" is authenticated")
	c.Redirect(http.StatusTemporaryRedirect, app.Params.Redirects["/admin"])

}

// getOauth2UserEmail здесь выполняется вся грязная работа по извлечению email
// из API конкретного провайдера.
func getOauth2UserEmail(provider string, state string, code string) (string, error) {
	if state != oauthStateString {
		return "", fmt.Errorf("invalid oauth state")
	}

	oauthConfig := getOauthConfig(provider)

	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", fmt.Errorf("code exchange failed: %s", err.Error())
	}
	params := Oauth2Params[provider]
	response, err := http.Get(params.UserInfoURI + token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf("failed reading response body: %s", err.Error())
	}

	fmt.Printf("content=%s\n\n", contents)

	var userInfo map[string]interface{}

	// FIXME: a dirty fix fo github
	if provider == "github" {
		userInfo = getGithubUserInfo(contents)
		if userInfo == nil {
			return "", fmt.Errorf("can not get GitHub user info")
		}
	} else {
		userInfo = jsonStringToMap(string(contents))
	}

	email, ok := userInfo[params.EmailFieldName]
	if !ok {
		return "", fmt.Errorf("There is no '" + params.EmailFieldName + "' field in user info")
	}

	return email.(string), nil

	// return contents, nil
}

// getGithubUserInfo a dirty patch for github
func getGithubUserInfo(contents []byte) map[string]interface{} {
	a := make([]interface{}, 0)
	_ = json.Unmarshal(contents, &a)
	ghUserInfo, ok := a[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return ghUserInfo
}

func redirectWithMessage(c *gin.Context, message string) {
	u := app.Params.Redirects["/admin"] + "&oauth2error=" + url.QueryEscape(message)
	fmt.Println("Redirect url=", u)
	c.Redirect(http.StatusTemporaryRedirect, u)
	c.Abort()
}
