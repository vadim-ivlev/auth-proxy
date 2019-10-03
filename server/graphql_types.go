package server

import (
	// "encoding/json"

	gq "github.com/graphql-go/graphql"
)

// F U N C S ***********************************************

// JSONParamToMap - возвращает параметр paramName в map[string]interface{}.
// Второй параметр возврата - ошибка.
// Применяется для сериализации поля JSON таблицы postgres в map.
func JSONParamToMap(params gq.ResolveParams, paramName string) (interface{}, error) {

	source := params.Source.(map[string]interface{})
	param := source[paramName]

	// TODO: may be it's better to check if it can be converted to map[string]interface{}
	paramBytes, ok := param.([]byte)
	if !ok {
		return param, nil
	}
	var paramMap []map[string]interface{}
	err := json.Unmarshal(paramBytes, &paramMap)
	return paramMap, err
}

// FIELDS **************************************************
var userFields = gq.Fields{
	"username": &gq.Field{
		Type:        gq.String,
		Description: "Уникальный идентификатор пользователя",
	},
	"email": &gq.Field{
		Type:        gq.String,
		Description: "Емайл",
	},
	"fullname": &gq.Field{
		Type:        gq.String,
		Description: "Полное имя пользователя",
	},
	"description": &gq.Field{
		Type:        gq.String,
		Description: "Дополнительная информация",
	},
	"disabled": &gq.Field{
		Type:        gq.Int,
		Description: "Если не равно 0, пользователь отключен",
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
}

var app_user_roleFields = gq.Fields{
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

var app_user_role_extendedFields = gq.Fields{
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
		Description: "Емайл пользователя",
	},
	"user_fullname": &gq.Field{
		Type:        gq.String,
		Description: "Полное имя пользователя",
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

var app_user_roleObject = gq.NewObject(gq.ObjectConfig{
	Name:        "AppUserRole",
	Description: "Роль пользователя в приложении",
	Fields:      app_user_roleFields,
})

var app_user_role_extendedObject = gq.NewObject(gq.ObjectConfig{
	Name:        "AppUserRoleExtended",
	Description: "Роль пользователя в приложении с дополнительными полями из справочных таблиц",
	Fields:      app_user_role_extendedFields,
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
	Description: "Сообщения аутентификации",
	Fields: gq.Fields{
		"username": &gq.Field{
			Type:        gq.String,
			Description: "Имя пользователя",
		},
		"message": &gq.Field{
			Type:        gq.String,
			Description: "Сообщение",
		},
	},
})
