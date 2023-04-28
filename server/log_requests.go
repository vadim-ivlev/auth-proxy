// функции для работы с логами запросов
package server

import (
	"auth-proxy/pkg/app"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

// logCreateUser - логирование запроса на создание пользователя
func logCreateUser(params graphql.ResolveParams) {
	ArgToLowerCase(params, "username")
	c, _ := params.Context.Value("ginContext").(*gin.Context)
	email, _ := params.Args["email"].(string)
	password, _ := params.Args["password"].(string)
	fullname, _ := params.Args["fullname"].(string)
	description, _ := params.Args["description"].(string)

	fullPath := c.FullPath()
	headers := app.Serialize(c.Request.Header)
	// request := app.Serialize(c.Request)

	log.Printf("email=%s, password=%s fullname=%s description=%s \n", email, password, fullname, description)
	log.Printf("ip=%s, user_agent=%s\n", c.ClientIP(), c.Request.UserAgent())
	log.Printf(" FullPath: %v\n", fullPath)
	log.Printf(" Headers: %v\n", headers)

}
