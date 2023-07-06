package server

import (
	"auth-proxy/pkg/app"
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/authenticator"
	"auth-proxy/pkg/counter"
	"strconv"

	"auth-proxy/pkg/db"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func login() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.String,
		Description: "Войти по имени или email и паролю",
		Args: graphql.FieldConfigArgument{
			"username": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя пользователя",
			},
			"password": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Пароль",
			},
			"captcha": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Captcha",
			},
			"pin": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "PIN by Google Authenticator",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ArgToLowerCase(params, "username")
			c, _ := params.Context.Value("ginContext").(*gin.Context)
			username, _ := params.Args["username"].(string)
			password, _ := params.Args["password"].(string)
			captcha, _ := params.Args["captcha"].(string)

			// небольшая задержка чтобы усложнить перебор паролей
			time.Sleep(500 * time.Millisecond)

			// проверить PIN если установлен глобальный флаг
			if app.Params.UsePin {
				// нужен ли PIN для данного пользователя?
				fmt.Println("Auth-proxy uses PIN !!!!!!")
				pinRequired, _, _, err := authenticator.GetUserPinFields(username)
				if err != nil {
					return nil, errors.New("email или пароль введен неверно")
				}
				if pinRequired {
					fmt.Printf("PIN is required for user %s !!! \n", username)
					pin, _ := params.Args["pin"].(string)
					fmt.Println("pin = ", pin)
					// правильный ли пин?
					err := authenticator.IsPinGood(username, pin, false)
					if err != nil {
						fmt.Println("pin is bad!!!!!!")
						return "", err
					}
					fmt.Println("pin is good!!!!!!")
				}
			}

			// проверить капчу если превышено число допустимых ошибок входа /* или это админ */
			if UseCaptcha && counter.IsTooBig(username) /* || auth.AppUserRoleExist("auth", username, "authadmin")*/ {
				if captcha == "" {
					return "", errors.New("Вы должны ввести картинку. uri=/captcha ")
				}

				sessionCaptcha := GetSessionVariable(c, "captcha")
				if sessionCaptcha != captcha {
					return "", errors.New("Картинка введена с ошибкой")
				}
			}

			r, dbUsername := auth.CheckUserPassword2(username, password)
			// Пользователь не найден
			if r == auth.NO_USER {
				return nil, errors.New("email или пароль введен неверно")
			} else
			// Пользователь найден, но пароль неверный
			if r == auth.WRONG_PASSWORD {
				counter.IncrementCounter(username)
				return nil, errors.New("email или пароль введен неверно")
			} else
			// Пользователь найден, но неактивирован
			if r == auth.USER_DISABLED {
				return nil, errors.New(username + " деактивирован.")
			} else
			// Пользователь найден, но не подтвержден email
			if r == auth.EMAIL_NOT_CONFIRMED {
				// Если включена проверка подтверждения email, то
				if !app.Params.LoginNotConfirmedEmail {
					// получаем имя пользователя по email на случай, если введен email
					username1 := auth.GetUserNameByEmail(username)
					// проверяем, является ли пользователь админом
					isAdmin := auth.AppUserRoleExist("auth", username, "authadmin") || auth.AppUserRoleExist("auth", username1, "authadmin")
					fmt.Printf("User %s is admin: %v \n", username, isAdmin)
					// Если пользователь не админ, то возвращаем ошибку и посылаем письмо для подтверждения email
					if !isAdmin {
						UpdateHashAndSendEmail(username, username, password, false)
						return nil, errors.New("email not confirmed")
					}
				}
			}

			// Все проверки пройдены. Устанавливаем переменные сессии
			_ = SetSessionVariable(c, "user", dbUsername)

			user, err := auth.GetUser(username)
			if err != nil {
				return nil, errors.New("cant get record for " + username)
			}
			id := strconv.FormatInt(user["id"].(int64), 10)
			_ = SetSessionVariable(c, "id", id)
			email := user["email"].(string)
			_ = SetSessionVariable(c, "email", email)

			counter.ResetCounter(username)
			return "Success. " + dbUsername + " is authenticated.", nil
		},
	}
}

