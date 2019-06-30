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
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r := gin.New()

	r.StaticFile("/favicon.ico", "./templates/favicon.ico")
	r.Static("/templates", "./templates")
	r.LoadHTMLGlob("templates/*.*")

	r.Use(middleware.HeadersMiddleware())
	store := sessions.NewCookieStore([]byte("secret"))
	// store.Options(sessions.Options{MaxAge: 0})
	r.Use(sessions.Sessions("auth-proxy", store))

	r.GET("/", controller.LandingPage)
	r.POST("/login", controller.Login)
	r.GET("/logout", controller.Logout)

	r.GET("/ping", controller.PingHandler)

	apps := r.Group("/apps")
	apps.Use(middleware.CheckUser())
	{
		apps.Any("/app1/*proxypath", controller.ReverseProxy("http://localhost:3001", "/apps/app1"))
		apps.Any("/app2/*proxypath", controller.ReverseProxy("http://localhost:3002", "/apps/app2"))
		apps.Any("/onlinebc/*proxypath", controller.ReverseProxy("http://localhost:7700", "/apps/onlinebc"))
		apps.Any("/rg/*proxypath", controller.ReverseProxy("https://rg.ru", "/apps/rg"))
	}

	return r
}

// **************************GIN ****************************************************************
// BasicAuthForRealm returns a Basic HTTP Authorization middleware. It takes as arguments a map[string]string where
// the key is the user name and the value is the password, as well as the name of the Realm.
// If the realm is empty, "Authorization Required" will be used by default.
// (see http://tools.ietf.org/html/rfc2617#section-1.2)
// func BasicAuthForRealm(accounts gin.Accounts, realm string) gin.HandlerFunc {
// 	if realm == "" {
// 		realm = "Authorization Required"
// 	}
// 	realm = "Basic realm=" + strconv.Quote(realm)
// 	// pairs := processAccounts(accounts)
// 	return func(c *gin.Context) {
// 		// Search user in the slice of allowed credentials
// 		user, found := pairs.searchCredential(c.Request.Header.Get("Authorization"))
// 		if !found {
// 			// Credentials doesn't match, we return 401 and abort handlers chain.
// 			c.Header("WWW-Authenticate", realm)
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 			return
// 		}

// 		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
// 		// c.MustGet(gin.AuthUserKey).
// 		c.Set(gin.AuthUserKey, user)
// 	}
// }
