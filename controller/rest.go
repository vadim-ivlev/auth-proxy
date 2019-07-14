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

func LandingPage(c *gin.Context) {
	// c.HTML(200, "index.html", map[string]interface{}{})

	htmlFile, _ := ioutil.ReadFile("./templates/index.html")
	c.Data(http.StatusOK, "text/html; charset=utf-8", htmlFile)
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	err := auth.Login(c, username, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()}) // http.StatusUnauthorized
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated " + username})
	}
}

func Logout(c *gin.Context) {
	username := auth.GetUserName(c)
	auth.Logout(c)
	c.JSON(http.StatusOK, gin.H{"message": username + " successfully logged out"})
}

// ReverseProxy перенаправляет запросы к другому серверу
func ReverseProxy(target string, pathPrefix string) gin.HandlerFunc {
	url, err := url.Parse(target)
	if err != nil {
		log.Println("ERR", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, pathPrefix)
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