func login_by_email() *graphql.Field {
	return &graphql.Field{
		Type:        userObject,
		Description: "Войти по email и паролю",
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Email пользователя",
			},
			"password": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Пароль",
			},
			"captcha": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Captcha",
			},
			"pin": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "PIN by Google Authenticator",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ArgToLowerCase(params, "username")
			c, _ := params.Context.Value("ginContext").(*gin.Context)
			username, _ := params.Args["email"].(string)
			password, _ := params.Args["password"].(string)
			captcha, _ := params.Args["captcha"].(string)

			// небольшая задержка чтобы усложнить перебор паролей
			time.Sleep(500 * time.Millisecond)

			// проверить PIN если установлен глобальный флаг
			if app.Params.UsePin {
				// нужен ли PIN для данного пользователя?
				fmt.Println("UsePin!!!!!!")
				pinRequired, _, _, err := authenticator.GetUserPinFields(username)
				if err != nil {
					return "", err
				}
				fmt.Println("pin required!!!!!!")
				if pinRequired {
					pin, _ := params.Args["pin"].(string)
					fmt.Println("pin = ", pin)
					// правильный ли пин?
					err := authenticator.IsPinGood(username, pin, false)
					if err != nil {
						fmt.Println("pin is bad!!!!!!")
						return "", err
					}
					fmt.Println("pin is good!!!!!!")
				}
			}

			// проверить капчу если превышено число допустимых ошибок входа /* или это админ */
			if UseCaptcha && counter.IsTooBig(username) /* || auth.AppUserRoleExist("auth", username, "authadmin")*/ {
				if captcha == "" {
					return "", errors.New("Вы должны ввести картинку. uri=/captcha ")
				}

				sessionCaptcha := GetSessionVariable(c, "captcha")
				if sessionCaptcha != captcha {
					return "", errors.New("Картинка введена с ошибкой")
				}
			}

			r, dbUsername := auth.CheckUserPassword2(username, password)
			if r == auth.NO_USER {
				// return nil, errors.New(username + " не зарегистрирован или БД не доступна. Обратитесь к администратору.")
				return nil, errors.New("email или пароль введен неверно")
			} else if r == auth.WRONG_PASSWORD {
				counter.IncrementCounter(username)
				return nil, errors.New("email или пароль введен неверно")
			} else if r == auth.USER_DISABLED {
				return nil, errors.New(username + " деактивирован.")
			} else if r == auth.EMAIL_NOT_CONFIRMED {
				if !app.Params.LoginNotConfirmedEmail {
					UpdateHashAndSendEmail(username, username, password, false)
					return nil, errors.New("email not confirmed")
				}
			}

			// Все проверки пройдены. Устанавливаем переменные сессии
			_ = SetSessionVariable(c, "user", dbUsername)

			user, err := auth.GetUser(username)
			if err != nil {
				return nil, errors.New("cant get record for " + username)
			}
			id := strconv.FormatInt(user["id"].(int64), 10)
			_ = SetSessionVariable(c, "id", id)
			email := user["email"].(string)
			_ = SetSessionVariable(c, "email", email)

			counter.ResetCounter(username)

			fields := getSelectedFields([]string{"login_by_email"}, params)
			res, err := db.QueryRowMap("SELECT "+fields+` FROM "user" WHERE email = $1 ;`, email)
			return res, err

		},
	}
}

func logout() *graphql.Field {
	return &graphql.Field{
		Type:        authMessageObject,
		Description: "Выйти",
		Args:        graphql.FieldConfigArgument{},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			c, _ := params.Context.Value("ginContext").(*gin.Context)
			email := GetSessionVariable(c, "email")
			DeleteSession(c)
			return gin.H{"email": email, "message": "Successfully logged out"}, nil
		},
	}
}

func is_selfreg_allowed() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Возможна ли саморегистрация пользователей",
		Args:        graphql.FieldConfigArgument{},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return SelfRegistrationAllowed, nil
		},
	}
}

func get_stat() *graphql.Field {
	return &graphql.Field{
		Type:        statObject,
		Description: "records statistics about the memory allocator in megabytes and requests",
		Args:        graphql.FieldConfigArgument{},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return app.GetStat(), nil
		},
	}
}

func list_oauth_providers() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(oauthProviderObject),
		Description: "Показать список провайдеров Oauth2",
		Args:        graphql.FieldConfigArgument{},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			results := []map[string]string{}

			for provider := range Oauth2Params {
				rec := make(map[string]string)
				rec["provider_name"] = provider
				rec["login_endpoint"] = "/oauthlogin/" + provider
				rec["logout_endpoint"] = "/oauthlogout/" + provider
				results = append(results, rec)
			}

			return results, nil
		},
	}
}

