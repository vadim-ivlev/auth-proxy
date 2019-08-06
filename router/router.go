package router

import (
	"auth-proxy/controller"
	"auth-proxy/middleware"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Serve запускает сервер на заданном порту. ============================================================
func Serve(port string) {
	r := Setup()
	r.Run(port)
}

// Setup определяет пути и присоединяет функции middleware.
func Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r := gin.New()

	controller.CreateProxies()

	r.StaticFile("/favicon.ico", "./templates/favicon.ico")
	r.Static("/templates", "./templates")
	r.LoadHTMLGlob("templates/*.html")

	r.Use(middleware.HeadersMiddleware())

	store := sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 86400 * 365 * 5}) //0 - for session life
	r.Use(sessions.Sessions("auth-proxy", store))

	r.GET("/testapp", controller.LandingPage)
	r.POST("/graphql", controller.GraphQL)
	r.POST("/schema", controller.GraphQL)

	apps := r.Group("/apps")
	apps.Use(middleware.CheckUser())
	apps.Any("/:appname/*proxypath", controller.Proxy)
	return r
}
