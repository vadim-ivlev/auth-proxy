package controller

import (
	"encoding/json"

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

// addFields возвращает новую структуру graphql.Fields с суммой полей входных структур.
func addFields(fields1 gq.Fields, fields2 gq.Fields) gq.Fields {
	sumFields := gq.Fields{}

	for key, field := range fields1 {
		sumFields[key] = field
	}
	for key, field := range fields2 {
		sumFields[key] = field
	}
	return sumFields
}

// FIELDS **************************************************
var broadcastFields = gq.Fields{
	"id": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор трансляции",
	},
	"title": &gq.Field{
		Type:        gq.String,
		Description: "Название трансляции",
	},
	"time_created": &gq.Field{
		Type:        gq.Int,
		Description: "Время создания",
	},
	"time_begin": &gq.Field{
		Type:        gq.Int,
		Description: "Время начала",
	},
	"is_ended": &gq.Field{
		Type:        gq.Int,
		Description: "Завершена 0 1",
	},
	"show_date": &gq.Field{
		Type:        gq.Int,
		Description: "Показывать дату 0 1",
	},
	"show_time": &gq.Field{
		Type:        gq.Int,
		Description: "Показывать время 0 1",
	},
	"show_main_page": &gq.Field{
		Type:        gq.Int,
		Description: "Показывать на главной странице 01",
	},
	"link_article": &gq.Field{
		Type:        gq.String,
		Description: "Ссылка на статью",
	},
	"link_img": &gq.Field{
		Type:        gq.String,
		Description: "Ссылка на изображение",
	},
	"groups_create": &gq.Field{
		Type:        gq.Int,
		Description: "",
	},
	"is_diary": &gq.Field{
		Type:        gq.Int,
		Description: "Дневник 01",
	},
	"diary_author": &gq.Field{
		Type:        gq.String,
		Description: "Автор дневника",
	},
}

var postFields = gq.Fields{
	"id": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор поста",
	},
	"id_parent": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор родительского поста если это ответ на другой пост",
	},
	"id_broadcast": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор трансляции",
	},
	"text": &gq.Field{
		Type:        gq.String,
		Description: "Текст поста",
	},
	"post_time": &gq.Field{
		Type:        gq.Int,
		Description: "Текст поста",
	},
	"post_type": &gq.Field{
		Type:        gq.Int,
		Description: "Тип поста 1,2,3,4...",
	},
	"link": &gq.Field{
		Type:        gq.String,
		Description: "Ссылка",
	},
	"has_big_img": &gq.Field{
		Type:        gq.Int,
		Description: "Есть ли большое изображение 0,1",
	},
	"author": &gq.Field{
		Type:        gq.String,
		Description: "ФИО автора поста",
	},
}

var imageFields = gq.Fields{
	"id": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор медиа",
	},
	"post_id": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор поста",
	},
	"filepath": &gq.Field{
		Type:        gq.String,
		Description: "URI изображения",
	},
	"thumbs": &gq.Field{
		Type:        gq.NewList(thumbType),
		Description: "Превью и изображение для видео - jsonb ",
		Resolve: func(params gq.ResolveParams) (interface{}, error) {
			return JSONParamToMap(params, "thumbs")
		},
	},
	"source": &gq.Field{
		Type:        gq.String,
		Description: "Источник медиа",
	},
	"width": &gq.Field{
		Type:        gq.Int,
		Description: "Ширина в пикселях",
	},
	"height": &gq.Field{
		Type:        gq.Int,
		Description: "Высота в пикселях",
	},
}

var listBroadcastFields = gq.Fields{
	"length": &gq.Field{
		Type:        gq.Int,
		Description: "Количество элементов в списке",
	},
	"list": &gq.Field{
		Type:        gq.NewList(broadcastType),
		Description: "Список трансляций",
	},
}

// FULL FIELDS поля с древовидной структурой  ****************************************************

var fullAnswerFields = addFields(postFields, gq.Fields{

	"images": &gq.Field{
		Type:        gq.NewList(imageType),
		Description: "Медиа ответа",
		Resolve: func(params gq.ResolveParams) (interface{}, error) {
			return JSONParamToMap(params, "images")
		},
	},
})

var fullPostFields = addFields(postFields, gq.Fields{
	"images": &gq.Field{
		Type:        gq.NewList(imageType),
		Description: "Медиа поста",
		Resolve: func(params gq.ResolveParams) (interface{}, error) {
			return JSONParamToMap(params, "images")
		},
	},
	"answers": &gq.Field{
		Type:        gq.NewList(fullAnswerType),
		Description: "Ответы к посту",
		Resolve: func(params gq.ResolveParams) (interface{}, error) {
			return JSONParamToMap(params, "answers")
		},
	},
})

var fullBroadcastFields = addFields(broadcastFields, gq.Fields{
	"posts": &gq.Field{
		Type:        gq.NewList(fullPostType),
		Description: "Посты бродкаста",
		Resolve: func(params gq.ResolveParams) (interface{}, error) {
			return JSONParamToMap(params, "posts")
		},
	},
})

var fullListBroadcastFields = gq.Fields{
	"length": &gq.Field{
		Type:        gq.Int,
		Description: "Количество элементов в списке",
	},
	"list": &gq.Field{
		Type:        gq.NewList(fullBroadcastType),
		Description: "Список трансляций c постами, ответами и медиа",
	},
}

// TYPES ****************************************************

var postType = gq.NewObject(gq.ObjectConfig{
	Name:        "Post",
	Description: "Пост трансляции",
	Fields:      postFields,
})

var imageType = gq.NewObject(gq.ObjectConfig{
	Name:        "Image",
	Description: "Медиа поста трансляции",
	Fields:      imageFields,
})

var broadcastType = gq.NewObject(gq.ObjectConfig{
	Name:        "Broadcast",
	Description: "Онлайн трансляция",
	Fields:      broadcastFields,
})

var listBroadcastType = gq.NewObject(gq.ObjectConfig{
	Name:        "ListBroadcast",
	Description: "Список трансляций и количество элементов в списке",
	Fields:      listBroadcastFields,
})

// FULL TYPES типы с древовидной структурой *************

var fullPostType = gq.NewObject(gq.ObjectConfig{
	Name:        "FullPost",
	Description: "Пост трансляции с медиа и ответами к посту",
	Fields:      fullPostFields,
})

var fullAnswerType = gq.NewObject(gq.ObjectConfig{
	Name:        "FullAnswer",
	Description: "Ответ к посту с медиа ответа",
	Fields:      fullAnswerFields,
})

var fullBroadcastType = gq.NewObject(gq.ObjectConfig{
	Name:        "FullBroadcast",
	Description: "Трансляция с постами",
	Fields:      fullBroadcastFields,
})

var fullListBroadcastType = gq.NewObject(gq.ObjectConfig{
	Name:        "FullListBroadcast",
	Description: "Список трансляций c постами, ответами и медиа,  и количество элементов в списке",
	Fields:      fullListBroadcastFields,
})

var thumbType = gq.NewObject(gq.ObjectConfig{
	Name:        "Thumb",
	Description: "Уменьшенное изображение для видео",
	Fields: gq.Fields{
		"type": &gq.Field{
			Type:        gq.String,
			Description: "Тип (small, middle, large)",
		},
		"filepath": &gq.Field{
			Type:        gq.String,
			Description: "Ссылка на файл на сервере",
		},
		"width": &gq.Field{
			Type:        gq.Int,
			Description: "Ширина изображения в пикселях.",
		},
		"height": &gq.Field{
			Type:        gq.Int,
			Description: "Высота изображения в пикселях.",
		},
	},
})
