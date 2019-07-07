package controller

import (
	"auth-proxy/model/auth"
	"context"
	"encoding/json"
	"errors"
	"log"

	// "go/ast"
	"net/http"
	"strings"

	"auth-proxy/model/db"

	"github.com/gin-gonic/gin"
	gq "github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// FUNCTIONS *******************************************************

// jsonStringToMap преобразует строку JSON в map[string]interface{}
func jsonStringToMap(s string) map[string]interface{} {
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(s), &m)
	return m
}

// getParamsFromBody извлекает параметры запроса из тела запроса
func getParamsFromBody(c *gin.Context) (map[string]interface{}, error) {
	r := c.Request
	mb := make(map[string]interface{})
	if r.ContentLength > 0 {
		errBodyDecode := json.NewDecoder(r.Body).Decode(&mb)
		return mb, errBodyDecode
	}
	return mb, errors.New("No body")
}

// getPayload3 извлекает "query", "variables", "operationName".
// Decoded body has precedence over POST over GET.
func getPayload3(c *gin.Context) (query string, variables map[string]interface{}) {

	// Проверяем на существование данных из Form Data
	query = c.PostForm("query")
	variables = jsonStringToMap(c.PostForm("variables"))

	// если есть тело запроса то берем из Request Payload (для Altair)
	params, errBody := getParamsFromBody(c)
	if errBody == nil {
		query, _ = params["query"].(string)
		variables, _ = params["variables"].(map[string]interface{})
	}

	return
}

// createRecord вставляет запись в таблицу tableToUpdate,
// и возвращает вставленную  запись из таблицы tableToSelectFrom,
// которая является представлением с более богатым содержимым,
// чем обновленная таблица.
// Используется в запросах GraphQL на вставку записей.
func createRecord(keyFieldName string, params gq.ResolveParams, tableToUpdate string, tableToSelectFrom string) (interface{}, error) {
	// вставляем запись
	fieldValues, err := db.CreateRow(tableToUpdate, params.Args)
	if err != nil {
		return fieldValues, err
	}
	// извлекаем id вставленной записи
	id := fieldValues[keyFieldName]

	// возвращаем ответ
	path := params.Info.FieldName
	fields := getSelectedFields([]string{path}, params)
	return db.QueryRowMap("SELECT "+fields+" FROM "+tableToSelectFrom+" WHERE "+keyFieldName+" = $1 ;", id)
}

// updateRecord обновляет запись в таблице tableToUpdate,
// и возвращает обновленную запись из таблицы tableToSelectFrom,
// которая является представлением с более богатым содержимым,
// чем обновленная таблица.
// Используется в запросах GraphQL на обновление записей.
func updateRecord(keyFieldName string, params gq.ResolveParams, tableToUpdate string, tableToSelectFrom string) (interface{}, error) {
	id := params.Args[keyFieldName]
	fieldValues, err := db.UpdateRowByID(keyFieldName, tableToUpdate, id, params.Args)
	if err != nil {
		return fieldValues, err
	}
	path := params.Info.FieldName
	fields := getSelectedFields([]string{path}, params)
	return db.QueryRowMap("SELECT "+fields+" FROM "+tableToSelectFrom+" WHERE "+keyFieldName+" = $1 ;", id)
}

// deleteRecord удаляет запись из таблицы tableToUpdate,
// и возвращает удаленную запись из таблицы tableToSelectFrom,
// которая является представлением с более богатым содержимым,
// чем обновленная таблица.
// Используется в запросах GraphQL на удаление записей.
func deleteRecord(keyFieldName string, params gq.ResolveParams, tableToUpdate string, tableToSelectFrom string) (interface{}, error) {
	// сохраняем запись, которую собираемся удалять
	id := params.Args[keyFieldName]
	path := params.Info.FieldName
	fields := getSelectedFields([]string{path}, params)
	fieldValues, err := db.QueryRowMap("SELECT "+fields+" FROM "+tableToSelectFrom+" WHERE "+keyFieldName+" = $1 ;", id)
	if err != nil {
		return nil, err
	}
	// удаляем запись
	_, err = db.DeleteRowByID(keyFieldName, tableToUpdate, id)
	if err != nil {
		return nil, err
	}
	// возвращаем только что удаленную запись
	return fieldValues, nil
}

// getSelectedFields - returns list of selected fields defined in GraphQL query.
/*

	Invoke it like this:

	getSelectedFields([]string{"companies"}, resolveParams)
	// this will return []string{"id", "name"}
	In case you have a "path" you want to select from, e.g.

	query {
	a {
		b {
		x,
		y,
		z
		}
	}
	}
	Then you'd call it like this:

	getSelectedFields([]string{"a", "b"}, resolveParams)
	// Returns []string{"x", "y", "z"}

	import "github.com/graphql-go/graphql/language/ast" is added by hands.
	source: https://github.com/graphql-go/graphql/issues/125
*/
func getSelectedFields(selectionPath []string, resolveParams gq.ResolveParams) string {
	fields := resolveParams.Info.FieldASTs
	for _, propName := range selectionPath {
		found := false
		for _, field := range fields {
			if field.Name.Value == propName {
				selections := field.SelectionSet.Selections
				fields = make([]*ast.Field, 0)
				for _, selection := range selections {
					fields = append(fields, selection.(*ast.Field))
				}
				found = true
				break
			}
		}
		if !found {
			return ""
		}
	}
	var collect []string
	for _, field := range fields {
		name := field.Name.Value
		if name != "__typename" {
			collect = append(collect, field.Name.Value)
		}
	}
	s := strings.Join(collect, ", ")
	return s
}

func getLoginedUserName(params gq.ResolveParams) string {
	c, ok := params.Context.Value("ginContext").(*gin.Context)
	if !ok {
		log.Println("Not OK: getLoginedUserName")
		return ""
	}
	return auth.GetUserName(c)
}

func getOwnerUserName(params gq.ResolveParams) string {
	c, ok := params.Context.Value("ginContext").(*gin.Context)
	if !ok {
		log.Println("Not OK: getLoginedUserName")
		return ""
	}
	return auth.GetUserName(c)
}

func panicIfNotOwnerOrAdmin(params gq.ResolveParams) {
	uname := getLoginedUserName(params)
	pname, ok := params.Args["username"].(string)
	if ok && pname == uname {
		return
	}
	panicIfNotAdmin(params)
}

func panicIfNotAdmin(params gq.ResolveParams) {
	username := getLoginedUserName(params)
	if auth.AppUserRoleExist("auth", username, "admin") {
		return
	}
	panic("Sorry. You have to be an admin.")
}

func panicIfNotUser(params gq.ResolveParams) {
	username := getLoginedUserName(params)
	if username == "" {
		panic("Sorry. You have to log in.")
	}
}

// G R A P H Q L ********************************************************************************

var schema, _ = gq.NewSchema(gq.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

// GraphQL исполняет GraphQL запрос
func GraphQL(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 100*1024*1024)

	query, variables := getPayload3(c)

	result := gq.Do(gq.Params{
		Schema:         schema,
		RequestString:  query,
		Context:        context.WithValue(context.Background(), "ginContext", c),
		VariableValues: variables,
	})

	c.JSON(200, result)
}