func get_params() *graphql.Field {
	return &graphql.Field{
		Type:        appParamsObject,
		Description: "Показать параметры приложения",
		Args:        graphql.FieldConfigArgument{},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return app.GetParams(), nil
		},
	}
}

func set_params() *graphql.Field {
	return &graphql.Field{
		Type:        appParamsObject,
		Description: "Установить параметры приложения",
		Args: graphql.FieldConfigArgument{

			"selfreg": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "Могут ли пользователи регистрироваться самостоятельно",
				// DefaultValue: false,
			},
			"use_captcha": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "Нужно ли вводить капчу при входе в систему",
				// DefaultValue: true,
			},
			"use_pin": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "Нужно ли вводить PIN при входе в систему",
				// DefaultValue: false,
			},
			"login_not_confirmed_email": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "Разрешить авторизацию пользователей не подтвердивших email",
				// DefaultValue: true,
			},
			"no_schema": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "подавлять чтение схемы GraphQL",
				// DefaultValue: false,
			},
			"max_attempts": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Максимально допустимое число ошибок ввода пароля",
				// DefaultValue: 5,
			},
			"reset_time": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Время сброса счетчика ошибок пароля в минутах",
				// DefaultValue: 60,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if !isAuthAdmin(params) {
				return nil, errors.New("sorry. you have no admin rights")
			}

			// перед установкой параметров проверяем, что они переданы
			if _, ok := params.Args["selfreg"]; ok {
				app.Params.Selfreg = params.Args["selfreg"].(bool)
			}
			if _, ok := params.Args["use_captcha"]; ok {
				app.Params.UseCaptcha = params.Args["use_captcha"].(bool)
			}
			if _, ok := params.Args["use_pin"]; ok {
				app.Params.UsePin = params.Args["use_pin"].(bool)
			}
			if _, ok := params.Args["login_not_confirmed_email"]; ok {
				app.Params.LoginNotConfirmedEmail = params.Args["login_not_confirmed_email"].(bool)
			}
			if _, ok := params.Args["no_schema"]; ok {
				app.Params.NoSchema = params.Args["no_schema"].(bool)
				if app.Params.NoSchema {
					disableSchemaIntrospection()
				} else {
					enableSchemaIntrospection()
				}
			}
			if _, ok := params.Args["max_attempts"]; ok {
				app.Params.MaxAttempts = int64(params.Args["max_attempts"].(int))
			}
			if _, ok := params.Args["reset_time"]; ok {
				app.Params.ResetTime = int64(params.Args["reset_time"].(int))
			}

			// TODO: удалить эти глобальные переменные
			SelfRegistrationAllowed = app.Params.Selfreg
			UseCaptcha = app.Params.UseCaptcha

			counter.MAX_ATTEMPTS = app.Params.MaxAttempts
			counter.RESET_TIME = time.Duration(app.Params.ResetTime)

			return app.GetParams(), nil
		},
	}
}

func is_captcha_required() *graphql.Field {
	return &graphql.Field{
		Type:        isCaptchaRequiredObject,
		Description: "Нужно ли пользователю вводить каптчу",
		Args: graphql.FieldConfigArgument{
			"username": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Идентификатор пользователя",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if !UseCaptcha {
				return gin.H{"is_required": false, "path": ""}, nil
			}
			ArgToLowerCase(params, "username")
			username := params.Args["username"].(string)
			if counter.IsTooBig(username) /* || auth.AppUserRoleExist("auth", username, "authadmin")*/ {
				return gin.H{"is_required": true, "path": "/captcha"}, nil
			}
			return gin.H{"is_required": false, "path": ""}, nil
		},
	}
}

func is_pin_required() *graphql.Field {
	return &graphql.Field{
		Type:        isPinRequiredObject,
		Description: "Информация о необходимости вводить PIN для входа в систему",
		Args: graphql.FieldConfigArgument{
			"username": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Идентификатор пользователя",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			pinRequired, _, _, err := authenticator.GetUserPinFields(params.Args["username"].(string))
			if err != nil {
				return nil, err
			}

			return gin.H{"use_pin": app.Params.UsePin, "pinrequired": pinRequired}, err
		},
	}
}

func get_user() *graphql.Field {
	return &graphql.Field{
		Type:        userObject,
		Description: "Показать пользователя",
		Args: graphql.FieldConfigArgument{
			"username": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя пользователя или email",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ArgToLowerCase(params, "username")
			panicIfNotOwnerOrAdminOrAuditor(params)
			fields := getSelectedFields([]string{"get_user"}, params)
			return db.QueryRowMap("SELECT "+fields+` FROM "user" WHERE username = $1 OR email = $1 ;`, params.Args["username"])
		},
	}
}

