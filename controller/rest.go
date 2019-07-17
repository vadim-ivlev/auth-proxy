package controller

import (
	"auth-proxy/model/auth"

	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	// gsessions "github.com/gorilla/sessions"

	"auth-proxy/primitiveproxy"

	"github.com/gin-gonic/gin"
)

// proxies - Отображение appname -> proxy.
// Перечень прокси серверов предзаготовленных для каждого приложения.
var proxies map[string]*httputil.ReverseProxy

func LandingPage(c *gin.Context) {
	htmlFile, _ := ioutil.ReadFile("./templates/index.html")
	c.Data(http.StatusOK, "text/html; charset=utf-8", htmlFile)
}

func Proxy(c *gin.Context) {
	appname := c.Param("appname")
	proxypath := c.Param("proxypath")
	proxy, ok := proxies[appname]
	// log.Printf(` appname=%v proxypath=%v proxy=%v`, appname, proxypath, proxy)
	if ok {
		c.Request.URL.Path = proxypath
		proxy.ServeHTTP(c.Writer, c.Request)
	} else {
		log.Println("No proxy for appname:", appname)
		c.JSON(200, gin.H{"error": "No proxy for " + appname})
	}
}

// createProxy создает прокси сервер для конкретного URL
func createProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Println("ERR", err)
	}
	return httputil.NewSingleHostReverseProxy(targetURL)
}

// CreateProxies создает глобальный массив proxies в соответствии с таблицей app
func CreateProxies() {
	proxies = make(map[string]*httputil.ReverseProxy)
	appUrls, err := auth.GetAppURLs()
	if err != nil {
		return
	}
	for app, url := range appUrls {
		proxies[app] = createProxy(url)
	}
}

// func getProxy(appname string) *httputil.ReverseProxy {
// 	proxy, ok := proxies[appname]
// 	if !ok {
// 		target, _ := auth.GetAppURL(appname)
// 		proxies[appname] = createProxy(target)
// 		proxy = proxies[appname]
// 	}
// 	return proxy
// }

// ReverseProxy перенаправляет запросы к другому серверу
func ReverseProxy(target string, pathPrefix string) gin.HandlerFunc {
	url, err := url.Parse(target)
	if err != nil {
		log.Println("ERR", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		// p := c.Request.URL.Path
		// pp := strings.TrimPrefix(p, pathPrefix)
		proxypath := c.Param("proxypath")
		// if pp == proxypath {
		// 	fmt.Println("equal:", proxypath)
		// }
		c.Request.URL.Path = proxypath
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func PrimitiveReverseProxy(target string, pathPrefix string) gin.HandlerFunc {
	proxy := primitiveproxy.NewPrimitiveProxy(target)
	return func(c *gin.Context) {
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, pathPrefix)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func ReverseProxy2(scheme string, targ string, pathPrefix string) gin.HandlerFunc {
	target := targ
	return func(c *gin.Context) {
		director := func(req *http.Request) {
			r := c.Request
			// req = r
			req.URL.Scheme = scheme
			req.URL.Host = target
			req.URL.Path = strings.TrimPrefix(req.URL.Path, pathPrefix)
			// req.Header["my-header"] = []string{r.Header.Get("my-header")}
			// Golang camelcases headers
			// delete(req.Header, "My-Header")

			// println("-------------------------")
			// for name, value := range r.Header {
			// 	println(name, value[0])
			// }
			for name, value := range r.Header {
				req.Header.Set(name, value[0])
			}
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
