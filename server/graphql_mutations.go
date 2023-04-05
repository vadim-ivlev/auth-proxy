package server

import (
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/db"
	"auth-proxy/pkg/mail"
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func create_user() *graphql.Field {
	return &graphql.Field{
		Description: "Создать пользователя",
		Type:        userObject,
		Args: graphql.FieldConfigArgument{
			"password": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Пароль",
			},
			"email": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Email пользователя",
			},
			"fullname": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
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
			"noemail": &graphql.ArgumentConfig{
				Type:         graphql.Boolean,
				Description:  "Не посылать пользователю письмо о регистрации.",
				DefaultValue: false,
			},
			"emailconfirmed": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "подтвержеден ли email пользователя",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if SelfRegistrationAllowed || isAuthAdmin(params) {
				panicIfEmpty(params.Args["password"], "Введите пароль")
				panicIfEmpty(params.Args["email"], "Заполните поле Email")
				panicIfEmpty(params.Args["fullname"], "Заполните имя")
				// Только админ может включить/отключить проверку пина и установить emailconfirmed
				if !isAuthAdmin(params) {
					delete(params.Args, "pinrequired")
					delete(params.Args, "emailconfirmed")
				}
				noemail, _ := params.Args["noemail"].(bool)
				delete(params.Args, "noemail")

				ArgToLowerCase(params, "email")
				TrimParamValue(params, "email")
				params.Args["username"] = params.Args["email"]
				convertPasswordToHash(params)
				clearCache()
				res, err := createRecord("username", params, "user", "user")
				if err == nil {
					// Отправляем письмо пользователю
					if !noemail {
						UpdateHashAndSendEmail(params.Args["email"].(string), params.Args["fullname"].(string))
					}
					// Добавляем пользователя в группу по умолчанию
					userID, ok := res.(map[string]interface{})["id"].(int64)
					if !ok {
						log.Println("create_user addUserToGroup error: can't get user id")
					}
					err = addUserToDefaultGroup(userID)
					if err != nil {
						log.Println("create_user addUserToGroup error:", err)
					}
				}
				return res, err
			}
			return nil, errors.New("self registration is not allowed. Please ask administrators")
		},
	}
}

func send_confirm_email() *graphql.Field {
	return &graphql.Field{
		Description: "Послать пользователю письмо для подтверждения регистрации",
		Type:        userObject,
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Email пользователя",
			},
			"password": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Пароль",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			panicIfEmpty(params.Args["email"], "Заполните поле Email")
			panicIfEmpty(params.Args["password"], "Введите пароль")

			ArgToLowerCase(params, "email")
			TrimParamValue(params, "email")

			email, _ := params.Args["email"].(string)
			password, _ := params.Args["password"].(string)

			// проверить пароль
			r, dbUsername := auth.CheckUserPassword2(email, password)
			fmt.Println("r=", r, "dbUsername=", dbUsername)
			if r == auth.NO_USER {
				return nil, errors.New("email или пароль введен неверно")
			} else if r == auth.WRONG_PASSWORD {
				// counter.IncrementCounter(email)
				return nil, errors.New("email или пароль введен неверно")
			} else if r == auth.USER_DISABLED {
				return nil, errors.New(email + " деактивирован.")
			}

			return UpdateHashAndSendEmail(email, dbUsername)

		},
	}
}

func send_email() *graphql.Field {
	return &graphql.Field{
		Description: "Послать пользователю email с произвольным текстом",
		Type:        graphql.String,
		Args: graphql.FieldConfigArgument{
			"email_from": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Email от кого отправляется письмо",
			},
			"email_to": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Email кому отправляется письмо",
			},
			"subject": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "тема письма",
			},
			"text": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "текст письма",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if !isAuthAdmin(params) {
				return nil, errors.New("access denied")
			}
			panicIfEmpty(params.Args["email_from"], "Заполните поле email_from")
			panicIfEmpty(params.Args["email_to"], "Заполните поле email_to")

			email_from, _ := params.Args["email_from"].(string)
			email_to, _ := params.Args["email_to"].(string)
			subject, _ := params.Args["subject"].(string)
			text, _ := params.Args["text"].(string)

			err := mail.SendEmailTo(email_from, email_to, subject, text)
			if err != nil {
				return nil, err
			}

			return fmt.Sprintf("the letter to %s is sent", email_to), nil
		},
	}
}

