package controller

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// PingHandler нужен для фронта, так как сначала отправляется метод с OPTIONS
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

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

func ReverseProxy2(scheme string, targ string, pathPrefix string) gin.HandlerFunc {

	target := targ

	return func(c *gin.Context) {
		director := func(req *http.Request) {
			// r := c.Request
			// req = r
			req.URL.Scheme = scheme
			req.URL.Host = target
			req.URL.Path = strings.TrimPrefix(req.URL.Path, pathPrefix)
			// req.Header["my-header"] = []string{r.Header.Get("my-header")}
			// Golang camelcases headers
			// delete(req.Header, "My-Header")
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
