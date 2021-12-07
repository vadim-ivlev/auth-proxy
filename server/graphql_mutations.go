package server

import (
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/db"
	"auth-proxy/pkg/mail"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

func create_user() *graphql.Field {
	return &graphql.Field{
		Description: "Создать пользователя",
		Type:        userObject,
		Args: graphql.FieldConfigArgument{
			// "username": &graphql.ArgumentConfig{
			// 	Type:        graphql.NewNonNull(graphql.String),
			// 	Description: "Имя пользователя (уникальное)",
			// },
			"password": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Пароль",
			},
			"email": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Email пользователя",
			},
			"fullname": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Полное имя (Фамилия Имя)",
			},
			"description": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Описание",
			},
			"disabled": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Если не равно 0, пользователь отключен",
			},
			"pinrequired": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "требуется ли PIN Google Authenticator",
			},
			// "pinhash_temp": &graphql.ArgumentConfig{
			// 	Type:        graphql.String,
			// 	Description: "хэш для установки Google Authenticator",
			// },
			// "pinhash": &graphql.ArgumentConfig{
			// 	Type:        graphql.String,
			// 	Description: "хэш для для проверки PIN после установки Google Authenticator",
			// },
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if SelfRegistrationAllowed || isAuthAdmin(params) {
				// panicIfEmpty(params.Args["username"], "Заполните поле Имя пользователя")
				panicIfEmpty(params.Args["password"], "Введите пароль")
				panicIfEmpty(params.Args["email"], "Заполните поле Email")

				// Только админ может включить/отключить проверку пина
				if !isAuthAdmin(params) {
					delete(params.Args, "pinrequired")
				}
				// // TODO: может нужно вообще запретить изменение pinhash, pinhash_temp
				// // если значение пусто не меняем его
				// pinhash, _ := params.Args["pinhash"].(string)
				// if pinhash == "" {
				// 	delete(params.Args, "pinhash")
				// }
				// // если значение пусто не меняем его
				// pinhash_temp, _ := params.Args["pinhash_temp"].(string)
				// if pinhash_temp == "" {
				// 	delete(params.Args, "pinhash_temp")
				// }

				ArgToLowerCase(params, "email")
				params.Args["username"] = params.Args["email"]
				convertPasswordToHash(params)
				emailhash := uuid.New().String()
				params.Args["emailhash"] = emailhash

				clearUserCache(params)
				res, err := createRecord("username", params, "user", "user")
				if err == nil {
					err := mail.SendNewUserEmail(params.Args["email"].(string), emailhash)
					if err != nil {
						log.Println("create_user SendNewUserEmail error:", err)
					}
				}
				return res, err
			}
			return nil, errors.New("self registration is not allowed. Please ask administrators")
		},
	}
}

func update_user() *graphql.Field {
	return &graphql.Field{
		Description: "Обновить пользователя",
		Type:        userObject,
		Args: graphql.FieldConfigArgument{
			// "old_username": &graphql.ArgumentConfig{
			// 	Type:         graphql.NewNonNull(graphql.String),
			// 	Description:  "Имя пользователя до обновления (уникальное)",
			// 	DefaultValue: "",
			// },
			// "username": &graphql.ArgumentConfig{
			// 	Type:         graphql.NewNonNull(graphql.String),
			// 	Description:  "Имя пользователя (уникальное)",
			// 	DefaultValue: "",
			// },
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор пользователя",
			},
			"password": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Пароль",
			},
			// "email": &graphql.ArgumentConfig{
			// 	Type:        graphql.String,
			// 	Description: "Email пользователя",
			// },
			"fullname": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Полное имя (Фамилия Имя)",
			},
			"description": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Описание",
			},
			"disabled": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Если не равно 0, пользователь отключен",
			},
			"pinrequired": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "требуется ли PIN Google Authenticator",
			},
			// "pinhash_temp": &graphql.ArgumentConfig{
			// 	Type:        graphql.String,
			// 	Description: "хэш для установки Google Authenticator",
			// },
			// "pinhash": &graphql.ArgumentConfig{
			// 	Type:        graphql.String,
			// 	Description: "хэш для для проверки PIN после установки Google Authenticator",
			// },
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotOwnerOrAdmin(params)

			// Только админ может включить/отключить проверку пина
			if !isAuthAdmin(params) {
				delete(params.Args, "pinrequired")
			}
			// // TODO: может нужно вообще запретить изменение pinhash, pinhash_temp
			// // если значение пусто не меняем его
			// pinhash, _ := params.Args["pinhash"].(string)
			// if pinhash == "" {
			// 	delete(params.Args, "pinhash")
			// }
			// // если значение пусто не меняем его
			// pinhash_temp, _ := params.Args["pinhash_temp"].(string)
			// if pinhash_temp == "" {
			// 	delete(params.Args, "pinhash_temp")
			// }

			// ArgToLowerCase(params, "old_username")
			// ArgToLowerCase(params, "username")
			// ArgToLowerCase(params, "email")

			convertPasswordToHash(params)
			clearUserCache(params)

			// oldUsername, _ := params.Args["old_username"].(string)
			// if oldUsername == "" {
			// 	oldUsername, _ = params.Args["username"].(string)
			// }
			// if oldUsername == "" {
			// 	return nil, errors.New("username is blank")
			// }

			// delete(params.Args, "old_username")
			// return updateRecord(oldUsername, "username", params, "user", "user")

			id := params.Args["id"].(int)
			return db.UpdateRowByID("id", "user", id, params.Args)
		},
	}
}

