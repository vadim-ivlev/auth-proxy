package middleware

import (
	"auth-proxy/model/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// HeadersMiddleware добавляет HTTP заголовки к ответу сервера
func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Header("Access-Control-Allow-Origin", "https://editor.rg.ru")
		// if hostIsAllowed(c.Request.Host) {
		// c.Header("Access-Control-Allow-Origin", "*")
		// }
		// c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, HEAD")
		// c.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		// // c.Header("Content-Type", "application/json; charset=utf-8")
		// c.Header("Access-Control-Max-Age", "600")
		c.Next()
	}
}

// hostIsAllowed - проверяем пришел ли запрос с разрешенного хоста
func hostIsAllowed(host string) bool {
	if strings.HasPrefix(host, "localhost") ||
		strings.HasPrefix(host, "127.0.0.1") ||
		strings.Contains(host, ".rg.ru:") ||
		strings.HasSuffix(host, ".rg.ru") {
		return true
	}
	return false
}

// CheckUser проверяет залогинен ли пользователь
func CheckUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Please login: /login "})
		} else {
			userName := user.(string)
			roles := auth.GetUserRoles(userName, strings.TrimSuffix(strings.TrimPrefix(c.Request.URL.Path, "/apps/"), "/"))
			info := auth.GetUserInfo(userName)

			c.Request.Header.Set("user-roles", roles)
			c.Request.Header.Set("user-info", info)
			c.Next()
		}
	}
}

// // CheckUser проверяет залогинен ли пользователь
// func CheckUser1() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		session, err := auth.Store.Get(c.Request, "auth-proxy")
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "auth.Store.Get(c.Request, auth-proxy)"})
// 		} else {
// 			userName, ok := session.Values["user"].(string)
// 			if !ok {
// 				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Please login: /login "})
// 			}
// 			roles := auth.GetUserRoles(userName, c.Request.URL.Path)
// 			info := auth.GetUserInfo(userName)

// 			c.Request.Header.Set("user-roles", roles)
// 			c.Request.Header.Set("user-info", info)
// 			c.Next()
// 		}
// 	}
// }
