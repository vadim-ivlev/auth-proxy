package router

import (
	"auth-proxy/controller"
	"auth-proxy/middleware"
	"encoding/base64"
	"log"
	"net/http"
	"strings"

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

	r.Use(middleware.HeadersMiddleware())

	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/login", login)
	r.GET("/logout", logout)

	r.GET("/ping", controller.PingHandler)

	apps := r.Group("/apps")
	// apps.Use(gin.BasicAuthForRealm(gin.Accounts{"q": "q"}, "myrealm"))
	// apps.Use(gin.BasicAuth(gin.Accounts{}))
	// apps.Use(basicAuth())

	apps.Use(AuthRequired())
	{
		apps.Any("/app1/*proxypath", controller.ReverseProxy("http://localhost:3001", "/apps/app1"))
		apps.Any("/app2/*proxypath", controller.ReverseProxy("http://localhost:3002", "/apps/app2"))
		apps.Any("/onlinebc/*proxypath", controller.ReverseProxy("http://localhost:7700", "/apps/onlinebc"))
		apps.Any("/rg/*proxypath", controller.ReverseProxy("https://rg.ru", "/apps/rg"))
	}

	return r
}

// *****************************************************************************************************
// https://github.com/Depado/gin-auth-example/blob/master/main.go

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			// You'd normally redirect to login page
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
			c.Abort()
		} else {
			// Continue down the chain to handler etc
			c.Next()
		}
	}
}

func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Parameters can't be empty"})
		return
	}
	if username == "hello" && password == "itsme" {
		session.Set("user", username) //In real world usage you'd set this to the users ID
		err := session.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	}
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
	} else {
		log.Println(user)
		session.Delete("user")
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}
}

//***********************************************************************************
// https://www.pandurang-waghulde.com/2018/09/custom-http-basic-authentication-using.html

func basicAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			respondWithError(401, "Unauthorized", c)
			return
			// c.Header("WWW-Authenticate", "abc")
			// c.AbortWithStatus(http.StatusUnauthorized)
			// return

		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !authenticateUser(pair[0], pair[1]) {

			respondWithError(401, "Unauthorized", c)
			return
		}

		c.Next()
	}
}

func authenticateUser(username, password string) bool {
	// var user models.User
	// err := db.Client().Where(models.User{Login: username, Password: password}).FirstOrCreate(&user)
	// if err.Error != nil {
	// 	return false
	// }

	if username == "a" && password == "a" {
		return true
	}

	return false
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
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
