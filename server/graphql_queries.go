package server

import (
	"auth-proxy/pkg/app"
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/authenticator"
	"auth-proxy/pkg/counter"

	"auth-proxy/pkg/db"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	gq "github.com/graphql-go/graphql"
)

// ************************************************************************

var rootQuery = gq.NewObject(gq.ObjectConfig{
	Name: "Query",
	Fields: gq.Fields{

		"login": &gq.Field{
			Type:        gq.String,
			Description: "Войти по имени или email и паролю",
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя",
				},
				"password": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Пароль",
				},
				"captcha": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Captcha",
				},
				"pin": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "PIN by Google Authenticator",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
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
						err := authenticator.IsPinGood(username, pin)
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
					return nil, errors.New(username + " не зарегистрирован или БД не доступна. Обратитесь к администратору.")
				} else if r == auth.WRONG_PASSWORD {
					counter.IncrementCounter(username)
					return nil, errors.New("Неверный пароль")
				} else if r == auth.USER_DISABLED {
					return nil, errors.New(username + " деактивирован.")
				}
				counter.ResetCounter(username)
				_ = SetSessionVariable(c, "user", dbUsername)
				return "Success. " + dbUsername + " is authenticated.", nil
			},
		},

		"logout": &gq.Field{
			Type:        authMessageObject,
			Description: "Выйти",
			Args:        gq.FieldConfigArgument{},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				c, _ := params.Context.Value("ginContext").(*gin.Context)
				username := GetSessionVariable(c, "user")
				DeleteSession(c)
				return gin.H{"username": username, "message": "Successfully logged out"}, nil
			},
		},

		"is_selfreg_allowed": &gq.Field{
			Type:        gq.Boolean,
			Description: "Возможна ли саморегистрация пользователей",
			Args:        gq.FieldConfigArgument{},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				return SelfRegistrationAllowed, nil
			},
		},

		"get_stat": &gq.Field{
			Type:        statObject,
			Description: "records statistics about the memory allocator in megabytes and requests",
			Args:        gq.FieldConfigArgument{},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				return app.GetStat(), nil
			},
		},

		"list_oauth_providers": &gq.Field{
			Type:        gq.NewList(oauthProviderObject),
			Description: "Показать список провайдеров Oauth2",
			Args:        gq.FieldConfigArgument{},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {

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
		},

		"get_params": &gq.Field{
			Type:        appParamsObject,
			Description: "Показать параметры приложения",
			Args:        gq.FieldConfigArgument{},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				return app.GetParams(), nil
			},
		},

		"set_params": &gq.Field{
			Type:        appParamsObject,
			Description: "Установить параметры приложения",
			Args: gq.FieldConfigArgument{

				"selfreg": &gq.ArgumentConfig{
					Type:         gq.Boolean,
					Description:  "Могут ли пользователи регистрироваться самостоятельно",
					DefaultValue: false,
				},
				"use_captcha": &gq.ArgumentConfig{
					Type:         gq.Boolean,
					Description:  "Нужно ли вводить капчу при входе в систему",
					DefaultValue: true,
				},
				"use_pin": &gq.ArgumentConfig{
					Type:         gq.Boolean,
					Description:  "Нужно ли вводить PIN при входе в систему",
					DefaultValue: false,
				},
				"max_attempts": &gq.ArgumentConfig{
					Type:         gq.Int,
					Description:  "Максимально допустимое число ошибок ввода пароля",
					DefaultValue: 5,
				},
				"reset_time": &gq.ArgumentConfig{
					Type:         gq.Int,
					Description:  "Время сброса счетчика ошибок пароля в минутах",
					DefaultValue: 60,
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				if !isAuthAdmin(params) {
					return nil, errors.New("Sorry. You have no admin rights")
				}

				app.Params.Selfreg = params.Args["selfreg"].(bool)
				app.Params.UseCaptcha = params.Args["use_captcha"].(bool)
				app.Params.UsePin = params.Args["use_pin"].(bool)
				app.Params.MaxAttempts = int64(params.Args["max_attempts"].(int))
				app.Params.ResetTime = int64(params.Args["reset_time"].(int))

				SelfRegistrationAllowed = app.Params.Selfreg
				UseCaptcha = app.Params.UseCaptcha
				counter.MAX_ATTEMPTS = app.Params.MaxAttempts
				counter.RESET_TIME = time.Duration(app.Params.ResetTime)

				return app.GetParams(), nil
			},
		},

		"is_captcha_required": &gq.Field{
			Type:        isCaptchaRequiredObject,
			Description: "Нужно ли пользователю вводить каптчу",
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Идентификатор пользователя",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
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
		},

		"is_pin_required": &gq.Field{
			Type:        isPinRequiredObject,
			Description: "Информация о необходимости вводить PIN для входа в систему",
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Идентификатор пользователя",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				pinRequired, pinSet, _, err := authenticator.GetUserPinFields(params.Args["username"].(string))
				return gin.H{"use_pin": app.Params.UsePin, "pinrequired": pinRequired, "pinset": pinSet}, err
			},
		},

		"get_user": &gq.Field{
			Type:        userObject,
			Description: "Показать пользователя",
			Args: gq.FieldConfigArgument{
				"username": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя пользователя",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				ArgToLowerCase(params, "username")
				panicIfNotOwnerOrAdminOrAuditor(params)
				fields := getSelectedFields([]string{"get_user"}, params)
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
				res, err := db.QueryRowMap("SELECT "+fields+` FROM "user" WHERE username = $1 ;`, username)
				return res, err
			},
		},

		"get_app": &gq.Field{
			Type:        appObject,
			Description: "Показать приложение",
			Args: gq.FieldConfigArgument{
				"appname": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.String),
					Description: "Имя приложения",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				panicIfNotAdminOrAuditor(params)
				fields := getSelectedFields([]string{"get_app"}, params)
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
		},

		"list_user_by_ids": &gq.Field{
			Type:        gq.NewList(userObject),
			Description: "Получить список пользователей по их id.",
			Args: gq.FieldConfigArgument{
				"ids": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.NewList(gq.Int)),
					Description: "Массив идентификаторов id",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
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
		},

		"list_user_by_usernames": &gq.Field{
			Type:        gq.NewList(userObject),
			Description: "Получить список пользователей по их username.",
			Args: gq.FieldConfigArgument{
				"usernames": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.NewList(gq.String)),
					Description: "Массив идентификаторов username",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
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
			Type:        gq.NewList(appUserRoleExtendedObject),
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
				panicIfNotOwnerOrAdminOrAuditor(params)
				fields := getSelectedFields([]string{"list_app_user_role"}, params)
				wherePart := listAppUserRoleWherePart(params)
				query := fmt.Sprintf(`SELECT DISTINCT %s FROM app_user_role_extended %s `, fields, wherePart)
				return db.QuerySliceMap(query)
			},
		},
	},
})

func listAppUserRoleWherePart(params gq.ResolveParams) (wherePart string) {
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

func QueryEnd(params gq.ResolveParams, fieldList string) (wherePart string, orderAndLimits string) {
	var searchConditions []string

	// Если запрос к таблице app и это не админский или аудиторский запрос
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
	return wherePart, orderAndLimits
}

func getListOfAllowedAppNames(userName string) string {
	slice, err := db.QuerySliceMap(`
		SELECT DISTINCT appname
		FROM app_user_role
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