func update_user() *graphql.Field {
	return &graphql.Field{
		Description: "Обновить пользователя",
		Type:        userObject,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор пользователя",
			},
			"password": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Пароль",
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
			"emailconfirmed": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "подтвержеден ли email пользователя",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotOwnerOrAdmin(params)
			// Только админ может включить/отключить проверку пина и установить emailconfirmed
			if !isAuthAdmin(params) {
				delete(params.Args, "pinrequired")
				delete(params.Args, "emailconfirmed")
			}
			convertPasswordToHash(params)
			clearCache()
			id := params.Args["id"].(int)
			// return db.UpdateRowByID("id", "user", id, params.Args)
			return updateRecord(id, "id", params, "user", "user")
		},
	}
}

func delete_user() *graphql.Field {
	return &graphql.Field{
		Description: "Удалить пользователя",
		Type:        userObject,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор пользователя",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotOwnerOrAdmin(params)
			clearCache()

			// Получаем данные удаляемого пользователя
			id := params.Args["id"]
			deletingRecord, _ := db.QueryRowMap(`SELECT * FROM "user" WHERE id = $1`, id)
			deletingRecord["deleted_by"] = getLoginedUserEmail(params)

			// Удаляем пользователя
			res, err := deleteRecord("id", params, "user", "user")

			if err == nil {
				// Записываем в таблицу удаленных пользователей
				_, createErr := db.CreateRow("user_deleted", deletingRecord)
				if createErr != nil {
					log.Println(`db.CreateRow("user_deleted") error:`, createErr)
				}

				// Если пользователь удаляет сам себя, то разлогиниваем его
				deletedUserID := fmt.Sprintf("%v", params.Args["id"])
				loginedUserID := getLoginedUserID(params)
				if loginedUserID == deletedUserID {
					c, _ := params.Context.Value("ginContext").(*gin.Context)
					DeleteSession(c)
					log.Printf("USER %v LOGGED OUT \n", deletedUserID)
				}
			}
			return res, err
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
			TrimParamValue(params, "appname")
			res, err := createRecord("appname", params, "app", "app")
			if err == nil {
				app, _ := params.Args["appname"].(string)
				url, _ := params.Args["url"].(string)
				rebase, _ := params.Args["rebase"].(string)
				if url != "" {
					proxies[app] = createProxy(url, app, rebase)
					log.Printf("Proxy created appname=%v target=%v rebase=%v", app, url, rebase)
				}
				clearCache()
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
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор приложения",
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

			appID := params.Args["id"].(int)
			app, err := auth.GetApp(appID)
			if err != nil {
				return nil, err
			}
			oldAppname := app["appname"].(string)
			// oldAppname, _ := params.Args["old_appname"].(string)
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
				TrimParamValue(params, "appname")
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
				clearCache()
			}
			return res, err
		},
	}
}

func delete_app() *graphql.Field {
	return &graphql.Field{
		Description: "Удалить приложение",
		Type:        appObject,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор приложения",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdmin(params)

			appID := params.Args["id"].(int)
			app, _ := auth.GetApp(appID)
			appName, _ := app["appname"].(string)
			delete(proxies, appName)
			log.Printf("Proxy deleted appname=%s", appName)
			clearCache()

			res, err := deleteRecord("id", params, "app", "app")
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
			TrimParamValue(params, "username")
			TrimParamValue(params, "appname")
			TrimParamValue(params, "rolename")
			clearCache()
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
			clearCache()
			return map[string]interface{}{"appname": a, "username": u, "rolename": r}, nil
		},
	}
}

func clearCache() {
	auth.Cache.Flush()
}

// ----------------------------------------------------------------
func create_group() *graphql.Field {
	return &graphql.Field{
		Description: "Создать группу",
		Type:        groupObject,
		Args: graphql.FieldConfigArgument{
			"groupname": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя группы (уникальное)",
			},
			"description": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Описание",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdmin(params)
			panicIfEmpty(params.Args["groupname"], "Имя группы не должно быть пустым")
			TrimParamValue(params, "groupname")
			res, err := createRecord("groupname", params, "group", "group")
			if err == nil {
				clearCache()
			}
			return res, err
		},
	}
}

