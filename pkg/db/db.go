package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"

	//blank import
	_ "github.com/lib/pq"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func getDBFromPool() (*sqlx.DB, error) {
	if DBPool != nil {
		return DBPool, nil
	}

	var err error
	DBPool, err = sqlx.Connect("postgres", Params.connectStr)
	printIf("getDBFromPool(): postgres", err)
	if err != nil {
		return nil, err
	}
	DBPool.SetMaxOpenConns(1000) // The default is 0 (unlimited)
	DBPool.SetMaxIdleConns(4)    // defaultMaxIdleConns = 2
	DBPool.SetConnMaxLifetime(0) // 0, connections are reused forever.
	return DBPool, err
}

// DbAvailable проверяет, доступна ли база данных
func DbAvailable() bool {
	conn, err := getDBFromPool()
	if err != nil {
		fmt.Println(err)
		return false
	}
	err = conn.Ping()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// QueryExec исполняет запросы заданные в sqlText.
func QueryExec(sqlText string, args ...interface{}) (sql.Result, error) {
	conn, err := getDBFromPool()
	panicIf(err)
	// log.Println("QueryExec SQL=", sqlText, args)
	return conn.Exec(sqlText, args...)
}

// QuerySliceMap возвращает результат запроса заданного sqlText, как срез отображений ключ - значение.
// Применяется для запросов SELECT возвращающих набор записей.
func QuerySliceMap(sqlText string, args ...interface{}) ([]map[string]interface{}, error) {
	conn, err := getDBFromPool()
	panicIf(err)
	// log.Println("QuerySliceMap SQL=", sqlText, args)
	rows, err := conn.Queryx(sqlText, args...) //.MapScan(result)

	if err != nil {
		fmt.Println("QuerySliceMap():", err.Error())
		return nil, err
	}

	results := []map[string]interface{}{}

	for rows.Next() {
		row := make(map[string]interface{})
		err = rows.MapScan(row)
		if err != nil {
			log.Println("QuerySliceMap(): ", err)
		}
		results = append(results, row)

	}

	return results, nil
}

// QueryRowMap возвращает результат запроса заданного sqlText, с возможными параметрами args.
// Применяется для исполнения запросов , INSERT, SELECT.
func QueryRowMap(sqlText string, args ...interface{}) (map[string]interface{}, error) {
	conn, err := getDBFromPool()
	panicIf(err)
	result := make(map[string]interface{})
	// log.Println("QueryRowMap SQL=", sqlText, args)
	err = conn.QueryRowx(sqlText, args...).MapScan(result)
	printIf("QueryRowMap() response=", err)
	return result, err
}

// CreateRow Вставляет запись в таблицу tableName.
// fieldValues задает имена и значения полей таблицы.
// Возвращает map[string]interface{}, error новой записи таблицы.
func CreateRow(tableName string, fieldValues map[string]interface{}) (map[string]interface{}, error) {
	keys, values, dollars := getKeysAndValues(fieldValues)

	sqlText := fmt.Sprintf(`INSERT INTO "%s" ( %s ) VALUES ( %s ) ;`,
		tableName, strings.Join(keys, ", "), strings.Join(dollars, ", "))
	res, err := QueryExec(sqlText, values...)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	return map[string]interface{}{"RowsAffected": n}, nil
}

// UpdateRowByID обновляет запись в таблице tableName по ее id.
// map fieldValues задает имена и значения полей таблицы.
// Возвращает map[string]interface{}, error обновленной записи таблицы.
func UpdateRowByID(keyFieldName string, tableName string, id interface{}, fieldValues map[string]interface{}) (map[string]interface{}, error) {
	keys, values, dollars := getKeysAndValues(fieldValues)

	sqlText := fmt.Sprintf(`UPDATE "%s" SET ( %s ) = ( %s ) WHERE `+keyFieldName+` = '%v';`,
		RemoveDoubleQuotesStr(tableName), strings.Join(keys, ", "), strings.Join(dollars, ", "), RemoveSingleQuotes(id))
	res, err := QueryExec(sqlText, values...)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	return map[string]interface{}{"RowsAffected": n}, nil
}

// DeleteRowByID удаляет запись в таблице tableName по ее id.
// Возвращает map[string]interface{}, error удаленной записи таблицы.
func DeleteRowByID(keyFieldName string, tableName string, id interface{}) (map[string]interface{}, error) {

	sqlText := fmt.Sprintf(`DELETE FROM "%s" WHERE `+keyFieldName+` = '%v' ;`, tableName, id)
	res, err := QueryExec(sqlText)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	return map[string]interface{}{"RowsAffected": n}, nil

}

// SerializeIfArray - возвращает JSON строку, если входной параметр массив.
// в противном случае не изменяет значения.
// Используется при отдаче данных пользователю в GraphQL.
func SerializeIfArray(val interface{}) interface{} {
	switch vv := val.(type) {
	case []interface{}:
		jsonBytes, _ := json.Marshal(val)
		jsonString := string(jsonBytes)
		return jsonString
	default:
		return vv
	}
}

// ToPostgresArrayLiteral - преобразует массив в строковое выражение,
// для INSERT, UPDATE запросов к Postgres.
func ToPostgresArrayLiteral(arr interface{}) string {
	tags, ok := arr.([]interface{})
	if ok {
		var ss []string
		for _, v := range tags {
			ss = append(ss, v.(string))
		}
		return "{" + strings.Join(ss, ",") + "}"
	}
	return "{}"
}

// getKeysAndValues возвращает срезы ключей, значений и символов доллара $n.
func getKeysAndValues(vars map[string]interface{}) ([]string, []interface{}, []string) {
	ln := len(vars)
	keys := make([]string, ln)
	values := make([]interface{}, ln)
	dollars := make([]string, ln)
	var n int64 = 0
	for key, val := range vars {
		values[n] = SerializeIfArray(val)
		keys[n] = key
		dollars[n] = "$" + strconv.FormatInt(n+1, 10)
		n++
	}
	return keys, values, dollars
}

// CreateDatabaseIfNotExists порождает объекты базы данных и наполняет базу тестовыми данными
func CreateDatabaseIfNotExists() {
	fmt.Println("Миграция ...")
	MigrateUp("./migrations/")
}

// MigrateUp миграция БД
func MigrateUp(dirname string) {
	files, err := ioutil.ReadDir(dirname)
	panicIf(err)

	for _, file := range files {
		fileName := file.Name()

		if strings.HasSuffix(fileName, "up.sql") {
			sqlBytes, err := ioutil.ReadFile(dirname + fileName)
			if err != nil {
				log.Println("Cannot read file: ", fileName)
				continue
			}
			sqlText := string(sqlBytes)
			_, err = QueryExec(sqlText)
			if err != nil {
				log.Printf("Query Execution error: %s.\t Error:  %v ", fileName, err)
				continue
			}
			log.Printf("Executed: %s \n", fileName)
		}
	}
}

// RemoveSingleQuotes - sanitizes SQL
func RemoveSingleQuotes(text interface{}) interface{} {
	// проверяем строка ли это
	s, ok := text.(string)
	if !ok {
		return text
	}
	return strings.Replace(s, "'", "", -1)
}

// RemoveSingleQuotesStr - sanitizes SQL
func RemoveSingleQuotesStr(text string) string {
	return strings.ReplaceAll(text, "'", "")
}

// RemoveDoubleQuotesStr - - sanitizes SQL
func RemoveDoubleQuotesStr(text string) string {
	return strings.Replace(text, "\"", "", -1)
}

// SanitizeOrderClause - sanitizes SQL ORDER BY clause.
// допускает выражения типа field1 ASC
func SanitizeOrderClause(text string) string {
	s := strings.Trim(text, " ")
	validExpr := regexp.MustCompile(`(?i)^\w+ +(ASC|DESC)$`)
	if validExpr.MatchString(s) {
		return s
	}
	return ""
}
