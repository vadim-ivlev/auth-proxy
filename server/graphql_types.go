package server

import (
	"github.com/graphql-go/graphql"
)

// FIELDS **************************************************
var userFields = graphql.Fields{
	"id": &graphql.Field{
		Type:        graphql.Int,
		Description: "автогенерируемый id пользователя",
	},
	"username": &graphql.Field{
		Type:        graphql.String,
		Description: "Уникальный идентификатор пользователя",
	},
	"email": &graphql.Field{
		Type:        graphql.String,
		Description: "Email пользователя",
	},
	"fullname": &graphql.Field{
		Type:        graphql.String,
		Description: "Полное имя пользователя (Фамилия Имя)",
	},
	"description": &graphql.Field{
		Type:        graphql.String,
		Description: "Описание пользователя",
	},
	"disabled": &graphql.Field{
		Type:        graphql.Int,
		Description: "Если не равно 0, пользователь отключен",
	},
	"pinrequired": &graphql.Field{
		Type:        graphql.Boolean,
		Description: "требуется ли PIN Google Authenticator",
	},
	"pinhash_temp": &graphql.Field{
		Type:        graphql.String,
		Description: "хэш для установки Google Authenticator",
	},
	"pinhash": &graphql.Field{
		Type:        graphql.String,
		Description: "хэш для проверки PIN после установки Google Authenticator",
	},
	"emailhash": &graphql.Field{
		Type:        graphql.String,
		Description: "хэш для подтверждения email",
	},
	"emailconfirmed": &graphql.Field{
		Type:        graphql.Boolean,
		Description: "подтвержеден ли email",
	},
}

var appFields = graphql.Fields{
	"id": &graphql.Field{
		Type:        graphql.Int,
		Description: "Уникальный атогенерируемый идентификатор",
	},
	"appname": &graphql.Field{
		Type:        graphql.String,
		Description: "Уникальное имя приложения",
	},
	"description": &graphql.Field{
		Type:        graphql.String,
		Description: "Дополнительная информация",
	},
	"url": &graphql.Field{
		Type:        graphql.String,
		Description: "URL приложения относительно сервера авторизации",
	},
	"rebase": &graphql.Field{
		Type:        graphql.String,
		Description: "Y - чтобы поменять ссылки на HTML страницах на относительные",
	},
	"public": &graphql.Field{
		Type:        graphql.String,
		Description: "Y - чтобы сделать приложение доступным для пользователей без роли",
	},
	"sign": &graphql.Field{
		Type:        graphql.String,
		Description: "Y - для цифровой подписи запросов к приложению",
	},
	"xtoken": &graphql.Field{
		Type:        graphql.String,
		Description: "Используется для проверки доверенного источника. Отправляется вместе с запросом в HTTP-заголовке X-AuthProxy-Token",
	},
}

var groupFields = graphql.Fields{
	"id": &graphql.Field{
		Type:        graphql.Int,
		Description: "Уникальный идентификатор группы",
	},
	"groupname": &graphql.Field{
		Type:        graphql.String,
		Description: "Уникальное  имя группы",
	},
	"description": &graphql.Field{
		Type:        graphql.String,
		Description: "Дополнительная информация",
	},
}

var groupUserRoleFields = graphql.Fields{
	"group_id": &graphql.Field{
		Type:        graphql.Int,
		Description: "Идентификатор группы",
	},
	"user_id": &graphql.Field{
		Type:        graphql.Int,
		Description: "Идентификатор пользователя",
	},
	"rolename": &graphql.Field{
		Type:        graphql.String,
		Description: "Роль пользователя в группе",
	},

	"group_groupname": &graphql.Field{
		Type:        graphql.String,
		Description: "Имя группы",
	},
	"group_description": &graphql.Field{
		Type:        graphql.String,
		Description: "Описание группы",
	},

	"user_email": &graphql.Field{
		Type:        graphql.String,
		Description: "Email пользователя",
	},
	"user_fullname": &graphql.Field{
		Type:        graphql.String,
		Description: "Фамилия Имя пользователя",
	},
	"user_description": &graphql.Field{
		Type:        graphql.String,
		Description: "Описание пользователя",
	},
	"user_disabled": &graphql.Field{
		Type:        graphql.Int,
		Description: "Если не равно 0, пользователь отключен",
	},
}

