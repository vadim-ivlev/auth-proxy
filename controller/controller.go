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

// PingHandler нужен для фронта, так как сначала отправляется метод с OPTIONS
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func LandingPage(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.HTML(200, "index.html", nil)
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if auth.CheckUserPassword(username, password) {
		session.Set("user", username) //In real world usage you'd set this to the users ID
		session.Options(sessions.Options{MaxAge: 0})

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
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
	} else {
		session.Delete("user")
		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": user.(string) + " successfully logged out"})
	}
}

// func Login1(c *gin.Context) {
// 	r := c.Request
// 	w := c.Writer
// 	session, err := auth.Store.Get(c.Request, "auth-proxy")
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "auth.Store.Get(c.Request, auth-proxy)"})
// 	}

// 	username := c.PostForm("username")
// 	password := c.PostForm("password")

// 	if auth.CheckUserPassword(username, password) {
// 		session.Values["user"] = username
// 		// FIXME: MaxAge=0 does not work
// 		session.Options = &gsessions.Options{
// 			MaxAge: 0,
// 		}

// 		err := session.Save(r, w)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		} else {
// 			c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated " + username})
// 		}
// 	} else {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
// 	}
// }

// func Logout1(c *gin.Context) {
// 	r := c.Request
// 	w := c.Writer
// 	session, err := auth.Store.Get(c.Request, "auth-proxy")
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "auth.Store.Get(c.Request, auth-proxy)"})
// 	}
// 	user, ok := session.Values["user"].(string)

// 	if !ok {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
// 	} else {
// 		delete(session.Values, "user")
// 		session.Options = &gsessions.Options{
// 			MaxAge: -1,
// 		}
// 		session.Save(r, w)
// 		c.JSON(http.StatusOK, gin.H{"message": user + " successfully logged out"})
// 	}
// }

// ReverseProxy перенаправляет запросы к другому серверу
func ReverseProxy(target string, pathPrefix string) gin.HandlerFunc {
	url, err := url.Parse(target)
	if err != nil {
		log.Println("ERR", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		// session := sessions.Default(c)
		// user := session.Get("user")
		// println("USER:", user.(string))

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