func get_user_by_id() *graphql.Field {
	return &graphql.Field{
		Type:        userObject,
		Description: "Показать пользователя",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "id пользователя",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ArgToLowerCase(params, "username")
			panicIfNotOwnerOrAdminOrAuditor(params)
			fields := getSelectedFields([]string{"get_user_by_id"}, params)
			return db.QueryRowMap("SELECT "+fields+` FROM "user" WHERE id = $1 ;`, params.Args["id"])
		},
	}
}

func get_logined_user() *graphql.Field {
	return &graphql.Field{
		Type:        userObject,
		Description: "Показать пользователя",
		Args:        graphql.FieldConfigArgument{},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotUser(params)
			email := getLoginedUserEmail(params)
			fields := getSelectedFields([]string{"get_logined_user"}, params)
			res, err := db.QueryRowMap("SELECT "+fields+` FROM "user" WHERE email = $1 ;`, email)
			return res, err
		},
	}
}

func get_group() *graphql.Field {
	return &graphql.Field{
		Type:        groupObject,
		Description: "Показать группу по id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "id группы",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdminOrAuditor(params)
			fields := getSelectedFields([]string{"get_group"}, params)
			return db.QueryRowMap("SELECT "+fields+" FROM \"group\" WHERE id = $1 ;", params.Args["id"])
		},
	}
}

func get_group_by_name() *graphql.Field {
	return &graphql.Field{
		Type:        groupObject,
		Description: "Показать группу по имени",
		Args: graphql.FieldConfigArgument{
			"groupname": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя группы",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdminOrAuditor(params)
			fields := getSelectedFields([]string{"get_group_by_name"}, params)
			return db.QueryRowMap("SELECT "+fields+" FROM \"group\" WHERE groupname = $1 ;", params.Args["groupname"])
		},
	}
}

func list_group() *graphql.Field {
	return &graphql.Field{
		Type:        listGroupGQType,
		Description: "Получить список групп.",
		Args: graphql.FieldConfigArgument{
			"search": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Строка полнотекстового поиска.",
			},
			"order": &graphql.ArgumentConfig{
				Type:         graphql.String,
				Description:  "сортировка строк в определённом порядке. По умолчанию 'appname ASC'",
				DefaultValue: "groupname ASC",
			},
			"limit": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				Description:  "возвратить не больше заданного числа строк. По умолчанию 100.",
				DefaultValue: 1000,
			},
			"offset": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				Description:  "пропустить указанное число строк, прежде чем начать выдавать строки. По умолчанию 0.",
				DefaultValue: 0,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdminOrAuditor(params)

			wherePart, orderAndLimits := QueryEnd(params, "groupname,description")
			fields := getSelectedFields([]string{"list_group", "list"}, params)
			fmt.Printf("wherePart:%v, orderAndLimits:%v", wherePart, orderAndLimits)
			list, err := db.QuerySliceMap("SELECT " + fields + ` FROM "group"` + wherePart + orderAndLimits)
			if err != nil {
				return nil, err
			}

			count, err := db.QueryRowMap(`SELECT count(*) AS count FROM "group"` + wherePart)
			if err != nil {
				return nil, err
			}

			m := map[string]interface{}{
				"length": count["count"],
				"list":   list,
			}

			return m, nil

		},
	}
}

func list_group_user_role() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(groupUserRoleObject),
		Description: "Получить список групп пользователей и их ролей.",
		Args: graphql.FieldConfigArgument{
			"group_id": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Идентификатор группы",
			},
			"user_id": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Идентификатор пользователя",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotOwnerOrAdminOrAuditor(params)
			fields := getSelectedFields([]string{"list_group_user_role"}, params)
			wherePart := listGroupUserRoleWherePart(params)
			query := fmt.Sprintf(`SELECT DISTINCT %s FROM group_user_role_extended %s `, fields, wherePart)
			return db.QuerySliceMap(query)
		},
	}
}

func listGroupUserRoleWherePart(params graphql.ResolveParams) (wherePart string) {
	var searchConditions []string
	for paramName, v := range params.Args {
		n, ok := v.(int)
		if ok {
			searchConditions = append(searchConditions, fmt.Sprintf(" %s = %d ", paramName, n))
		}
	}
	if len(searchConditions) > 0 {
		wherePart = " WHERE " + strings.Join(searchConditions, " AND ")
	}
	return wherePart
}

