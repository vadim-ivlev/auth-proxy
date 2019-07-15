package controller

import (
	"auth-proxy/model/db"

	gq "github.com/graphql-go/graphql"
)

var rootMutation = gq.NewObject(gq.ObjectConfig{
	Name: "Mutation",
	Fields: gq.Fields{

		"create_user": &gq.Field{
			Description: "Создать пользователя",
			// Type:        fullUserObject,
			Type: userObject,
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя (уникальное)",
				},
				"password": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Пароль",
				},
				"email": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Емайл",
				},
				"fullname": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Полное имя",
				},
				"description": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Описание",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfEmpty(params.Args["username"], "Имя пользователя не должно быть пустым")
				// return createRecord("username", params, "user", "full_user")
				return createRecord("username", params, "user", "user")
			},
		},

		"update_user": &gq.Field{
			Description: "Обновить пользователя",
			// Type:        fullUserObject,
			Type: userObject,
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя (уникальное)",
				},
				"password": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Пароль",
				},
				"email": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Емайл",
				},
				"fullname": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Полное имя",
				},
				"description": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Описание",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotOwnerOrAdmin(params)
				// return updateRecord("username", params, "user", "full_user")
				return updateRecord("username", params, "user", "user")
			},
		},

		"delete_user": &gq.Field{
			Description: "Удалить пользователя",
			// Type:        fullUserObject,
			Type: userObject,
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя (уникальное)",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotOwnerOrAdmin(params)
				// return deleteRecord("username", params, "user", "full_user")
				return deleteRecord("username", params, "user", "user")
			},
		},

		"create_app": &gq.Field{
			Description: "Создать приложение",
			// Type:        fullAppObject,
			Type: appObject,
			Args: gq.FieldConfigArgument{
				"appname": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя приложения (уникальное)",
				},
				"description": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Описание",
				},
				"url": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "url",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)
				panicIfEmpty(params.Args["appname"], "Имя приложения не должно быть пустым")
				// return createRecord("appname", params, "app", "full_app")
				return createRecord("appname", params, "app", "app")
			},
		},

		"update_app": &gq.Field{
			Description: "Обновить пользователя",
			// Type:        fullAppObject,
			Type: appObject,
			Args: gq.FieldConfigArgument{
				"appname": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя (уникальное)",
				},
				"description": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Описание",
				},
				"url": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "url",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)
				// return updateRecord("appname", params, "app", "full_app")
				return updateRecord("appname", params, "app", "app")
			},
		},

		"delete_app": &gq.Field{
			Description: "Удалить пользователя",
			// Type:        fullAppObject,
			Type: appObject,
			Args: gq.FieldConfigArgument{
				"appname": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя (уникальное)",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)
				// return deleteRecord("appname", params, "app", "full_app")
				return deleteRecord("appname", params, "app", "app")
			},
		},

		"create_app_user_role": &gq.Field{
			Description: "Создать роль пользователя для приложения",
			Type:        app_user_roleObject,
			Args: gq.FieldConfigArgument{
				"appname": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя приложения (уникальное)",
				},
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя (уникальное)",
				},
				"rolename": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя роли",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)
				return db.CreateRow("app_user_role", params.Args)
			},
		},

		"delete_app_user_role": &gq.Field{
			Description: "Удалить роль пользователя для приложения",
			Type:        app_user_roleObject,
			Args: gq.FieldConfigArgument{
				"appname": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя приложения (уникальное)",
				},
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя (уникальное)",
				},
				"rolename": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя роли",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)
				// sqlText := fmt.Sprintf(
				// 	`DELETE FROM app_user_role WHERE appname = '%v' AND username = '%v' AND rolename = '%v' RETURNING * ;`,
				// 	params.Args["appname"],
				// 	params.Args["username"],
				// 	params.Args["rolename"],
				// )
				// return db.QueryRowMap(sqlText)
				a, u, r := params.Args["appname"], params.Args["username"], params.Args["rolename"]

				_, err := db.QueryExec(
					`DELETE FROM app_user_role WHERE appname = $1 AND username = $2 AND rolename = $3 ;`, a, u, r)
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"appname": a, "username": u, "rolename": r}, nil
			},
		},
	},
})
