package controller

import (
	"auth-proxy/model/auth"
	"auth-proxy/model/db"
	"auth-proxy/model/mail"
	"errors"
	"log"

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
				"disabled": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Если не равно 0, пользователь отключен",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfEmpty(params.Args["username"], "Имя пользователя не должно быть пустым")
				panicIfEmpty(params.Args["password"], "Пароль не должен быть пустым")
				processPassword(params)

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
				"disabled": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Если не равно 0, пользователь отключен",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotOwnerOrAdmin(params)
				processPassword(params)
				// return updateRecord("username", params, "user", "full_user")
				return updateRecord("username", params, "user", "user")
			},
		},

		"generate_password": &gq.Field{
			Description: "Сгенерировать новый пароль пользователю и выслать по электронной почте",
			// Type:        userObject,
			Type: gq.String,
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя (уникальное)",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				username := params.Args["username"].(string)
				if username == "" {
					return "Поле username должно быть заполнено", errors.New("Не указаны имя или емайл пользователя")
				}
				foundUsername, email, password, err := auth.GenerateNewPassword(username)
				if err != nil {
					return "Could not generate new password", err
				}
				err = mail.SendPassword(foundUsername, email, password)
				if err != nil {
					return "Could not send email to:" + email, err
				}

				return "Новый пароль для "+foundUsername+" выслан по адресу: " + email, nil

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
				"rebase": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Y - чтобы иправить ссылки на относительные на HTML страницах",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)
				panicIfEmpty(params.Args["appname"], "Имя приложения не должно быть пустым")
				// return createRecord("appname", params, "app", "full_app")
				res, err := createRecord("appname", params, "app", "app")
				if err == nil {
					app, _ := params.Args["appname"].(string)
					url, _ := params.Args["url"].(string)
					rebase, _ := params.Args["rebase"].(string)
					if url != "" {
						proxies[app] = createProxy(url, app, rebase)
						log.Printf("Proxy created appname=%v target=%v rebase=%v", app, url, rebase)
					}
				}
				return res, err
			},
		},

		"update_app": &gq.Field{
			Description: "Обновить приложение",
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
					Description: "url проксируемого приложения",
				},
				"rebase": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Y - чтобы иправить ссылки на относительные на HTML страницах",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)
				// return updateRecord("appname", params, "app", "full_app")
				res, err := updateRecord("appname", params, "app", "app")
				if err == nil {
					app, _ := params.Args["appname"].(string)
					rebase, _ := params.Args["rebase"].(string)
					u, ok := params.Args["url"]
					if ok {
						url, _ := u.(string)
						if url == "" {
							delete(proxies, app)
							log.Printf("Proxy deleted appname=%s", app)
						} else {
							proxies[app] = createProxy(url, app, rebase)
							log.Printf("Proxy created appname=%v target=%v", app, url)
						}
					}
				}
				return res, err
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
				res, err := deleteRecord("appname", params, "app", "app")
				if err == nil {
					app, _ := params.Args["appname"].(string)
					delete(proxies, app)
					log.Printf("Proxy deleted appname=%s", app)
				}
				return res, err
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
