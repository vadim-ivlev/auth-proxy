package controller

import (
	"auth-proxy/model/auth"
	"auth-proxy/model/db"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	gq "github.com/graphql-go/graphql"
)

// ************************************************************************

var rootQuery = gq.NewObject(gq.ObjectConfig{
	Name: "Query",
	Fields: gq.Fields{

		"login": &gq.Field{
			Type:        gq.String,
			Description: "Войти по имени и паролю",
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя",
				},
				"password": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Пароль",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				c, _ := params.Context.Value("ginContext").(*gin.Context)
				username, _ := params.Args["username"].(string)
				password, _ := params.Args["password"].(string)

				err := auth.Login(c, username, password)
				if err != nil {
					return "Authentication failed", err
				} else {
					return gin.H{"username": username, "message": "Successfully authenticated"}, nil
				}
			},
		},

		"logout": &gq.Field{
			Type:        authMessageObject,
			Description: "Выйти",
			Args:        gq.FieldConfigArgument{},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				c, _ := params.Context.Value("ginContext").(*gin.Context)
				username := auth.GetUserName(c)
				auth.Logout(c)
				return gin.H{"username": username, "message": "Successfully logged out"}, nil
			},
		},

		"get_user": &gq.Field{
			// Type:        fullUserObject,
			Type:        userObject,
			Description: "Показать пользователя",
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotOwnerOrAdmin(params)
				fields := getSelectedFields([]string{"get_user"}, params)
				// return db.QueryRowMap("SELECT "+fields+" FROM full_user WHERE username = $1 ;", params.Args["username"])
				return db.QueryRowMap("SELECT "+fields+` FROM "user" WHERE username = $1 ;`, params.Args["username"])
			},
		},

		"get_logined_user": &gq.Field{
			Type:        userObject,
			Description: "Показать пользователя",
			Args:        gq.FieldConfigArgument{},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotUser(params)
				username := getLoginedUserName(params)
				fields := getSelectedFields([]string{"get_logined_user"}, params)
				return db.QueryRowMap("SELECT "+fields+` FROM "user" WHERE username = $1 ;`, username)
			},
		},

		"get_app": &gq.Field{
			// Type:        fullAppObject,
			Type:        appObject,
			Description: "Показать приложение",
			Args: gq.FieldConfigArgument{
				"appname": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя приложения",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)
				fields := getSelectedFields([]string{"get_app"}, params)
				// return db.QueryRowMap("SELECT "+fields+" FROM full_app WHERE appname = $1 ;", params.Args["appname"])
				return db.QueryRowMap("SELECT "+fields+" FROM app WHERE appname = $1 ;", params.Args["appname"])
			},
		},

		"list_user": &gq.Field{
			Type:        listUserGQType,
			Description: "Получить список пользователей.",
			Args: gq.FieldConfigArgument{
				"search": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Строка полнотекстового поиска.",
				},
				"order": &gq.ArgumentConfig{
					Type:         gq.String,
					Description:  "сортировка строк в определённом порядке. По умолчанию 'username ASC'",
					DefaultValue: "username ASC",
				},
				"limit": &gq.ArgumentConfig{
					Type:         gq.Int,
					Description:  "возвратить не больше заданного числа строк. По умолчанию 100.",
					DefaultValue: 1000,
				},
				"offset": &gq.ArgumentConfig{
					Type:         gq.Int,
					Description:  "пропустить указанное число строк, прежде чем начать выдавать строки. По умолчанию 0.",
					DefaultValue: 0,
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)
				wherePart, orderAndLimits := QueryEnd(params, "fullname,description,email,username")
				fields := getSelectedFields([]string{"list_user", "list"}, params)

				list, err := db.QuerySliceMap("SELECT " + fields + ` FROM "user"` + wherePart + orderAndLimits)
				if err != nil {
					return nil, err
				}
				count, err := db.QueryRowMap(`SELECT count(*) AS count FROM "user"` + wherePart)
				if err != nil {
					return nil, err
				}

				m := map[string]interface{}{
					"length": count["count"],
					"list":   list,
				}

				return m, nil

			},
		},

		"list_app": &gq.Field{
			Type:        listAppGQType,
			Description: "Получить список приложений.",
			Args: gq.FieldConfigArgument{
				"search": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Строка полнотекстового поиска.",
				},
				"order": &gq.ArgumentConfig{
					Type:         gq.String,
					Description:  "сортировка строк в определённом порядке. По умолчанию 'appname ASC'",
					DefaultValue: "appname ASC",
				},
				"limit": &gq.ArgumentConfig{
					Type:         gq.Int,
					Description:  "возвратить не больше заданного числа строк. По умолчанию 100.",
					DefaultValue: 1000,
				},
				"offset": &gq.ArgumentConfig{
					Type:         gq.Int,
					Description:  "пропустить указанное число строк, прежде чем начать выдавать строки. По умолчанию 0.",
					DefaultValue: 0,
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotUser(params)

				wherePart, orderAndLimits := QueryEnd(params, "appname,description")
				fields := getSelectedFields([]string{"list_app", "list"}, params)

				list, err := db.QuerySliceMap("SELECT " + fields + ` FROM "app"` + wherePart + orderAndLimits)
				if err != nil {
					return nil, err
				}
				count, err := db.QueryRowMap(`SELECT count(*) AS count FROM "app"` + wherePart)
				if err != nil {
					return nil, err
				}

				m := map[string]interface{}{
					"length": count["count"],
					"list":   list,
				}

				return m, nil

			},
		},

		"list_app_user_role": &gq.Field{
			Type:        gq.NewList(app_user_role_extendedObject),
			Description: "Получить список приложений пользователей и их ролей.",
			Args: gq.FieldConfigArgument{
				"appname": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Идентификатор приложения",
				},
				"username": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Идентификатор пользователя",
				},
				"rolename": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Идентификатор роли",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotOwnerOrAdmin(params)
				fields := getSelectedFields([]string{"list_app_user_role"}, params)
				wherePart := list_app_user_roleWherePart(params)
				query := fmt.Sprintf(`SELECT DISTINCT %s FROM app_user_role_extended %s `, fields, wherePart)
				// println(query)
				return db.QuerySliceMap(query)
			},
		},
	},
})