// func generate_password() *graphql.Field {
// 	return &graphql.Field{
// 		Description: "Сгенерировать новый пароль пользователю и выслать по электронной почте",
// 		Type:        graphql.String,
// 		Args: graphql.FieldConfigArgument{
// 			"username": &graphql.ArgumentConfig{
// 				Type:        graphql.NewNonNull(graphql.String),
// 				Description: "Имя пользователя (уникальное)",
// 			},
// 		},
// 		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
// 			ArgToLowerCase(params, "username")

// 			username := params.Args["username"].(string)
// 			if username == "" {
// 				return "Поле username должно быть заполнено", errors.New("Не указаны имя или Email пользователя")
// 			}
// 			foundUsername, email, password, err := auth.GenerateNewPassword(username)
// 			if err != nil {
// 				return "Could not generate new password", err
// 			}
// 			// err = mail.SendMessage("new_password", foundUsername, email, password)
// 			err = mail.SendNewPasswordEmail(foundUsername, email, password)
// 			if err != nil {
// 				return "Could not send email to:" + email, err
// 			}

// 			return "Новый пароль для " + foundUsername + " выслан.", nil

// 		},
// 	}
// }

func delete_user() *graphql.Field {
	return &graphql.Field{
		Description: "Удалить пользователя",
		Type:        userObject,
		Args: graphql.FieldConfigArgument{
			"username": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя пользователя (уникальное)",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ArgToLowerCase(params, "username")

			panicIfNotOwnerOrAdmin(params)
			clearUserCache(params)
			return deleteRecord("username", params, "user", "user")
		},
	}
}

func create_app() *graphql.Field {
	return &graphql.Field{
		Description: "Создать приложение",
		Type:        appObject,
		Args: graphql.FieldConfigArgument{
			"appname": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя приложения (уникальное)",
			},
			"description": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Описание",
			},
			"url": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "url",
			},
			"rebase": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Y - чтобы иправить ссылки на относительные на HTML страницах",
			},
			"public": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Y - чтобы сделать приложение доступным для пользователей без роли",
			},
			"sign": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Y - для цифровой подписи запросов к приложению",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
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
	}
}

