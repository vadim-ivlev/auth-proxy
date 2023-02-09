package server

import (
	"auth-proxy/pkg/auth"
	"errors"
	"fmt"
	"log"
	"strconv"

	// "go/ast"

	"strings"

	"auth-proxy/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
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
	return mb, errors.New("no body")
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
func createRecord(keyFieldName string, params graphql.ResolveParams, tableToUpdate string, tableToSelectFrom string) (interface{}, error) {
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
func updateRecord(oldKeyValue interface{}, keyFieldName string, params graphql.ResolveParams, tableToUpdate string, tableToSelectFrom string) (interface{}, error) {
	id := params.Args[keyFieldName]
	fieldValues, err := db.UpdateRowByID(keyFieldName, tableToUpdate, oldKeyValue, params.Args)
	if err != nil {
		return fieldValues, err
	}
	path := params.Info.FieldName
	fields := getSelectedFields([]string{path}, params)
	return db.QueryRowMap("SELECT "+fields+" FROM \""+tableToSelectFrom+"\" WHERE \""+db.RemoveDoubleQuotes(keyFieldName)+"\" = $1 ;", id)

}

// deleteRecord удаляет запись из таблицы tableToUpdate,
// и возвращает удаленную запись из таблицы tableToSelectFrom,
// которая является представлением с более богатым содержимым,
// чем обновленная таблица.
// Используется в запросах GraphQL на удаление записей.
func deleteRecord(keyFieldName string, params graphql.ResolveParams, tableToUpdate string, tableToSelectFrom string) (interface{}, error) {
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
func getSelectedFields(selectionPath []string, resolveParams graphql.ResolveParams) string {
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

func getLoginedUserName(params graphql.ResolveParams) string {
	c, ok := params.Context.Value("ginContext").(*gin.Context)
	if !ok {
		log.Println("Not OK: getLoginedUserName")
		return ""
	}
	return GetSessionVariable(c, "user")
}
func getLoginedUserID(params graphql.ResolveParams) string {
	c, ok := params.Context.Value("ginContext").(*gin.Context)
	if !ok {
		log.Println("Not OK: getLoginedUserID")
		return ""
	}
	return GetSessionVariable(c, "id")
}

func getLoginedUserEmail(params graphql.ResolveParams) string {
	c, ok := params.Context.Value("ginContext").(*gin.Context)
	if !ok {
		log.Println("Not OK: getLoginedUserEmail")
		return ""
	}
	return GetSessionVariable(c, "email")
}

func panicIfNotOwnerOrAdmin(params graphql.ResolveParams) {
	uName := getLoginedUserName(params)
	pName, ok := params.Args["username"].(string)
	if ok && pName == uName {
		return
	}
	uID := getLoginedUserID(params)
	pID := strconv.Itoa(params.Args["id"].(int))
	if pID == uID {
		return
	}
	panicIfNotAdmin(params)
}

func panicIfNotOwnerOrAdminOrAuditor(params graphql.ResolveParams) {
	uname := getLoginedUserName(params)
	pname, ok := params.Args["username"].(string)
	if ok && pname == uname {
		return
	}
	panicIfNotAdminOrAuditor(params)
}

func isAuthAdmin(params graphql.ResolveParams) bool {
	username := getLoginedUserName(params)
	return auth.AppUserRoleExist("auth", username, "authadmin")
}

func isAuditor(params graphql.ResolveParams) bool {
	username := getLoginedUserName(params)
	return auth.AppUserRoleExist("auth", username, "auditor")
}

func panicIfNotAdmin(params graphql.ResolveParams) {
	if isAuthAdmin(params) {
		return
	}
	panic("Sorry. No admin rights.")
}

func panicIfNotAdminOrAuditor(params graphql.ResolveParams) {
	if isAuthAdmin(params) {
		return
	}
	if isAuditor(params) {
		return
	}
	panic("Sorry. No admin or auditor rights.")
}

func panicIfNotUser(params graphql.ResolveParams) {
	username := getLoginedUserName(params)
	if username == "" {
		panic("Войдите в приложение")
	}
}

func panicIfEmpty(v interface{}, message string) {
	username, _ := v.(string)
	username = strings.TrimSpace(username)
	if len(username) == 0 {
		panic(errors.New("Внимание: " + message))
	}
}

func convertPasswordToHash(params graphql.ResolveParams) string {
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
func ArgToLowerCase(params graphql.ResolveParams, argName string) {
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

// TrimParamValue удаляет ведущие и последние пробелы
// из значения параметра с именем argName.
func TrimParamValue(params graphql.ResolveParams, argName string) {
	v, ok := params.Args[argName]
	if !ok {
		return
	}
	s, ok := v.(string)
	if !ok {
		return
	}
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		panic(errors.New("пустое значение недопустимо"))
	}
	params.Args[argName] = s
}

// getGroupByName возвращает группу с именем groupName.
// Возвращает группу и ошибку, если группа не найдена.
func getGroupByName(groupName string) (interface{}, error) {
	return db.QueryRowMap("SELECT * FROM \"group\" WHERE groupname = $1 ;", groupName)
}

// createNewGroup создает новую группу
// c именем groupName и описанием groupDescription.
// Возвращает созданную группу и ошибку, если группа с таким именем уже существует.
func createNewGroup(groupName, groupDescription string) (interface{}, error) {
	params := graphql.ResolveParams{}
	params.Args = make(map[string]interface{})
	params.Args["name"] = groupName
	ArgToLowerCase(params, "name")
	TrimParamValue(params, "name")
	params.Args["description"] = groupDescription
	return createRecord("name", params, "group", "group")
}

// createGroupUserRole добавляет пользователя в группу,
// создавая запись в таблице group_user_role.
// Возвращает ошибку, если пользователь уже состоит в группе.
func createGroupUserRole(groupID, userID int, rolename string) error {
	_, err := db.CreateRow("group_user_role", map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
		"rolename": rolename,
	})
	if err != nil {
		return err
	}
	clearCache()
	return nil
}

// addUserToDefaultGroup добавляет пользователя с идеинтификатором userID
// в группу для пользователей по умолчанию.
func addUserToDefaultGroup(userID int) error {
	groupName := "users"
	groupDescription := "Группа для новых пользователей"

	group, err := getGroupByName(groupName)
	if err != nil {
		fmt.Println("addUserToGroup getGroupByName() Не удалось найти группу users.: ", err)
		fmt.Println(`Создаем новую группу "users".`)
		group, err = createNewGroup(groupName, groupDescription)
		if err != nil {
			fmt.Println(`Не удалось создать новую группу "users".`, err)
			return err
		}
	}
	groupID, _ := group.(map[string]interface{})["id"].(int)
	return createGroupUserRole(groupID, userID, "new_user")
}
