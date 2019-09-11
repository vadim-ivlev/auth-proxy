package middleware

import (
	"auth-proxy/model/auth"
	"auth-proxy/model/session"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HeadersMiddleware добавляет HTTP заголовки к ответу сервера
func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		// fmt.Println("origin=", origin)
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// Пока что кроссдоменность действует только для /graphql этого приложения.
		// if c.Request.URL.Path == "/graphql" {
		// 	c.Header("Access-Control-Allow-Origin", "*")
		// }

		c.Next()
	}
}

// CheckUser проверяет залогинен ли пользователь
func CheckUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Если это префлайт запрос, браузеры могут не посылать куки,
		// и пользователь  не определится.
		// Пропускаем этот запрос без изменений
		if c.Request.Method == "OPTIONS" {
			fmt.Println("c.Request.Method=OPTIONS")
			c.Next()
			return
		}

		userName := session.GetUserName(c)
		if userName == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Please login: /login "})
		} else {
			if auth.IsUserEnabled(userName) == false {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Sorry. " + userName + " is disabled."})
			}
			roles := auth.GetUserRoles(userName, strings.TrimSuffix(strings.TrimPrefix(c.Request.URL.Path, "/apps/"), "/"))
			info := auth.GetUserInfo(userName)

			c.Request.Header.Set("user-roles", roles)
			c.Request.Header.Set("user-info", info)
			c.Next()
		}
	}
}
