// функции для работы с логами запросов
package server

import (
	"auth-proxy/pkg/app"
	"auth-proxy/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

// logCreateUser - логирование запроса на создание пользователя
func logCreateUser(params graphql.ResolveParams) {

	c, ok := params.Context.Value("ginContext").(*gin.Context)
	if !ok {
		log.Println("logCreateUser error: ginContext not found")
		return
	}

	// получение параметров запроса
	email, _ := params.Args["email"].(string)
	// password, _ := params.Args["password"].(string)
	password := ""
	fullname, _ := params.Args["fullname"].(string)
	description, _ := params.Args["description"].(string)
	fullPath := c.FullPath()
	headers := app.Serialize(c.Request.Header)
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()
	referer := c.Request.Referer()

	// запись в лог
	_, err := db.QueryExec("INSERT INTO create_user_log (email, password, fullname, description, ip, user_agent, full_path, referer, headers) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		email, password, fullname, description, ip, userAgent, fullPath, referer, headers)
	if err != nil {
		log.Println("logCreateUser err=", err)
	}

}