func list_group_app_role() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(groupAppRoleObject),
		Description: "Получить список приложений групп и их ролей.",
		Args: graphql.FieldConfigArgument{
			"group_id": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Идентификатор группы",
			},
			"app_id": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Идентификатор приложения",
			},
			"rolename": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Идентификатор роли",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotOwnerOrAdminOrAuditor(params)
			fields := getSelectedFields([]string{"list_group_app_role"}, params)
			wherePart := listGroupAppRoleWherePart(params)
			query := fmt.Sprintf(`SELECT DISTINCT %s FROM group_app_role_extended %s `, fields, wherePart)
			return db.QuerySliceMap(query)
		},
	}
}

func listGroupAppRoleWherePart(params graphql.ResolveParams) (wherePart string) {
	var searchConditions []string
	for paramName, v := range params.Args {
		// if v is a number
		n, ok := v.(int)
		if ok {
			searchConditions = append(searchConditions, fmt.Sprintf(" %s = %d ", paramName, n))
		} else {
			// if v is a string
			s, ok := v.(string)
			s = strings.Trim(s, " ")
			if ok && len(s) > 0 {
				searchConditions = append(searchConditions, fmt.Sprintf(" %s = '%s' ", paramName, db.RemoveSingleQuotes(s)))
			}
		}

	}
	if len(searchConditions) > 0 {
		wherePart = " WHERE " + strings.Join(searchConditions, " AND ")
	}
	return wherePart
}

func get_app() *graphql.Field {
	return &graphql.Field{
		Type:        appObject,
		Description: "Показать приложение",
		Args: graphql.FieldConfigArgument{
			"appname": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Имя приложения",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdminOrAuditor(params)
			fields := getSelectedFields([]string{"get_app"}, params)
			return db.QueryRowMap("SELECT "+fields+" FROM app WHERE appname = $1 ;", params.Args["appname"])
		},
	}
}

func list_user() *graphql.Field {
	return &graphql.Field{
		Type:        listUserGQType,
		Description: "Получить список пользователей.",
		Args: graphql.FieldConfigArgument{
			"search": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Строка полнотекстового поиска.",
			},
			"order": &graphql.ArgumentConfig{
				Type:         graphql.String,
				Description:  "сортировка строк в определённом порядке. По умолчанию 'username ASC'",
				DefaultValue: "username ASC",
			},
			"limit": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				Description:  "возвратить не больше заданного числа строк. По умолчанию 100.",
				DefaultValue: 1000,
			},
			"offset": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				Description:  "пропустить указанное число строк, прежде чем начать выдавать строки. По умолчанию 0.",
				DefaultValue: 0,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdminOrAuditor(params)
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
	}
}

func list_user_by_ids() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(userObject),
		Description: "Получить список пользователей по их id.",
		Args: graphql.FieldConfigArgument{
			"ids": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.NewList(graphql.Int)),
				Description: "Массив идентификаторов id",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdminOrAuditor(params)
			fields := getSelectedFields([]string{"list_user_by_ids"}, params)
			ids, _ := params.Args["ids"]
			idss, _ := db.SerializeIfArray(ids).(string)
			idss = strings.Trim(idss, "[]")
			if idss == "" {
				return []interface{}{}, nil
			}
			query := fmt.Sprintf(`SELECT %s FROM "user" WHERE id IN (%s)`, fields, idss)
			return db.QuerySliceMap(query)
		},
	}
}

func list_user_by_usernames() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(userObject),
		Description: "Получить список пользователей по их username.",
		Args: graphql.FieldConfigArgument{
			"usernames": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.NewList(graphql.String)),
				Description: "Массив идентификаторов username",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotAdminOrAuditor(params)
			fields := getSelectedFields([]string{"list_user_by_usernames"}, params)
			ids, _ := params.Args["usernames"]
			idss, _ := db.SerializeIfArray(ids).(string)
			idss = db.RemoveSingleQuotes(idss)
			idss = strings.Trim(idss, "[]")
			if idss == "" {
				return []interface{}{}, nil
			}
			idss = strings.ReplaceAll(idss, "\"", "'")
			query := fmt.Sprintf(`SELECT %s FROM "user" WHERE username IN (%s)`, fields, idss)
			return db.QuerySliceMap(query)
		},
	}
}

