package server

import (
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/reqcounter"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CountersMiddleware считает запросы
func CountersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// считаем все кроме запросов на выдачу статистики
		if c.Request.URL.Path != "/stat" && c.Request.URL.Path != "/metrics" {
			reqcounter.IncrementCounters()
		}
		c.Next()
	}
}

// RedirectsMiddleware перенаправляет браузер на другие ресурсы
func RedirectsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// если url есть в списке редиректов, перенаправляем браузер
		url, ok := Redirects[c.Request.URL.Path]
		if ok {
			c.Redirect(http.StatusMovedPermanently, url)
			c.Abort()
		}
		c.Next()
	}
}

// HeadersMiddleware добавляет HTTP заголовки к ответу сервера
func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Если конечное приложение  не установило Access-Control-Allow-Origin добавляем его
		aoHeader := c.GetHeader("Access-Control-Allow-Origin")
		if aoHeader == "" || aoHeader == "*" {
			origin := c.GetHeader("Origin")
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}

		// Если конечное приложение  не установило Access-Control-Allow-Credentials добавляем его
		if c.GetHeader("Access-Control-Allow-Credentials") == "" {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Next()
	}
}

// CheckUserMiddleware проверяет залогинен ли пользователь
func CheckUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Если это префлайт запрос пропускаем запрос без изменений,
		// поскольку браузеры не посылают куки и пользователь не определится.
		if c.Request.Method == "OPTIONS" {
			fmt.Println("c.Request.Method=OPTIONS")
			c.Next()
			return
		}

		// К какому приложению делается запрос
		appPath := strings.TrimPrefix(c.Request.URL.Path, "/apps/")
		appName := strings.SplitN(appPath, "/", 2)[0]

		// Кто делает запрос
		userName := GetSessionVariable(c, "user")
		userInfo := auth.GetUserInfoString(userName, appName)

		// Добавляем заголовки к запросу
		c.Request.Header.Set("user-info", userInfo)

		// публичное ли это приложение? (доступно ли для неавторизованных пользователей?)
		isAppPublic := auth.IsAppPublic(appName)
		if isAppPublic {
			fmt.Println("Публичное приложение:", appName)
			c.Next()
			return
		}

		// !!! Если пользователь не залогинен ПРЕРЫВАЕМ ЗАПРОС
		if userName == "" {
			// Заголовки добавлены, чтобы пользователь получал внятный ответ
			// если он разлогинен и пытается достучаться до какого то приложения

			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			c.Header("Access-Control-Max-Age", "600")
			c.Header("Access-Control-Allow-Headers", "origin, content-type")
			c.Header("Connection", "keep-alive")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Please login: /login "})
			return
		}

		// !!! Если пользователь заблокирован ПРЕРЫВАЕМ ЗАПРОС
		if !auth.IsUserEnabled(userName) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Sorry. " + userName + " is disabled or DB is unavailable. Please ask admins."})
			return
		}

		// роли пользователя в приложении
		userRoles := auth.GetUserRolesString(userName, appName)

		// !!! Если приложение не публичное и пользователю не назначены роли ПРЕРЫВАЕМ ЗАПРОС
		if userRoles == "[]" && !isAppPublic {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Sorry. " + userName + " has no roles in " + appName + ". And " + appName + " is not public."})
			return
		}

		c.Next()
	}
}
