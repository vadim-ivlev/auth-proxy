package server

import (
	// "github.com/gin-gonic/contrib/sessions"
	// "github.com/gin-contrib/sessions"
	"auth-proxy/pkg/app"
	"auth-proxy/pkg/prometeo"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	r.Use(prometeo.CountersMiddleware())
	r.Use(CountersMiddleware())
	r.Use(RedirectsMiddleware())
	r.Use(HeadersMiddleware())

	Store = cookie.NewStore([]byte("secret"))
	Store.Options(sessions.Options{
		MaxAge: 86400 * 365 * 5, //0 - for session life
		Secure: SecureCookie,
		Path:   "/",
		// FIXME: Новая версия Chrome требует установки флага SameSite: http.SameSiteNoneMode.
		// что требует установки флага Secure: true,
		// что в свою очередь требует https протокола (tls=true),
		// что не работает в локальной версии без валидного сертификата.
		// Попросить валидный сертификат у админов или поставить envoy as an ssl proxy
		SameSite: http.SameSiteNoneMode,
	})
	r.Use(sessions.Sessions("auth-proxy", Store))

	r.GET("/captcha", Captcha)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.OPTIONS("/graphql", optionHandler)
	r.OPTIONS("/schema", optionHandler)

	r.POST("/graphql", GraphQL)
	r.POST("/schema", GraphQL)

	r.GET("/publickey", publicKeyHandler)
	r.GET("/stat", app.Stat)
	r.GET("/oauthproviders", ListOauthProviders)
	r.GET("/oauthlogin/:provider", OauthLogin)
	r.GET("/oauthlogout/:provider", OauthLogout)
	r.GET("/oauthcallback/:provider", OauthCallback)

	apps := r.Group("/apps")
	apps.Use(CheckUserMiddleware())
	apps.Any("/:appname/*proxypath", Proxy)

	return r
}