var groupAppRoleFields = graphql.Fields{
	"group_id": &graphql.Field{
		Type:        graphql.Int,
		Description: "Идентификатор группы",
	},
	"app_id": &graphql.Field{
		Type:        graphql.Int,
		Description: "Идентификатор приложения",
	},
	"rolename": &graphql.Field{
		Type:        graphql.String,
		Description: "Роль группы в приложении",
	},

	"group_groupname": &graphql.Field{
		Type:        graphql.String,
		Description: "Имя группы",
	},
	"group_description": &graphql.Field{
		Type:        graphql.String,
		Description: "Описание группы",
	},

	"app_appname": &graphql.Field{
		Type:        graphql.String,
		Description: "Имя приложения",
	},
	"app_description": &graphql.Field{
		Type:        graphql.String,
		Description: "Описание приложения",
	},
	"app_url": &graphql.Field{
		Type:        graphql.String,
		Description: "URL приложения",
	},
}

var appUserRoleFields = graphql.Fields{
	"appname": &graphql.Field{
		Type:        graphql.String,
		Description: "Имя приложения",
	},
	"username": &graphql.Field{
		Type:        graphql.String,
		Description: "Имя пользователя",
	},
	"rolename": &graphql.Field{
		Type:        graphql.String,
		Description: "Роль пользователя в данном приложении",
	},
}

var appUserRoleExtendedFields = graphql.Fields{
	"appname": &graphql.Field{
		Type:        graphql.String,
		Description: "Идентификатор приложения",
	},
	"username": &graphql.Field{
		Type:        graphql.String,
		Description: "Идентификатор пользователя",
	},
	"rolename": &graphql.Field{
		Type:        graphql.String,
		Description: "Роль пользователя в данном приложении",
	},

	"user_email": &graphql.Field{
		Type:        graphql.String,
		Description: "Email пользователя",
	},
	"user_fullname": &graphql.Field{
		Type:        graphql.String,
		Description: "Полное имя пользователя (Фамилия Имя)",
	},
	"user_description": &graphql.Field{
		Type:        graphql.String,
		Description: "Описание пользователя",
	},
	"user_disabled": &graphql.Field{
		Type:        graphql.Int,
		Description: "Если не равно 0, пользователь отключен",
	},
	"app_description": &graphql.Field{
		Type:        graphql.String,
		Description: "Описание приложения",
	},
	"app_url": &graphql.Field{
		Type:        graphql.String,
		Description: "URL приложения относительно сервера авторизации",
	},
}

// FULL FIELDS поля с древовидной структурой  ****************************************************

var listUserFields = graphql.Fields{
	"length": &graphql.Field{
		Type:        graphql.Int,
		Description: "Количество элементов в списке",
	},
	"list": &graphql.Field{
		Type:        graphql.NewList(userObject),
		Description: "Список пользователей",
	},
}

var listAppFields = graphql.Fields{
	"length": &graphql.Field{
		Type:        graphql.Int,
		Description: "Количество элементов в списке",
	},
	"list": &graphql.Field{
		Type:        graphql.NewList(appObject),
		Description: "Список приложений",
	},
}

var listGroupFields = graphql.Fields{
	"length": &graphql.Field{
		Type:        graphql.Int,
		Description: "Количество элементов в списке",
	},
	"list": &graphql.Field{
		Type:        graphql.NewList(groupObject),
		Description: "Список групп",
	},
}

// TYPES ****************************************************

var userObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "Пользователь",
	Fields:      userFields,
})

var appObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "App",
	Description: "Приложение к которому требуется авторизация",
	Fields:      appFields,
})

var groupObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Group",
	Description: "Группа пользователей",
	Fields:      groupFields,
})

var groupUserRoleObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "GroupUserRole",
	Description: "Роль пользователя в группе",
	Fields:      groupUserRoleFields,
})

var groupAppRoleObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "GroupAppRole",
	Description: "Роль группы в приложении",
	Fields:      groupAppRoleFields,
})

var appUserRoleObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "AppUserRole",
	Description: "Роль пользователя в приложении",
	Fields:      appUserRoleFields,
})

var appUserRoleExtendedObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "AppUserRoleExtended",
	Description: "Роль пользователя в приложении с дополнительными полями из справочных таблиц",
	Fields:      appUserRoleExtendedFields,
})

// LISTS *************************************************************

var listUserGQType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "ListUser",
	Description: "Список пользователей и их количество",
	Fields:      listUserFields,
})

var listAppGQType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "ListApp",
	Description: "Список приложений и их количество",
	Fields:      listAppFields,
})

var listGroupGQType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "ListGroup",
	Description: "Список групп и их количество",
	Fields:      listGroupFields,
})

