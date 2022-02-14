package server

import (
	// "encoding/json"

	gq "github.com/graphql-go/graphql"
)

// FIELDS **************************************************
var userFields = gq.Fields{
	"id": &gq.Field{
		Type:        gq.Int,
		Description: "автогенерируемый id пользователя",
	},
	"username": &gq.Field{
		Type:        gq.String,
		Description: "Уникальный идентификатор пользователя",
	},
	"email": &gq.Field{
		Type:        gq.String,
		Description: "Email пользователя",
	},
	"fullname": &gq.Field{
		Type:        gq.String,
		Description: "Полное имя пользователя (Фамилия Имя)",
	},
	"description": &gq.Field{
		Type:        gq.String,
		Description: "Дополнительная информация",
	},
	"disabled": &gq.Field{
		Type:        gq.Int,
		Description: "Если не равно 0, пользователь отключен",
	},
	"pinrequired": &gq.Field{
		Type:        gq.Boolean,
		Description: "требуется ли PIN Google Authenticator",
	},
	"pinhash_temp": &gq.Field{
		Type:        gq.String,
		Description: "хэш для установки Google Authenticator",
	},
	"pinhash": &gq.Field{
		Type:        gq.String,
		Description: "хэш для проверки PIN после установки Google Authenticator",
	},
	"emailhash": &gq.Field{
		Type:        gq.String,
		Description: "хэш для подтверждения email",
	},
	"emailconfirmed": &gq.Field{
		Type:        gq.Boolean,
		Description: "подтвержеден ли email",
	},
}

var appFields = gq.Fields{
	"appname": &gq.Field{
		Type:        gq.String,
		Description: "Уникальный идентификатор приложения",
	},
	"description": &gq.Field{
		Type:        gq.String,
		Description: "Дополнительная информация",
	},
	"url": &gq.Field{
		Type:        gq.String,
		Description: "URL приложения относительно сервера авторизации",
	},
	"rebase": &gq.Field{
		Type:        gq.String,
		Description: "Y - чтобы поменять ссылки на HTML страницах на относительные",
	},
	"public": &gq.Field{
		Type:        gq.String,
		Description: "Y - чтобы сделать приложение доступным для пользователей без роли",
	},
	"sign": &gq.Field{
		Type:        gq.String,
		Description: "Y - для цифровой подписи запросов к приложению",
	},
}

var groupFields = gq.Fields{
	"id": &gq.Field{
		Type:        gq.Int,
		Description: "Уникальный идентификатор группы",
	},
	"groupname": &gq.Field{
		Type:        gq.String,
		Description: "Уникальное  имя группы",
	},
	"description": &gq.Field{
		Type:        gq.String,
		Description: "Дополнительная информация",
	},
}

var groupUserRoleFields = gq.Fields{
	"group_id": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор группы",
	},
	"user_id": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор пользователя",
	},
	"rolename": &gq.Field{
		Type:        gq.String,
		Description: "Роль пользователя в группе",
	},
}

var groupAppRoleFields = gq.Fields{
	"group_id": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор группы",
	},
	"app_id": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор приложения",
	},
	"rolename": &gq.Field{
		Type:        gq.String,
		Description: "Роль группы в приложении",
	},
}

var appUserRoleFields = gq.Fields{
	"appname": &gq.Field{
		Type:        gq.String,
		Description: "Идентификатор приложения",
	},
	"username": &gq.Field{
		Type:        gq.String,
		Description: "Идентификатор пользователя",
	},
	"rolename": &gq.Field{
		Type:        gq.String,
		Description: "Роль пользователя в данном приложении",
	},
}

var appUserRoleExtendedFields = gq.Fields{
	"appname": &gq.Field{
		Type:        gq.String,
		Description: "Идентификатор приложения",
	},
	"username": &gq.Field{
		Type:        gq.String,
		Description: "Идентификатор пользователя",
	},
	"rolename": &gq.Field{
		Type:        gq.String,
		Description: "Роль пользователя в данном приложении",
	},

	"user_email": &gq.Field{
		Type:        gq.String,
		Description: "Email пользователя",
	},
	"user_fullname": &gq.Field{
		Type:        gq.String,
		Description: "Полное имя пользователя (Фамилия Имя)",
	},
	"user_description": &gq.Field{
		Type:        gq.String,
		Description: "Описание пользователя",
	},
	"user_disabled": &gq.Field{
		Type:        gq.Int,
		Description: "Если не равно 0, пользователь отключен",
	},
	"app_description": &gq.Field{
		Type:        gq.String,
		Description: "Описание приложения",
	},
	"app_url": &gq.Field{
		Type:        gq.String,
		Description: "URL приложения относительно сервера авторизации",
	},
}

// FULL FIELDS поля с древовидной структурой  ****************************************************

var listUserFields = gq.Fields{
	"length": &gq.Field{
		Type:        gq.Int,
		Description: "Количество элементов в списке",
	},
	"list": &gq.Field{
		Type:        gq.NewList(userObject),
		Description: "Список пользователей",
	},
}

var listAppFields = gq.Fields{
	"length": &gq.Field{
		Type:        gq.Int,
		Description: "Количество элементов в списке",
	},
	"list": &gq.Field{
		Type:        gq.NewList(appObject),
		Description: "Список приложений",
	},
}

// TYPES ****************************************************

var userObject = gq.NewObject(gq.ObjectConfig{
	Name:        "User",
	Description: "Пользователь",
	Fields:      userFields,
})

