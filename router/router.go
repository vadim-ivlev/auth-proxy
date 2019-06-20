package router

import (
	"auth-proxy/controller"
	"auth-proxy/middleware"

	"github.com/gin-gonic/gin"
)

// Setup определяет пути и присоединяет функции middleware.
func Setup() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r := gin.New()

	// подключаем Middleware
	r.Use(middleware.HeadersMiddleware())

	// defineRoutes
	// r.GET("/ping", PingHandler)

	// TODO: experiment with groups?
	r.Any("/*proxypath", controller.ReverseProxy("http://localhost:3000", "/app1"))
	// r.Any("/app1/*proxypath", controller.ReverseProxy2("http", "localhost:3000", "/app1"))

	return r
}

// Serve запускает сервер на заданном порту.
func Serve(port string) {
	r := Setup()
	r.Run(port)
}
