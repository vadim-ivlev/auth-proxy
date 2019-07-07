package controller

import (
	"auth-proxy/model/db"
	"fmt"
	"strings"

	gq "github.com/graphql-go/graphql"
)

// ************************************************************************

var rootQuery = gq.NewObject(gq.ObjectConfig{
	Name: "Query",
	Fields: gq.Fields{

		"login": &gq.Field{
			Type:        fullUserObject,
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
				panic("Not implemented")
			},
		},

		"logout": &gq.Field{
			Type:        fullUserObject,
			Description: "Выйти",
			Args:        gq.FieldConfigArgument{},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panic("Not implemented")
			},
		},

		"get_user": &gq.Field{
			Type:        fullUserObject,
			Description: "Показать пользователя",
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				fields := getSelectedFields([]string{"get_user"}, params)
				return db.QueryRowMap("SELECT "+fields+" FROM full_user WHERE username = $1 ;", params.Args["username"])
			},
		},

		"get_app": &gq.Field{
			Type:        fullAppObject,
			Description: "Показать приложение",
			Args: gq.FieldConfigArgument{
				"appname": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя приложения",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				fields := getSelectedFields([]string{"get_app"}, params)
				return db.QueryRowMap("SELECT "+fields+" FROM full_app WHERE appname = $1 ;", params.Args["appname"])
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
				wherePart, orderAndLimits := userQueryEnd(params)
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
				wherePart, orderAndLimits := appQueryEnd(params)
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
			Type:        listAppUserRoleGQType,
			Description: "Получить список приложений пользователей и их ролей.",
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
				wherePart, orderAndLimits := app_user_roleQueryEnd(params)
				fields := getSelectedFields([]string{"list_app_user_role", "list"}, params)

				list, err := db.QuerySliceMap("SELECT " + fields + ` FROM "app_user_role"` + wherePart + orderAndLimits)
				if err != nil {
					return nil, err
				}
				count, err := db.QueryRowMap(`SELECT count(*) AS count FROM "app_user_role"` + wherePart)
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
	},
})

func userQueryEnd(params gq.ResolveParams) (wherePart string, orderAndLimits string) {
	var searchConditions []string
	search, ok := params.Args["search"].(string)
	search = strings.Trim(search, " ")
	if ok && len(search) > 0 {
		searchConditions = append(searchConditions,
			fmt.Sprintf("to_tsvector('russian', fullname || ' ' || description  || ' ' || email  || ' ' || username ) @@ plainto_tsquery('russian','%s') ", search))
	}
	// addIntSearchConditionForField(&searchConditions, params, "is_ended")
	if len(searchConditions) > 0 {
		wherePart = " WHERE " + strings.Join(searchConditions, " AND ")
	}
	orderAndLimits = fmt.Sprintf(" ORDER BY %v LIMIT %v OFFSET %v ;", params.Args["order"], params.Args["limit"], params.Args["offset"])
	return wherePart, orderAndLimits
}

func appQueryEnd(params gq.ResolveParams) (wherePart string, orderAndLimits string) {
	var searchConditions []string
	search, ok := params.Args["search"].(string)
	search = strings.Trim(search, " ")
	if ok && len(search) > 0 {
		searchConditions = append(searchConditions,
			fmt.Sprintf("to_tsvector('russian', appname || ' ' || description ) @@ plainto_tsquery('russian','%s') ", search))
	}
	// addIntSearchConditionForField(&searchConditions, params, "is_ended")
	if len(searchConditions) > 0 {
		wherePart = " WHERE " + strings.Join(searchConditions, " AND ")
	}
	orderAndLimits = fmt.Sprintf(" ORDER BY %v LIMIT %v OFFSET %v ;", params.Args["order"], params.Args["limit"], params.Args["offset"])
	return wherePart, orderAndLimits
}

func app_user_roleQueryEnd(params gq.ResolveParams) (wherePart string, orderAndLimits string) {
	var searchConditions []string
	search, ok := params.Args["search"].(string)
	search = strings.Trim(search, " ")
	if ok && len(search) > 0 {
		searchConditions = append(searchConditions,
			`   appname LIKE '%`+search+`%' 
			OR username LIKE '%`+search+`%' 
			OR rolename LIKE '%`+search+`%' `)
	}
	// addIntSearchConditionForField(&searchConditions, params, "is_ended")
	if len(searchConditions) > 0 {
		wherePart = " WHERE " + strings.Join(searchConditions, " AND ")
	}
	orderAndLimits = fmt.Sprintf(" ORDER BY %v LIMIT %v OFFSET %v ;", params.Args["order"], params.Args["limit"], params.Args["offset"])
	return wherePart, orderAndLimits
}
