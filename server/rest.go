package server

import (
	"auth-proxy/pkg/app"
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/signature"
	"net/http"

	"auth-proxy/pkg/primitiveproxy"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/steambap/captcha"
)

// proxies - Отображение appname -> proxy.
// Перечень прокси серверов предзаготовленных для каждого приложения.
// var proxies map[string]*httputil.ReverseProxy
var proxies map[string]*primitiveproxy.PrimitiveProxy

// Redirects перенаправления браузера для предоставления различных GUI
// var Redirects map[string]string

// Proxy проксирование
func Proxy(c *gin.Context) {
	appname := c.Param("appname")
	proxypath := c.Param("proxypath")
	proxy, ok := proxies[appname]
	log.Printf(` appname=%v proxypath=%v `, appname, proxypath)
	if ok {
		c.Request.URL.Path = proxypath
		proxy.ServeHTTP(c.Writer, c.Request)
	} else {
		c.JSON(200, gin.H{"error": "No proxy url for " + appname})
	}
}

// Captcha source https://github.com/steambap/captcha
func Captcha(c *gin.Context) {
	data, _ := captcha.New(120, 38, func(options *captcha.Options) {
		options.CharPreset = "123456789"
		options.FontScale = 1.3
		// options.TextLength = 5

	})
	// data, _ := captcha.NewMathExpr(100, 38, func(options *captcha.Options) {
	// 	options.FontScale = 1.4
	// })

	SetSessionVariable(c, "captcha", data.Text)
	// log.Println("Captcha text=", data.Text)

	if UseCaptcha {
		// send image data to client
		data.WriteImage(c.Writer)
	} else {
		// send empty image data to client
		c.Data(http.StatusOK, "image/png", []byte{})
	}
}

// createProxy создает прокси сервер для конкретного URL
func createProxy(target, appname, rebase, xtoken string) *primitiveproxy.PrimitiveProxy {
	return primitiveproxy.NewPrimitiveProxy(target, appname, rebase, xtoken)
}

// CreateProxies создает глобальный массив proxies в соответствии с таблицей app
func CreateProxies() {
	proxies = make(map[string]*primitiveproxy.PrimitiveProxy)
	appUrls, err := auth.GetAppURLs()
	if err != nil {
		return
	}
	for app, urlRebase := range appUrls {
		proxies[app] = createProxy(urlRebase[0], app, urlRebase[1], urlRebase[2])
	}
}

// optionHandler По просьбе Леши. Appolo требует этого
func optionHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "origin, content-type, accept, cookie, x-req-id")
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH, OPTIONS")
	c.JSON(http.StatusOK, "")
}

// publicKeyHandler returns text representaton of the public key used for verifying request signatures
func publicKeyHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.String(200, signature.PublicKeyText)
}

// Для тестирования работоспособности
func pingHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, "pong")
}

// listPublicApps - возвращает список публичных приложений не требующих авторизации пользователя.
func listPublicApps(c *gin.Context) {
	c.JSON(http.StatusOK, auth.ListPublicApps())
}

// redirectToAdminURL - направляет браузер на админку
func redirectToAdminURL(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, app.Params.AdminUrl)
}