func list_app() *graphql.Field {
	return &graphql.Field{
		Type:        listAppGQType,
		Description: "Получить список приложений.",
		Args: graphql.FieldConfigArgument{
			"search": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Строка полнотекстового поиска.",
			},
			"order": &graphql.ArgumentConfig{
				Type:         graphql.String,
				Description:  "сортировка строк в определённом порядке. По умолчанию 'appname ASC'",
				DefaultValue: "appname ASC",
			},
			"limit": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				Description:  "возвратить не больше заданного числа строк. По умолчанию 100.",
				DefaultValue: 1000,
			},
			"offset": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				Description:  "пропустить указанное число строк, прежде чем начать выдавать строки. По умолчанию 0.",
				DefaultValue: 0,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
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
	}
}

func list_app_user_role() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(appUserRoleExtendedObject),
		Description: "Получить список приложений пользователей и их ролей.",
		Args: graphql.FieldConfigArgument{
			"appname": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Идентификатор приложения",
			},
			"username": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Идентификатор пользователя",
			},
			"rolename": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Идентификатор роли",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			panicIfNotOwnerOrAdminOrAuditor(params)
			fields := getSelectedFields([]string{"list_app_user_role"}, params)
			wherePart := listAppUserRoleWherePart(params)
			query := fmt.Sprintf(`SELECT DISTINCT %s FROM app_user_role_extended %s `, fields, wherePart)
			return db.QuerySliceMap(query)
		},
	}
}

func listAppUserRoleWherePart(params graphql.ResolveParams) (wherePart string) {
	var searchConditions []string
	for paramName, v := range params.Args {
		s, ok := v.(string)
		s = strings.Trim(s, " ")
		if ok && len(s) > 0 {
			searchConditions = append(searchConditions, fmt.Sprintf(" %s = '%s' ", paramName, db.RemoveSingleQuotes(s)))
		}
	}
	if len(searchConditions) > 0 {
		wherePart = " WHERE " + strings.Join(searchConditions, " AND ")
	}
	return wherePart
}

func QueryEnd(params graphql.ResolveParams, fieldList string) (wherePart string, orderAndLimits string) {
	var searchConditions []string

	// Если запрос к таблице app и это не административный или аудиторский запрос
	// ограничиваем выборку
	isAdm := (isAuthAdmin(params) || isAuditor(params))
	isAppQuery := strings.Contains(fieldList, "appname")
	if isAppQuery && !isAdm {
		userName := getLoginedUserName(params)
		inClause := getListOfAllowedAppNames(userName)
		fmt.Println("inClause=", inClause)
		searchConditions = append(searchConditions, inClause)
	}

	search, ok := params.Args["search"].(string)
	search = strings.Trim(search, " ")
	search = db.RemoveSingleQuotes(search)
	if ok && len(search) > 0 {
		search = strings.ReplaceAll(search, " ", "%")
		searchConditions = append(searchConditions, Like(fieldList, search))
	}

	if len(searchConditions) > 0 {
		wherePart = " WHERE " + strings.Join(searchConditions, " AND ")
	}
	orderClause := db.SanitizeOrderClause(params.Args["order"].(string))
	orderAndLimits = fmt.Sprintf(" ORDER BY %s LIMIT %v OFFSET %v ;", orderClause, params.Args["limit"], params.Args["offset"])
	if orderClause == "" {
		orderAndLimits = fmt.Sprintf(" LIMIT %v OFFSET %v ;", params.Args["limit"], params.Args["offset"])
	}
	return wherePart, orderAndLimits
}

func getListOfAllowedAppNames(userName string) string {
	slice, err := db.QuerySliceMap(`
		SELECT DISTINCT appname
		FROM app_user_role_union
		WHERE username = $1
		UNION
		SELECT appname FROM app WHERE public = 'Y';	
	`, userName)
	if err != nil {
		return " TRUE "
	}
	appList := make([]string, 0)
	for _, element := range slice {
		appName, _ := element["appname"].(string)
		appList = append(appList, appName)
	}
	s := "('" + strings.Join(appList, "', '") + "')"
	return " appname IN " + s + " "
}

func Like(fieldsString, search string) string {
	fields := strings.Split(fieldsString, ",")
	var chunks []string
	for _, field := range fields {
		chunks = append(chunks, ` LOWER(`+field+`) LIKE LOWER('%`+search+`%') `)
	}
	s := strings.Join(chunks, " OR ")
	// fmt.Println(fieldsString, s)
	return s
}
