package server

import (
	"auth-proxy/pkg/prometeo"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

var queryObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"login":                  login(),
		"logout":                 logout(),
		"is_selfreg_allowed":     is_selfreg_allowed(),
		"get_stat":               get_stat(),
		"list_oauth_providers":   list_oauth_providers(),
		"get_params":             get_params(),
		"set_params":             set_params(),
		"is_captcha_required":    is_captcha_required(),
		"is_pin_required":        is_pin_required(),
		"get_user":               get_user(),
		"get_logined_user":       get_logined_user(),
		"get_app":                get_app(),
		"list_user":              list_user(),
		"list_user_by_ids":       list_user_by_ids(),
		"list_user_by_usernames": list_user_by_usernames(),
		"list_app":               list_app(),
		"list_app_user_role":     list_app_user_role(),
	},
})

var mutationObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"create_user":          create_user(),
		"update_user":          update_user(),
		"generate_password":    generate_password(),
		"delete_user":          delete_user(),
		"create_app":           create_app(),
		"update_app":           update_app(),
		"delete_app":           delete_app(),
		"create_app_user_role": create_app_user_role(),
		"delete_app_user_role": delete_app_user_role(),
	},
})

var schema graphql.Schema

func SchemaInit(noIntrospection bool) {
	if noIntrospection {
		fmt.Println("!!!!!!!!!!!!!!! SUPPRESSING GraphQL INTROSPECTION !!!!!!!!!!!!!!!!!!")
		graphql.SchemaMetaFieldDef.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
			return nil, nil
		}
		graphql.TypeMetaFieldDef.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
			return nil, nil
		}
	}
	var err error
	schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryObject,
		Mutation: mutationObject,
	})

	if err != nil {
		log.Println("SchemaInit ERROR:", err)
	}
}

// graphqlResult HTTP handler. исполняет graphql запрос
func graphqlResult(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 100*1024*1024)

	query, variables := getPayload3(c)
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		Context:        context.WithValue(context.Background(), "ginContext", c),
		VariableValues: variables,
	})

	if len(result.Errors) > 0 {
		// инкрементируем счетчик ошибок GraphQL
		prometeo.GraphQLErrorsTotal.Inc()
	}

	c.JSON(200, result)
}
