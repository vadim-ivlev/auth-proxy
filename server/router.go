package server

import (
	// "github.com/gin-gonic/contrib/sessions"
	// "github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

// SecureCookie флаг secure на куки браузера
var SecureCookie = false
var Store cookie.Store

// var StoreCapt cookie.Store

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

	r.Use(CountersMiddleware())
	r.Use(RedirectsMiddleware())
	r.Use(HeadersMiddleware())

	Store = cookie.NewStore([]byte("secret"))
	Store.Options(sessions.Options{MaxAge: 86400 * 365 * 5, Secure: SecureCookie, Path: "/"}) //0 - for session life
	r.Use(sessions.Sessions("auth-proxy", Store))

	r.GET("/captcha", Captcha)
	r.POST("/graphql", GraphQL)
	r.POST("/schema", GraphQL)

	r.GET("/oauthproviders", ListOauthProviders)
	r.GET("/oauthlogin/:provider", OauthLogin)
	r.GET("/oauthlogout/:provider", OauthLogout)
	r.GET("/oauthcallback/:provider", OauthCallback)

	apps := r.Group("/apps")
	apps.Use(CheckUserMiddleware())
	apps.Any("/:appname/*proxypath", Proxy)

	return r
}
