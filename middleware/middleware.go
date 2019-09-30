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
		// Кроссдоменность действует только для /graphql и /schema этого приложения.
		// path := c.Request.URL.Path
		// if path == "/graphql" || path == "/schema" {
		// 	c.Header("Access-Control-Allow-Origin", "*")
		// }

		// Если конечное приложение  не установило Access-Control-Allow-Origin добавляем его
		if c.GetHeader("Access-Control-Allow-Origin") == "" {
			origin := c.GetHeader("Origin")
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}

		c.Next()
	}
}

// CheckUser проверяет залогинен ли пользователь
func CheckUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Если это префлайт запрос пропускаем запрос без изменений,
		// поскольку браузеры не посылают куки и пользователь не определится.
		if c.Request.Method == "OPTIONS" {
			fmt.Println("c.Request.Method=OPTIONS")
			c.Next()
			return
		}

		// // Если конечное приложение  не установило Access-Control-Allow-Origin добавляем его
		// if c.GetHeader("Access-Control-Allow-Origin") == "" {
		// 	origin := c.GetHeader("Origin")
		// 	if origin != "" {
		// 		c.Header("Access-Control-Allow-Origin", origin)
		// 	}
		// }

		// Кто делает запрос
		userName := session.GetUserName(c)

		// !!! Если пользователь не залогинен ПРЕРЫВАЕМ ЗАПРОС
		if userName == "" {
			// Заголовки добавлены, чтобы пользователь получал внятный ответ
			// если он разлогинен и пытается достучаться до какого то приложения
			// c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			c.Header("Access-Control-Max-Age", "600")
			c.Header("Access-Control-Allow-Headers", "origin, content-type")
			c.Header("Connection", "keep-alive")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Please login: /login "})
			return
		}

		// !!! Если пользователь заблокирован ПРЕРЫВАЕМ ЗАПРОС
		if auth.IsUserEnabled(userName) == false {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Sorry. " + userName + " is disabled."})
			return
		}
		appPath := strings.TrimPrefix(c.Request.URL.Path, "/apps/")

		// К какому приложению делается запрос
		appName := strings.SplitN(appPath, "/", 2)[0]

		// публичное ли это приложение (доступно ли для пользователей без роли)
		isAppPublic := auth.IsAppPublic(appName)

		// роли пользователя в приложении
		userRoles := auth.GetUserRoles(userName, appName)

		// !!! Если приложение не публичное и пользователю не назначены роли ПРЕРЫВАЕМ ЗАПРОС
		if userRoles == "[]" && !isAppPublic {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Sorry. " + userName + " has no roles in " + appName + ". And " + appName + " is not public."})
			return
		}

		// информация о текущем пользователе
		userInfo := auth.GetUserInfo(userName)

		// Добавляем заголовки к запросу
		c.Request.Header.Set("user-roles", userRoles)
		c.Request.Header.Set("user-info", userInfo)

		c.Next()
	}
}
