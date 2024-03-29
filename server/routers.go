package server

import (
	// "github.com/gin-gonic/contrib/sessions"
	// "github.com/gin-contrib/sessions"
	"auth-proxy/pkg/app"
	"auth-proxy/pkg/authenticator"
	"auth-proxy/pkg/db"
	"auth-proxy/pkg/prometeo"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
)

// SecureCookie флаг secure на куки браузера
var SecureCookie = false

// Store куки
var Store cookie.Store

// Up запускает сервер на заданном порту
func Up(port string, tls bool) {
	r := setup()
	if tls {
		err := r.RunTLS(":"+port, "./certificates/cert.pem", "./certificates/key.pem")
		if err != nil {
			fmt.Println("Error RunTLS", err)
			os.Exit(1)
		}
	} else {
		err := r.Run(":" + port)
		if err != nil {
			fmt.Println("Error server Run", err)
			os.Exit(1)
		}
	}
}

// Setup определяет пути и присоединяет функции middleware.
func setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r := gin.New()

	// CreateProxies()
	// непрерывно обновляем список проксируемых приложений,
	// на случай если БД была изменена извне.
	go keepCreatingProxies()

	r.StaticFile("/favicon.ico", "./templates/favicon.ico")
	r.Static("/admin/", "./admin")

	r.Use(prometeo.CountersMiddleware())
	r.Use(CountersMiddleware())
	// r.Use(RedirectsMiddleware())
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
	r.Use(sessions.Sessions(app.Params.CookieName, Store))

	// r.GET("/admin", redirectToAdminURL)
	// r.GET("/admin", redirectToAdminURL) // - удален 10.03.23 по соображениям безопасности
	// r.GET("/graphql_test", redirectToGraphqlTestURL) // - удален 10.03.23 по соображениям безопасности

	r.GET("/captcha", Captcha)
	r.GET("/metrics", gin.WrapH(promhttp.Handler())) // - удален 10.03.23 по соображениям безопасности

	r.OPTIONS("/graphql", optionHandler)
	r.OPTIONS("/schema", optionHandler)

	r.POST("/graphql", graphqlResult)
	r.POST("/schema", graphqlResult)

	r.GET("/publickey", publicKeyHandler)
	r.GET("/stat", app.Stat)
	r.GET("/logmessage/:message", app.LogMessage)
	r.GET("/publicapps", listPublicApps)
	r.GET("/oauthproviders", ListOauthProviders)
	r.GET("/oauthlogin/:provider", OauthLogin)
	r.GET("/oauthlogout/:provider", OauthLogout)
	r.GET("/oauthcallback/:provider", OauthCallback)

	r.GET("/authenticator_barcode", authenticator.AuthenticatorBarcode)
	r.GET("/authenticator_manualcode", authenticator.AuthenticatorManualCode)
	r.GET("/set_authenticator", authenticator.SetAuthenticator)
	r.GET("/reset_authenticator", authenticator.ResetAuthenticator)

	r.GET("/set_password", authenticator.SetPassword)
	r.GET("/reset_password", authenticator.ResetPassword)
	r.GET("/confirm-email", authenticator.ConfirmEmail)

	r.GET("/ping", pingHandler)

	apps := r.Group("/apps")
	apps.Use(CheckUserMiddleware())
	apps.Use(FlusherMiddleware())
	apps.Any("/:appname/*proxypath", Proxy)

	return r
}

// keepCreatingProxies перезачитывает список проксируемых приложений
// из базы данных каждые 5 минут. Запускается в отдельной горутине.
// Сделано для того чтобы поддерживать приложение в актуальном состоянии
// в случае если база данных была изменена внешним образом.
func keepCreatingProxies() {
	i := 0
	for {
		CreateProxies()
		i++
		fmt.Printf("*** List of proxied apps was refreshed %v times. Refresh interval is %v minutes. ***\n", i, db.Params.Refreshtime)
		time.Sleep(time.Duration(db.Params.Refreshtime) * time.Minute)
	}
}
