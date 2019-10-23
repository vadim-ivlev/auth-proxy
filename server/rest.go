package server

import (
	"auth-proxy/pkg/auth"
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
var AdminUrl = "https://auth-admin.now.sh"

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

func ProxyAdmin(c *gin.Context) {
	// строим строку запроса для того, чтобы
	// auth-admin обращался с запросами к этому серверу.
	host := c.Request.Host
	tls := c.Request.TLS
	prefix := "http://"
	if tls != nil {
		prefix = "https://"
	}

	log.Println(prefix, host)

	c.Redirect(http.StatusMovedPermanently, AdminUrl+"?url="+prefix+host)
	c.Abort()
}

// Captcha
// source https://github.com/steambap/captcha
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
	log.Println("Captcha text=", data.Text)
	// send image data to client
	data.WriteImage(c.Writer)
}

// createProxy создает прокси сервер для конкретного URL
func createProxy(target, appname, rebase string) *primitiveproxy.PrimitiveProxy {
	return primitiveproxy.NewPrimitiveProxy(target, appname, rebase)
}

// CreateProxies создает глобальный массив proxies в соответствии с таблицей app
func CreateProxies() {
	proxies = make(map[string]*primitiveproxy.PrimitiveProxy)
	appUrls, err := auth.GetAppURLs()
	if err != nil {
		return
	}
	for app, url_rebase := range appUrls {
		proxies[app] = createProxy(url_rebase[0], app, url_rebase[1])
	}
}
