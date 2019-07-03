package controller

import (
	"auth-proxy/model/db"

	gq "github.com/graphql-go/graphql"
)

var rootMutation = gq.NewObject(gq.ObjectConfig{
	Name: "Mutation",
	Fields: gq.Fields{

		// BROADCAST =====================================================

		"create_broadcast": &gq.Field{
			Type:        fullBroadcastType,
			Description: "Создать трансляцию",
			Args: gq.FieldConfigArgument{
				// "id":             &gq.ArgumentConfig{Type: gq.NewNonNull(gq.Int), Description: "Идентификатор трансляции"},
				"title": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Название трансляции",
				},
				"time_created": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Время создания",
				},
				"time_begin": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Время начала",
				},
				"is_ended": &gq.ArgumentConfig{
					Type:         gq.Int,
					Description:  "Завершена 0 1. По умолчанию 0",
					DefaultValue: 0,
				},
				"show_date": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Показывать дату 0 1",
				},
				"show_time": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Показывать время 0 1",
				},
				"show_main_page": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Показывать на главной странице 0 1",
				},
				"link_article": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Ссылка на статью",
				},
				"link_img": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Ссылка на изображение",
				},
				"groups_create": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "группа",
				},
				"is_diary": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Дневник 0 1",
				},
				"diary_author": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Автор дневника",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				return createRecord(params, "broadcast", "full_broadcast")
			},
		},

		"update_broadcast": &gq.Field{
			Type:        fullBroadcastType,
			Description: "Обновить трансляцию",
			Args: gq.FieldConfigArgument{
				"id": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.Int),
					Description: "Идентификатор трансляции",
				},
				"title": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Название трансляции",
				},
				"time_created": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Время создания",
				},
				"time_begin": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Время начала",
				},
				"is_ended": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Завершена 0 1",
				},
				"show_date": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Показывать дату 0 1",
				},
				"show_time": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Показывать время 0 1",
				},
				"show_main_page": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Показывать на главной странице 0 1",
				},
				"link_article": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Ссылка на статью",
				},
				"link_img": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Ссылка на изображение",
				},
				"groups_create": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "группа",
				},
				"is_diary": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Дневник 0 1",
				},
				"diary_author": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Автор дневника",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				res, err := updateRecord(params, "broadcast", "full_broadcast")
				if err == nil {
					// redis.ClearByBroadcastID(params.Args["id"])
				}
				return res, err
			},
		},

		"delete_broadcast": &gq.Field{
			Type:        fullBroadcastType,
			Description: "Удалить трансляцию",
			Args: gq.FieldConfigArgument{
				"id": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.Int),
					Description: "Идентификатор трансляции",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				res, err := deleteRecord(params, "broadcast", "full_broadcast")
				if err == nil {
					// redis.ClearByBroadcastID(params.Args["id"])
				}
				return res, err

			},
		},

		// POST =====================================================

		"create_post": &gq.Field{
			Type:        fullPostType,
			Description: "Создать пост к трансляции или ответ к посту",
			Args: gq.FieldConfigArgument{
				// "id":           &gq.ArgumentConfig{Type: gq.NewNonNull(gq.Int), Description: "Идентификатор поста"}                              ,
				"id_parent": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Идентификатор поста если это ответ на другой пост",
				},
				"id_broadcast": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Идентификатор трансляции",
				},
				"text": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Текст поста",
				},
				"post_time": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Текст поста",
				},
				"post_type": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Тип поста 1,2,3,4...",
				},
				"link": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Ссылка",
				},
				"has_big_img": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Есть ли большое изображение 0,1",
				},
				"author": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "ФИО автора поста",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				res, err := createRecord(params, "post", "full_post")
				if err == nil {
					// redis.ClearByBroadcastID(params.Args["id_broadcast"])
					// redis.ClearByPostID(params.Args["id_parent"])
				}
				return res, err

			},
		},

		"update_post": &gq.Field{
			Type:        fullPostType,
			Description: "Обновить пост или ответ к посту",
			Args: gq.FieldConfigArgument{
				"id": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.Int),
					Description: "Идентификатор поста",
				},
				"id_parent": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Идентификатор поста если это ответ на другой пост",
				},
				"id_broadcast": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Идентификатор трансляции",
				},
				"text": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Текст поста",
				},
				"post_time": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Текст поста",
				},
				"post_type": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Тип поста 1,2 ,3,4...",
				},
				"link": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Ссылка",
				},
				"has_big_img": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Есть ли большое изображение 0,1",
				},
				"author": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "ФИО автора поста",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				// // запоминаем старых родителей, на случай перепривязки
				// id_broadcast_old, id_parend_old, err := redis.GetPostParentIDs(params.Args["id"])
				// if err != nil {
				// 	log.Println("update_post:getPostParentIDs:", err)
				// }

				// // если id_parent отрицательный отвязываем запись
				// id_parent, ok := params.Args["id_parent"].(int)
				// if ok && id_parent < 0 {
				// 	params.Args["id_parent"] = nil
				// }

				res, err := updateRecord(params, "post", "full_post")
				// if err == nil {
				// 	redis.ClearByBroadcastID(id_broadcast_old)
				// 	redis.ClearByPostID(id_parend_old)
				// 	redis.ClearByPostID(params.Args["id"])
				// }
				return res, err
			},
		},

		"delete_post": &gq.Field{
			Type:        fullPostType,
			Description: "Удалить пост",
			Args: gq.FieldConfigArgument{
				"id": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.Int),
					Description: "Идентификатор поста"},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				// // запоминаем старых родителей, на случай перепривязки
				// id_broadcast_old, id_parend_old, err := redis.GetPostParentIDs(params.Args["id"])
				// if err != nil {
				// 	log.Println("delete_post:getPostParentIDs:", err)
				// }

				res, err := deleteRecord(params, "post", "full_post")
				// if err == nil {
				// 	redis.ClearByBroadcastID(id_broadcast_old)
				// 	redis.ClearByPostID(id_parend_old)
				// }
				return res, err
			},
		},

		// MEDIA =====================================================

		"create_image": &gq.Field{
			Type:        imageType,
			Description: "Создать медиа",
			Args: gq.FieldConfigArgument{
				// "id":           &gq.ArgumentConfig{Type: gq.NewNonNull(gq.Int), Description: "Идентификатор изображения"},
				"post_id": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Идентификатор поста",
				},
				"filepath": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "URI изображения",
				},
				"source": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Источник медиа",
				},
				"file_field_name": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Имя (name) поля формы для загрузки файла. <input name='fname' type='file' ...>",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {

				// _, ok := params.Args["file_field_name"].(string)
				// if ok {
				// 	path, width, height, thumbs, errMsg := SaveUploadedImage(params, "file_field_name")
				// 	if errMsg == "" {
				// 		params.Args["filepath"] = path
				// 		params.Args["thumbs"] = thumbs
				// 		params.Args["width"] = width
				// 		params.Args["height"] = height
				// 	} else {
				// 		msg := "create_image: Resolve(): " + errMsg
				// 		log.Println(msg)
				// 		return nil, errors.New(msg)
				// 	}
				// }
				delete(params.Args, "file_field_name")
				res, err := db.CreateRow("image", params.Args)
				// if err == nil {
				// 	redis.ClearByPostID(params.Args["post_id"])
				// }
				return res, err

			},
		},

		"update_image": &gq.Field{
			Type:        imageType,
			Description: "Обновить медиа по идентификатору",
			Args: gq.FieldConfigArgument{
				"id": &gq.ArgumentConfig{Type: gq.NewNonNull(gq.Int),
					Description: "Идентификатор медиа",
				},
				"post_id": &gq.ArgumentConfig{Type: gq.Int,
					Description: "Идентификатор поста",
				},
				"filepath": &gq.ArgumentConfig{Type: gq.String,
					Description: "URI изображения",
				},
				"source": &gq.ArgumentConfig{Type: gq.String,
					Description: "Источник медиа",
				},
				"file_field_name": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Имя (name) поля формы для загрузки файла. <input name='fname' type='file' ...>",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {

				// _, ok := params.Args["file_field_name"].(string)
				// if ok {
				// 	path, width, height, thumbs, errMsg := SaveUploadedImage(params, "file_field_name")
				// 	if errMsg == "" {
				// 		params.Args["filepath"] = path
				// 		params.Args["thumbs"] = thumbs
				// 		params.Args["width"] = width
				// 		params.Args["height"] = height
				// 	} else {
				// 		msg := "update_image: Resolve(): " + errMsg
				// 		log.Println(msg)
				// 		return nil, errors.New(msg)
				// 	}
				// }
				delete(params.Args, "file_field_name")

				res, err := db.UpdateRowByID("image", params.Args["id"].(int), params.Args)
				// if err == nil {
				// 	redis.ClearByImageID(params.Args["id"])
				// }
				return res, err

			},
		},

		"delete_image": &gq.Field{
			Type:        imageType,
			Description: "Удалить медиа по идентификатору",
			Args: gq.FieldConfigArgument{
				"id": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.Int),
					Description: "Идентификатор медиа",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				// post_id_old := redis.GetImagePostID(params.Args["id"])
				res, err := db.DeleteRowByID("image", params.Args["id"].(int))
				// if err == nil {
				// 	redis.ClearByPostID(post_id_old)
				// }
				return res, err
			},
		},
	},
})