func update_group() *graphql.Field {
	return &graphql.Field{
		Description: "Обновить группу",
		Type:        groupObject,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор группы",
			},
			"groupname": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Имя группы (уникальное)",
			},
			"description": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Описание",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdmin(params)
			TrimParamValue(params, "groupname")
			id, _ := params.Args["id"].(int)
			res, err := updateRecord(id, "id", params, "group", "group")
			if err == nil {
				clearCache()
			}
			return res, err
		},
	}
}

func delete_group() *graphql.Field {
	return &graphql.Field{
		Description: "Удалить группу",
		Type:        groupObject,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор группы",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdmin(params)
			res, err := deleteRecord("id", params, "group", "group")
			if err == nil {
				clearCache()
			}
			return res, err
		},
	}
}

func create_group_app_role() *graphql.Field {
	return &graphql.Field{
		Description: "Создать роль группы для приложения",
		Type:        groupAppRoleObject,
		// Type: insertResultObject,
		Args: graphql.FieldConfigArgument{
			"group_id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор группы",
			},
			"app_id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор приложения",
			},
			"rolename": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя роли",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdmin(params)
			TrimParamValue(params, "rolename")
			_, err := db.CreateRow("group_app_role", params.Args)
			if err != nil {
				return nil, err
			}
			clearCache()
			return map[string]interface{}{
				"app_id":   params.Args["app_id"],
				"group_id": params.Args["group_id"],
				"rolename": params.Args["rolename"],
			}, nil

		},
	}
}

func delete_group_app_role() *graphql.Field {
	return &graphql.Field{
		Description: "Удалить роль пользователя для приложения",
		Type:        groupAppRoleObject,
		Args: graphql.FieldConfigArgument{
			"group_id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор группы",
			},
			"app_id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор приложения",
			},
			"rolename": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя роли",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdmin(params)
			app_id, group_id, rolename := params.Args["app_id"], params.Args["group_id"], params.Args["rolename"]

			res, err := db.QueryExec(
				`DELETE FROM group_app_role WHERE app_id = $1 AND group_id = $2 AND rolename = $3 ;`, app_id, group_id, rolename)
			if err != nil {
				return nil, err
			}
			rowsAffected, _ := res.RowsAffected()
			if rowsAffected == 0 {
				return nil, errors.New("no records to delete")
			}
			clearCache()
			return map[string]interface{}{"app_id": app_id, "group_id": group_id, "rolename": rolename}, nil
		},
	}
}

func create_group_user_role() *graphql.Field {
	return &graphql.Field{
		Description: "Добавить пользователя в группу",
		Type:        groupUserRoleObject,
		Args: graphql.FieldConfigArgument{
			"group_id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор группы",
			},
			"user_id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор пользователя",
			},
			"rolename": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Имя роли",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdmin(params)
			TrimParamValue(params, "rolename")
			_, err := db.CreateRow("group_user_role", params.Args)
			if err != nil {
				return nil, err
			}
			clearCache()
			return map[string]interface{}{
				"user_id":  params.Args["user_id"],
				"group_id": params.Args["group_id"],
				"rolename": params.Args["rolename"],
			}, nil

		},
	}
}

func delete_group_user_role() *graphql.Field {
	return &graphql.Field{
		Description: "Удалить пользователя из группы",
		Type:        groupUserRoleObject,
		Args: graphql.FieldConfigArgument{
			"group_id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор группы",
			},
			"user_id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Идентификатор пользователя",
			},
			"rolename": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя роли",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdmin(params)
			user_id, group_id, rolename := params.Args["user_id"], params.Args["group_id"], params.Args["rolename"]
			res, err := db.QueryExec(`DELETE FROM group_user_role WHERE user_id = $1 AND group_id = $2 AND rolename = $3;`, user_id, group_id, rolename)
			if err != nil {
				return nil, err
			}
			rowsAffected, _ := res.RowsAffected()
			if rowsAffected == 0 {
				return nil, errors.New("no records to delete")
			}
			clearCache()
			return map[string]interface{}{"user_id": user_id, "group_id": group_id, "rolename": rolename}, nil
		},
	}
}
