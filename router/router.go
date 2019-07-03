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

	r.StaticFile("/favicon.ico", "./templates/favicon.ico")
	r.Static("/templates", "./templates")
	r.LoadHTMLGlob("templates/*.*")

	r.Use(middleware.HeadersMiddleware())

	store := sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 0})
	r.Use(sessions.Sessions("auth-proxy", store))

	r.GET("/", controller.LandingPage)
	r.POST("/login", controller.Login)
	r.GET("/logout", controller.Logout)

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