func update_app() *graphql.Field {
	return &graphql.Field{
		Description: "Обновить приложение",
		Type:        appObject,
		Args: graphql.FieldConfigArgument{
			"old_appname": &graphql.ArgumentConfig{
				Type:         graphql.NewNonNull(graphql.String),
				Description:  "Имя приложения до обновления (уникальное)",
				DefaultValue: "",
			},
			"appname": &graphql.ArgumentConfig{
				Type:         graphql.NewNonNull(graphql.String),
				Description:  "Имя приложения (уникальное)",
				DefaultValue: "",
			},
			"description": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Описание",
			},
			"url": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "url проксируемого приложения",
			},
			"rebase": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Y - чтобы иправить ссылки на относительные на HTML страницах",
			},
			"public": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Y - чтобы сделать приложение доступным для пользователей без роли",
			},
			"sign": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Y - для цифровой подписи запросов к приложению",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdmin(params)

			oldAppname, _ := params.Args["old_appname"].(string)
			if oldAppname == "" {
				oldAppname, _ = params.Args["appname"].(string)
			}
			if oldAppname == "" {
				return nil, errors.New("appname is blank")
			}

			delete(params.Args, "old_appname")
			res, err := updateRecord(oldAppname, "appname", params, "app", "app")

			// изменяем массив прокси
			if err == nil {
				appname, _ := params.Args["appname"].(string)

				// если appname изменилось перепорождаем прокси со старыми Url, Rebase
				if appname != oldAppname {
					oldProxy := *proxies[oldAppname]
					log.Println("renaming old_url, old_appname, old_rebase", oldProxy.Url, oldAppname, oldProxy.Rebase)
					delete(proxies, oldAppname)
					proxies[appname] = createProxy(oldProxy.Url, appname, oldProxy.Rebase)
					// clear old_appname cache
					auth.Cache.Delete("is-" + oldAppname + "-public")
					auth.Cache.Delete("is-request-to-" + oldAppname + "-signed")
				}

				// Если обновляется Url, Rebase перепорождаем прокси
				r, okRebase := params.Args["rebase"]
				u, okURL := params.Args["url"]
				if okRebase || okURL {
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
	}
}

func delete_app() *graphql.Field {
	return &graphql.Field{
		Description: "Удалить пользователя",
		Type:        appObject,
		Args: graphql.FieldConfigArgument{
			"appname": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя пользователя (уникальное)",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
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
	}
}

func create_app_user_role() *graphql.Field {
	return &graphql.Field{
		Description: "Создать роль пользователя для приложения",
		Type:        appUserRoleObject,
		Args: graphql.FieldConfigArgument{
			"appname": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя приложения (уникальное)",
			},
			"username": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя пользователя (уникальное)",
			},
			"rolename": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя роли",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdmin(params)
			ArgToLowerCase(params, "username")
			clearAppUserRolesCache(params)
			clearUserCache(params)
			return db.CreateRow("app_user_role", params.Args)
		},
	}
}

func delete_app_user_role() *graphql.Field {
	return &graphql.Field{
		Description: "Удалить роль пользователя для приложения",
		Type:        appUserRoleObject,
		Args: graphql.FieldConfigArgument{
			"appname": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя приложения (уникальное)",
			},
			"username": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя пользователя (уникальное)",
			},
			"rolename": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя роли",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdmin(params)
			ArgToLowerCase(params, "username")

			a, u, r := params.Args["appname"], params.Args["username"], params.Args["rolename"]

			_, err := db.QueryExec(
				`DELETE FROM app_user_role WHERE appname = $1 AND username = $2 AND rolename = $3 ;`, a, u, r)
			if err != nil {
				return nil, err
			}
			clearAppUserRolesCache(params)
			clearUserCache(params)
			return map[string]interface{}{"appname": a, "username": u, "rolename": r}, nil
		},
	}
}

// clearAppUserRolesCache чистим кэш
func clearAppUserRolesCache(params graphql.ResolveParams) {
	cacheKey := params.Args["username"].(string) + "-" + params.Args["appname"].(string) + "-roles"
	auth.Cache.Delete(cacheKey)
}

func clearUserCache(params graphql.ResolveParams) {
	username, _ := params.Args["username"].(string)
	auth.Cache.Delete(username + "-info")
	auth.Cache.Delete(username + "-enabled")
}

func clearAppCache(params graphql.ResolveParams) {
	app, _ := params.Args["appname"].(string)
	auth.Cache.Delete("is-" + app + "-public")
	auth.Cache.Delete("is-request-to-" + app + "-signed")
}

// func sendMessageToNewUser(params gq.ResolveParams, password string) {
// 	email := params.Args["email"].(string)
// 	if email == "" {
// 		return
// 	}
// 	// err := mail.SendMessage("new_user", params.Args["username"].(string), email, password)
// 	err := mail.SendNewUserEmail(params.Args["username"].(string), email, password)
// 	if err != nil {
// 		log.Println("create_user SendMessage error:", err)
// 	}
// }