func list_app_user_roleWherePart(params gq.ResolveParams) (wherePart string) {
	var searchConditions []string

	for paramName, v := range params.Args {
		s, ok := v.(string)
		s = strings.Trim(s, " ")
		if ok && len(s) > 0 {
			searchConditions = append(searchConditions, fmt.Sprintf(" %s = '%s' ", paramName, s))
		}
	}
	if len(searchConditions) > 0 {
		wherePart = " WHERE " + strings.Join(searchConditions, " AND ")
	}
	return wherePart
}

func QueryEnd(params gq.ResolveParams, fieldList string) (wherePart string, orderAndLimits string) {
	var searchConditions []string
	search, ok := params.Args["search"].(string)
	search = strings.Trim(search, " ")
	if ok && len(search) > 0 {
		search = strings.ReplaceAll(search, " ", "%")
		searchConditions = append(searchConditions, Like(fieldList, search))
	}
	// addIntSearchConditionForField(&searchConditions, params, "is_ended")
	if len(searchConditions) > 0 {
		wherePart = " WHERE " + strings.Join(searchConditions, " AND ")
	}
	orderAndLimits = fmt.Sprintf(" ORDER BY %v LIMIT %v OFFSET %v ;", params.Args["order"], params.Args["limit"], params.Args["offset"])
	return wherePart, orderAndLimits
}

func Like(fieldsString, search string) string {
	fields := strings.Split(fieldsString, ",")
	var chunks []string
	for _, field := range fields {
		chunks = append(chunks, ` LOWER(`+field+`) LIKE LOWER('%`+search+`%') `)
	}
	return strings.Join(chunks, " OR ")
}
