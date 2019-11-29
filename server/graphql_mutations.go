package server

import (
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/db"
	"auth-proxy/pkg/mail"
	"errors"
	"log"

	gq "github.com/graphql-go/graphql"
)

var rootMutation = gq.NewObject(gq.ObjectConfig{
	Name: "Mutation",
	Fields: gq.Fields{

		"create_user": &gq.Field{
			Description: "Создать пользователя",
			Type:        userObject,
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
				if SelfRegistrationAllowed || isAuthAdmin(params) {
					panicIfEmpty(params.Args["username"], "Заполните поле Имя пользователя")
					panicIfEmpty(params.Args["password"], "Введите пароль")
					panicIfEmpty(params.Args["email"], "Заполните поле Email")

					ArgToLowerCase(params, "username")
					ArgToLowerCase(params, "email")

					password := processPassword(params)
					clearUserCache(params)
					res, err := createRecord("username", params, "user", "user")
					if err == nil {
						sendMessageToNewUser(params, password)
					}
					return res, err
				} else {
					return nil, errors.New("Sorry. Self registration is not allowed. Please ask administrators.")
				}

			},
		},

		"update_user": &gq.Field{
			Description: "Обновить пользователя",
			Type:        userObject,
			Args: gq.FieldConfigArgument{
				"old_username": &gq.ArgumentConfig{
					Type:         gq.NewNonNull(gq.String),
					Description:  "Имя пользователя до обновления (уникальное)",
					DefaultValue: "",
				},
				"username": &gq.ArgumentConfig{
					Type:         gq.NewNonNull(gq.String),
					Description:  "Имя пользователя (уникальное)",
					DefaultValue: "",
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

				ArgToLowerCase(params, "old_username")
				ArgToLowerCase(params, "username")
				ArgToLowerCase(params, "email")

				processPassword(params)
				clearUserCache(params)

				old_username, _ := params.Args["old_username"].(string)
				if old_username == "" {
					old_username, _ = params.Args["username"].(string)
				}
				if old_username == "" {
					return nil, errors.New("username is blank")
				}

				delete(params.Args, "old_username")
				return updateRecord(old_username, "username", params, "user", "user")
			},
		},

		"generate_password": &gq.Field{
			Description: "Сгенерировать новый пароль пользователю и выслать по электронной почте",
			Type:        gq.String,
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя (уникальное)",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				ArgToLowerCase(params, "username")

				username := params.Args["username"].(string)
				if username == "" {
					return "Поле username должно быть заполнено", errors.New("Не указаны имя или емайл пользователя")
				}
				foundUsername, email, password, err := auth.GenerateNewPassword(username)
				if err != nil {
					return "Could not generate new password", err
				}
				err = mail.SendMessage("new_password", foundUsername, email, password)
				if err != nil {
					return "Could not send email to:" + email, err
				}

				return "Новый пароль для " + foundUsername + " выслан по адресу: " + email, nil

			},
		},

		"delete_user": &gq.Field{
			Description: "Удалить пользователя",
			Type:        userObject,
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя (уникальное)",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				ArgToLowerCase(params, "username")

				panicIfNotOwnerOrAdmin(params)
				clearUserCache(params)
				return deleteRecord("username", params, "user", "user")
			},
		},

		"create_app": &gq.Field{
			Description: "Создать приложение",
			Type:        appObject,
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
				"public": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Y - чтобы сделать приложение доступным для пользователей без роли",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)
				panicIfEmpty(params.Args["appname"], "Имя приложения не должно быть пустым")
				res, err := createRecord("appname", params, "app", "app")
				if err == nil {
					app, _ := params.Args["appname"].(string)
					url, _ := params.Args["url"].(string)
					rebase, _ := params.Args["rebase"].(string)
					if url != "" {
						proxies[app] = createProxy(url, app, rebase)
						log.Printf("Proxy created appname=%v target=%v rebase=%v", app, url, rebase)
					}
					clearAppCache(params)
				}
				return res, err
			},
		},

		"update_app": &gq.Field{
			Description: "Обновить приложение",
			Type:        appObject,
			Args: gq.FieldConfigArgument{
				"old_appname": &gq.ArgumentConfig{
					Type:         gq.NewNonNull(gq.String),
					Description:  "Имя приложения до обновления (уникальное)",
					DefaultValue: "",
				},
				"appname": &gq.ArgumentConfig{
					Type:         gq.NewNonNull(gq.String),
					Description:  "Имя приложения (уникальное)",
					DefaultValue: "",
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
				"public": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Y - чтобы сделать приложение доступным для пользователей без роли",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)

				old_appname, _ := params.Args["old_appname"].(string)
				if old_appname == "" {
					old_appname, _ = params.Args["appname"].(string)
				}
				if old_appname == "" {
					return nil, errors.New("appname is blank")
				}

				delete(params.Args, "old_appname")
				res, err := updateRecord(old_appname, "appname", params, "app", "app")

				// изменяем массив прокси
				if err == nil {
					appname, _ := params.Args["appname"].(string)

					// если appname изменилось перепорождаем прокси со старыми Url, Rebase
					if appname != old_appname {
						old_proxy := *proxies[old_appname]
						log.Println("renaming old_url, old_appname, old_rebase", old_proxy.Url, old_appname, old_proxy.Rebase)
						delete(proxies, old_appname)
						proxies[appname] = createProxy(old_proxy.Url, appname, old_proxy.Rebase)
						// clear old_appname cache
						auth.Cache.Delete("is-" + old_appname + "-public")
					}

					// Если обновляется Url, Rebase пререпорождаем прокси
					r, ok_r := params.Args["rebase"]
					u, ok_u := params.Args["url"]
					if ok_r || ok_u {
						rebase, _ := r.(string)
						url, _ := u.(string)
						if url == "" {
							delete(proxies, appname)
							log.Printf("Proxy deleted appname=%s", appname)
						} else {
							proxies[appname] = createProxy(url, appname, rebase)
							log.Printf("Proxy created appname=%v target=%v", appname, url)
						}
					}
					clearAppCache(params)
				}
				return res, err
			},
		},

		"delete_app": &gq.Field{
			Description: "Удалить пользователя",
			Type:        appObject,
			Args: gq.FieldConfigArgument{
				"appname": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя (уникальное)",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdmin(params)
				res, err := deleteRecord("appname", params, "app", "app")
				if err == nil {
					app, _ := params.Args["appname"].(string)
					delete(proxies, app)
					log.Printf("Proxy deleted appname=%s", app)
					clearAppCache(params)
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
				ArgToLowerCase(params, "username")
				clearAppUserRolesCache(params)
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
				ArgToLowerCase(params, "username")

				a, u, r := params.Args["appname"], params.Args["username"], params.Args["rolename"]

				_, err := db.QueryExec(
					`DELETE FROM app_user_role WHERE appname = $1 AND username = $2 AND rolename = $3 ;`, a, u, r)
				if err != nil {
					return nil, err
				}
				clearAppUserRolesCache(params)
				return map[string]interface{}{"appname": a, "username": u, "rolename": r}, nil
			},
		},
	},
})

// clearAppUserRolesCache чистим кэш
func clearAppUserRolesCache(params gq.ResolveParams) {
	cacheKey := params.Args["username"].(string) + "-" + params.Args["appname"].(string) + "-roles"
	auth.Cache.Delete(cacheKey)
}

func clearUserCache(params gq.ResolveParams) {
	username, _ := params.Args["username"].(string)
	auth.Cache.Delete(username + "-info")
	auth.Cache.Delete(username + "-enabled")
}

func clearAppCache(params gq.ResolveParams) {
	app, _ := params.Args["appname"].(string)
	auth.Cache.Delete("is-" + app + "-public")
}

func sendMessageToNewUser(params gq.ResolveParams, password string) {
	email := params.Args["email"].(string)
	if email == "" {
		return
	}
	err := mail.SendMessage("new_user", params.Args["username"].(string), email, password)
	if err != nil {
		log.Println("create_user SendMessage error:", err)
	}
}