// AUTH messages
var authMessageObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "AuthMessage",
	Description: "Сообщения AUTH",
	Fields: graphql.Fields{
		"email": &graphql.Field{
			Type:        graphql.String,
			Description: "Email пользователя",
		},
		"message": &graphql.Field{
			Type:        graphql.String,
			Description: "Сообщение",
		},
	},
})

// Сообщения GraphQL метода is_captcha_required()
var isCaptchaRequiredObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "IsCaptchaRequired",
	Description: "Сообщения метода is_captcha_required()",
	Fields: graphql.Fields{
		"is_required": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Будет ли анализироваться капча при следующем вызове метода login()",
		},
		"path": &graphql.Field{
			Type:        graphql.String,
			Description: "URI каптчи",
		},
	},
})

// Сообщения GraphQL метода is_pin_required()
var isPinRequiredObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "IsPinRequired",
	Description: "Сообщения метода is_pin_required()",
	Fields: graphql.Fields{
		"use_pin": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Глобальный флаг нужно ли вводить PIN Google Authenticator при входе в систему",
		},
		"pinrequired": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "требуется ли вводить PIN Google Authenticator для входа",
		},
	},
})

// Сообщения GraphQL метода get_stat()
var statObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Stat",
	Description: "statistics about the memory allocator and requests",
	Fields: graphql.Fields{
		"alloc": &graphql.Field{
			Type:        graphql.Float,
			Description: "megabytes of allocated heap objects",
		},
		"total_alloc": &graphql.Field{
			Type:        graphql.Float,
			Description: "cumulative megabytes allocated for heap objects",
		},
		"sys": &graphql.Field{
			Type:        graphql.Float,
			Description: "the total megabytes of memory obtained from the OS",
		},
		"requests_per_day": &graphql.Field{
			Type:        graphql.Int,
			Description: "requests per day",
		},
		"requests_per_hour": &graphql.Field{
			Type:        graphql.Int,
			Description: "requests per hour",
		},
		"requests_per_minute": &graphql.Field{
			Type:        graphql.Int,
			Description: "requests per minute",
		},
		"requests_per_second": &graphql.Field{
			Type:        graphql.Int,
			Description: "requests per second",
		},
	},
})

var appParamsObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "AppParams",
	Description: "Параметры приложения",
	Fields: graphql.Fields{
		"app_name": &graphql.Field{
			Type:        graphql.String,
			Description: "Имя приложения. Используется для генерации PIN Google Authenticator",
		},
		"selfreg": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Могут ли пользователи регистрироваться самостоятельно",
		},
		"use_captcha": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Нужно ли вводить капчу при входе в систему",
		},
		"use_pin": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Нужно ли вводить PIN при входе в систему",
		},
		"login_not_confirmed_email": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Разрешить авторизацию пользователей не подтвердивших email",
		},
		"no_schema": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "подавлять чтение схемы GraphQL",
		},
		"max_attempts": &graphql.Field{
			Type:        graphql.Int,
			Description: "Максимально допустимое число ошибок ввода пароля",
		},
		"reset_time": &graphql.Field{
			Type:        graphql.Int,
			Description: "Время сброса счетчика ошибок пароля в минутах",
		},
	},
})

var oauthProviderObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "OauthProvider",
	Description: "Сведения о провайдере сервиса аутентификации Oauth2",
	Fields: graphql.Fields{
		"provider_name": &graphql.Field{
			Type:        graphql.String,
			Description: "Имя сервиса",
		},
		"login_endpoint": &graphql.Field{
			Type:        graphql.String,
			Description: "URI точки входа в систему аутентификации",
		},
		"logout_endpoint": &graphql.Field{
			Type:        graphql.String,
			Description: "URI точки выхода в систему аутентификации",
		},
	},
})

// var insertResultObject = gq.NewObject(gq.ObjectConfig{
// 	Name:        "InsertResult",
// 	Description: "Результаты вставки записи",
// 	Fields: gq.Fields{
// 		"RowsAffected": &gq.Field{
// 			Type:        gq.Int,
// 			Description: "Количество порожденных записей",
// 		},
// 	},
// })

var enumAddGroupType = graphql.NewEnum(graphql.EnumConfig{
	Name:        "AdditionalGroup",
	Description: "Дополнительная группа пользователя нового пользователя ",
	Values: graphql.EnumValueConfigMap{
		"subsmag": &graphql.EnumValueConfig{
			Value: "subsmag",
		},
	},
})
