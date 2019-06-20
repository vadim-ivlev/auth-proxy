package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// r.GET("/ping", pong)

	// TODO: experiment with groups?
	r.Any("/*proxypath", ReverseProxy("http://localhost:3000", "/app1"))
	// r.Any("/app1/*proxypath", ReverseProxy2("http", "localhost:3000", "/app1"))

	r.Run(":4000") // listen and serve on 0.0.0.0:4000
}

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
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
