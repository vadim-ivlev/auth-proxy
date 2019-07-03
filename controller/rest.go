package controller

import (
	"auth-proxy/model/auth"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	// gsessions "github.com/gorilla/sessions"

	"auth-proxy/primitiveproxy"

	"github.com/gin-gonic/gin"
)

func LandingPage(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	username := auth.GetUserName(c)
	data := map[string]interface{}{
		"username": username,
		"info":     auth.GetUserInfo(username),
	}
	c.HTML(200, "index.html", data)
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if auth.CheckUserPassword(username, password) {
		session.Set("user", username)
		// session.Options(sessions.Options{MaxAge: 0})

		err := session.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated " + username})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	}
}

func Logout(c *gin.Context) {
	username := auth.GetUserName(c)
	auth.DeleteSession(c)
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
