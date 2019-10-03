package server

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SecureCookie флаг secure на куки браузера
var SecureCookie = false

// Serve запускает сервер на заданном порту. ============================================================
func Serve(port string, tls bool) {
	r := setup()
	if tls {
		_ = r.RunTLS(port, "./certificates/cert.pem", "./certificates/key.pem")
	} else {
		_ = r.Run(port)
	}
}

// Setup определяет пути и присоединяет функции middleware.
func setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r := gin.New()

	CreateProxies()

	r.StaticFile("/favicon.ico", "./templates/favicon.ico")
	r.Static("/templates", "./templates")
	r.LoadHTMLGlob("templates/*.html")

	r.Use(HeadersMiddleware())

	store := sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 86400 * 365 * 5, Secure: SecureCookie}) //0 - for session life
	r.Use(sessions.Sessions("auth-proxy", store))

	r.GET("/testapp", LandingPage)
	r.POST("/graphql", GraphQL)
	r.POST("/schema", GraphQL)

	apps := r.Group("/apps")
	apps.Use(CheckUserMiddleware())
	apps.Any("/:appname/*proxypath", Proxy)
	return r
}