var appObject = gq.NewObject(gq.ObjectConfig{
	Name:        "App",
	Description: "Приложение к которому требуется авторизация",
	Fields:      appFields,
})

var groupObject = gq.NewObject(gq.ObjectConfig{
	Name:        "Group",
	Description: "Группа пользователей",
	Fields:      groupFields,
})

var appUserRoleObject = gq.NewObject(gq.ObjectConfig{
	Name:        "AppUserRole",
	Description: "Роль пользователя в приложении",
	Fields:      appUserRoleFields,
})

var appUserRoleExtendedObject = gq.NewObject(gq.ObjectConfig{
	Name:        "AppUserRoleExtended",
	Description: "Роль пользователя в приложении с дополнительными полями из справочных таблиц",
	Fields:      appUserRoleExtendedFields,
})

// LISTS *************************************************************

var listUserGQType = gq.NewObject(gq.ObjectConfig{
	Name:        "ListUser",
	Description: "Список пользователей и их количество",
	Fields:      listUserFields,
})

var listAppGQType = gq.NewObject(gq.ObjectConfig{
	Name:        "ListApp",
	Description: "Список приложений и их количество",
	Fields:      listAppFields,
})

// AUTH messages
var authMessageObject = gq.NewObject(gq.ObjectConfig{
	Name:        "AuthMessage",
	Description: "Сообщения AUTH",
	Fields: gq.Fields{
		"email": &gq.Field{
			Type:        gq.String,
			Description: "Email пользователя",
		},
		"message": &gq.Field{
			Type:        gq.String,
			Description: "Сообщение",
		},
	},
})

// Сообщения GraphQL метода is_captcha_required()
var isCaptchaRequiredObject = gq.NewObject(gq.ObjectConfig{
	Name:        "IsCaptchaRequired",
	Description: "Сообщения метода is_captcha_required()",
	Fields: gq.Fields{
		"is_required": &gq.Field{
			Type:        gq.Boolean,
			Description: "Будет ли анализироваться капча при следующем вызове метода login()",
		},
		"path": &gq.Field{
			Type:        gq.String,
			Description: "URI каптчи",
		},
	},
})

// Сообщения GraphQL метода is_pin_required()
var isPinRequiredObject = gq.NewObject(gq.ObjectConfig{
	Name:        "IsPinRequired",
	Description: "Сообщения метода is_pin_required()",
	Fields: gq.Fields{
		"use_pin": &gq.Field{
			Type:        gq.Boolean,
			Description: "Глобальный флаг нужно ли вводить PIN Google Authenticator при входе в систему",
		},
		"pinrequired": &gq.Field{
			Type:        gq.Boolean,
			Description: "требуется ли вводить PIN Google Authenticator для входа",
		},
	},
})

// Сообщения GraphQL метода get_stat()
var statObject = gq.NewObject(gq.ObjectConfig{
	Name:        "Stat",
	Description: "statistics about the memory allocator and requests",
	Fields: gq.Fields{
		"alloc": &gq.Field{
			Type:        gq.Float,
			Description: "megabytes of allocated heap objects",
		},
		"total_alloc": &gq.Field{
			Type:        gq.Float,
			Description: "cumulative megabytes allocated for heap objects",
		},
		"sys": &gq.Field{
			Type:        gq.Float,
			Description: "the total megabytes of memory obtained from the OS",
		},
		"requests_per_day": &gq.Field{
			Type:        gq.Int,
			Description: "requests per day",
		},
		"requests_per_hour": &gq.Field{
			Type:        gq.Int,
			Description: "requests per hour",
		},
		"requests_per_minute": &gq.Field{
			Type:        gq.Int,
			Description: "requests per minute",
		},
		"requests_per_second": &gq.Field{
			Type:        gq.Int,
			Description: "requests per second",
		},
	},
})

var appParamsObject = gq.NewObject(gq.ObjectConfig{
	Name:        "AppParams",
	Description: "Параметры приложения",
	Fields: gq.Fields{
		"app_name": &gq.Field{
			Type:        gq.String,
			Description: "Имя приложения. Используется для генерации PIN Google Authenticator",
		},
		"selfreg": &gq.Field{
			Type:        gq.Boolean,
			Description: "Могут ли пользователи регистрироваться самостоятельно",
		},
		"use_captcha": &gq.Field{
			Type:        gq.Boolean,
			Description: "Нужно ли вводить капчу при входе в систему",
		},
		"use_pin": &gq.Field{
			Type:        gq.Boolean,
			Description: "Нужно ли вводить PIN при входе в систему",
		},
		"max_attempts": &gq.Field{
			Type:        gq.Int,
			Description: "Максимально допустимое число ошибок ввода пароля",
		},
		"reset_time": &gq.Field{
			Type:        gq.Int,
			Description: "Время сброса счетчика ошибок пароля в минутах",
		},
	},
})

var oauthProviderObject = gq.NewObject(gq.ObjectConfig{
	Name:        "OauthProvider",
	Description: "Сведения о провайдере сервиса аутентификации Oauth2",
	Fields: gq.Fields{
		"provider_name": &gq.Field{
			Type:        gq.String,
			Description: "Имя сервиса",
		},
		"login_endpoint": &gq.Field{
			Type:        gq.String,
			Description: "URI точки входа в систему аутентификации",
		},
		"logout_endpoint": &gq.Field{
			Type:        gq.String,
			Description: "URI точки выхода в систему аутентификации",
		},
	},
})
