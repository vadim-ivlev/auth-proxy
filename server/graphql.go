package server

import (
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/prometeo"
	"context"
	"errors"
	"log"

	// "go/ast"
	"net/http"
	"strings"

	"auth-proxy/pkg/db"

	"github.com/gin-gonic/gin"
	gq "github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	jsoniter "github.com/json-iterator/go"
)

// Для ускорения работы с json. (Drop in) вместо стандартного import "encoding/json"
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// SelfRegistrationAllowed может ли пользователь зарегистрироваться самостоятельно
var SelfRegistrationAllowed = false

// UseCaptcha нужно ли вводить каптчу во время входа
var UseCaptcha = true

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
	// id := fieldValues[keyFieldName]	// возвращаем ответ
	id := params.Args[keyFieldName] // возвращаем ответ

	path := params.Info.FieldName
	fields := getSelectedFields([]string{path}, params)
	return db.QueryRowMap("SELECT "+fields+" FROM \""+tableToSelectFrom+"\" WHERE "+keyFieldName+" = $1 ;", id)
}

// updateRecord обновляет запись в таблице tableToUpdate,
// и возвращает обновленную запись из таблицы tableToSelectFrom,
// которая является представлением с более богатым содержимым,
// чем обновленная таблица.
// Используется в запросах GraphQL на обновление записей.
func updateRecord(oldKeyValue string, keyFieldName string, params gq.ResolveParams, tableToUpdate string, tableToSelectFrom string) (interface{}, error) {
	id := params.Args[keyFieldName]
	fieldValues, err := db.UpdateRowByID(keyFieldName, tableToUpdate, oldKeyValue, params.Args)
	if err != nil {
		return fieldValues, err
	}
	path := params.Info.FieldName
	fields := getSelectedFields([]string{path}, params)
	return db.QueryRowMap("SELECT "+fields+" FROM \""+tableToSelectFrom+"\" WHERE "+keyFieldName+" = $1 ;", id)

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
	fieldValues, err := db.QueryRowMap("SELECT "+fields+" FROM \""+tableToSelectFrom+"\" WHERE "+keyFieldName+" = $1 ;", id)

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
	return GetSessionVariable(c, "user")
}

func panicIfNotOwnerOrAdmin(params gq.ResolveParams) {
	uname := getLoginedUserName(params)
	pname, ok := params.Args["username"].(string)
	if ok && pname == uname {
		return
	}
	panicIfNotAdmin(params)
}

func panicIfNotOwnerOrAdminOrAuditor(params gq.ResolveParams) {
	uname := getLoginedUserName(params)
	pname, ok := params.Args["username"].(string)
	if ok && pname == uname {
		return
	}
	panicIfNotAdminOrAuditor(params)
}

func isAuthAdmin(params gq.ResolveParams) bool {
	username := getLoginedUserName(params)
	return auth.AppUserRoleExist("auth", username, "authadmin")
}

func isAuditor(params gq.ResolveParams) bool {
	username := getLoginedUserName(params)
	return auth.AppUserRoleExist("auth", username, "auditor")
}

func panicIfNotAdmin(params gq.ResolveParams) {
	if isAuthAdmin(params) {
		return
	}
	panic("Sorry. No admin rights.")
}

func panicIfNotAdminOrAuditor(params gq.ResolveParams) {
	if isAuthAdmin(params) {
		return
	}
	if isAuditor(params) {
		return
	}
	panic("Sorry. No admin or auditor rights.")
}

func panicIfNotUser(params gq.ResolveParams) {
	username := getLoginedUserName(params)
	if username == "" {
		panic("Sorry. You have to login.")
	}
}

func panicIfEmpty(v interface{}, message string) {
	username, _ := v.(string)
	username = strings.Trim(username, " ")
	if len(username) == 0 {
		panic(errors.New("Внимание: " + message))
	}
}

func processPassword(params gq.ResolveParams) string {
	password, _ := params.Args["password"].(string)
	password = strings.Trim(password, " ")

	// remove empty field
	if password == "" {
		delete(params.Args, "password")
		return password
	}

	// check for length
	if len(password) < 6 {
		panic("Password must be 6 or more symbols long")
	}

	// encode
	params.Args["password"] = auth.GetHash(password)
	return password
}

// ArgToLowerCase преобразует значение параметра с именем argName к нижнему регистру
func ArgToLowerCase(params gq.ResolveParams, argName string) {
	v, ok := params.Args[argName]
	if !ok {
		return
	}
	s, ok := v.(string)
	if !ok {
		return
	}
	params.Args[argName] = strings.ToLower(s)
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

	if len(result.Errors) > 0 {
		// инкрементируем счетчик ошибок GraphQL
		prometeo.GraphQLErrorsTotal.Inc()
	}

	c.JSON(200, result)
}
