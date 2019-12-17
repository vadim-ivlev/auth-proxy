package server

import (
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/reqcounter"
	"auth-proxy/pkg/signature"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CountersMiddleware считает запросы
func CountersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// считаем все кроме запросов на выдачу статистики
		if c.Request.URL.Path != "/stat" {
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

		// Кто делает запрос
		userName := GetSessionVariable(c, "user")

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
		userRoles := auth.GetUserRolesString(userName, appName)

		// !!! Если приложение не публичное и пользователю не назначены роли ПРЕРЫВАЕМ ЗАПРОС
		if userRoles == "[]" && !isAppPublic {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Sorry. " + userName + " has no roles in " + appName + ". And " + appName + " is not public."})
			return
		}

		// информация о текущем пользователе
		userInfo := auth.GetUserInfoString(userName, appName)

		// Добавляем заголовки к запросу
		c.Request.Header.Set("user-roles", userRoles)
		c.Request.Header.Set("user-info", userInfo)

		// Подписываем запрос
		if auth.IsRequestToAppSigned(appName) {
			r := c.Request
			signature.Sign(r)
			signature.Verify(r)
		}

		c.Next()
	}
}
