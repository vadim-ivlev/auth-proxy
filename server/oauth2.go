package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type oauth2Provider struct {
	ClientID     string   `yaml:"client_id"`
	AuthURI      string   `yaml:"auth_uri"`
	TokenURI     string   `yaml:"token_uri"`
	ClientSecret string   `yaml:"client_secret"`
	RedirectURI  string   `yaml:"redirect_uri"`
	Scopes       []string `yaml:"scopes"`
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

	endpoint := oauth2.Endpoint{AuthURL: params.AuthURI, TokenURL: params.TokenURI}
	oauthConfig := &oauth2.Config{
		RedirectURL:  params.RedirectURI,
		ClientID:     params.ClientID,
		ClientSecret: params.ClientSecret,
		Scopes:       params.Scopes,
		// Endpoint: google.Endpoint,
		Endpoint: endpoint,
	}
	return oauthConfig
}

func ListOauthProviders(c *gin.Context) {
	providerURLs := make(map[string]string)
	for provider, _ := range Oauth2Params {
		providerURLs[provider] = "/oauthlogin/" + provider
	}
	c.JSON(200, providerURLs)
}

func OauthLogin(c *gin.Context) {
	provider := c.Param("provider")
	oauthConfig := getOauthConfig(provider)
	url := oauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
	// c.Abort()
}

func OauthCallback(c *gin.Context) {
	provider := c.Param("provider")

	r := c.Request
	content, err := getUserInfo(provider, r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(c.Writer, "Content: %s\n", content)

	// c.JSON(200, Oauth2Params[provider])
}

func getUserInfo(provider string, state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	oauthConfig := getOauthConfig(provider)

	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	str := string(contents)
	fmt.Println(str)

	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}
