package controller

import (
	"auth-proxy/model/auth"

	"auth-proxy/model/primitiveproxy"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// proxies - Отображение appname -> proxy.
// Перечень прокси серверов предзаготовленных для каждого приложения.
// var proxies map[string]*httputil.ReverseProxy
var proxies map[string]*primitiveproxy.PrimitiveProxy

func LandingPage(c *gin.Context) {
	htmlFile, _ := ioutil.ReadFile("./templates/index.html")
	c.Data(http.StatusOK, "text/html; charset=utf-8", htmlFile)
}

func Proxy(c *gin.Context) {
	appname := c.Param("appname")
	proxypath := c.Param("proxypath")
	proxy, ok := proxies[appname]
	log.Printf(` appname=%v proxypath=%v `, appname, proxypath)
	if ok {
		c.Request.URL.Path = proxypath
		proxy.ServeHTTP(c.Writer, c.Request)
	} else {
		log.Println("No proxy for appname:", appname)
		c.JSON(200, gin.H{"error": "No proxy for " + appname})
	}
}

// createProxy создает прокси сервер для конкретного URL
func createProxy(target, appname, rebase string) *primitiveproxy.PrimitiveProxy {
	return primitiveproxy.NewPrimitiveProxy(target, appname, rebase)
}

// CreateProxies создает глобальный массив proxies в соответствии с таблицей app
func CreateProxies() {
	// proxies = make(map[string]*httputil.ReverseProxy)
	proxies = make(map[string]*primitiveproxy.PrimitiveProxy)
	appUrls, err := auth.GetAppURLs()
	if err != nil {
		return
	}
	for app, url_rebase := range appUrls {
		proxies[app] = createProxy(url_rebase[0], app, url_rebase[1])
	}
}

// // createProxy создает прокси сервер для конкретного URL
// func createProxy(target string) *httputil.ReverseProxy {
// 	targetURL, err := url.Parse(target)
// 	if err != nil {
// 		log.Println("ERR", err)
// 	}
// 	return httputil.NewSingleHostReverseProxy(targetURL)
// }
