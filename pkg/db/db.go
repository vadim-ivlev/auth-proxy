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
	//blank import
	_ "github.com/mattn/go-sqlite3"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func getDB() (*sqlx.DB, error) {
	if UsePool {
		return getDBFromPool()
	}

	// Если пул не используется каждый раз коннектимся к БД заново
	if SQLite {
		return sqlx.Open("sqlite3", sqliteParams.Sqlitefile)
	}

	return sqlx.Open("postgres", params.connectStr)
}

func getDBFromPool() (*sqlx.DB, error) {
	if DBPool != nil {
		return DBPool, nil
	}

	var err error
	if SQLite {
		DBPool, err = sqlx.Connect("sqlite3", sqliteParams.Sqlitefile)
		printIf("getDBFromPool(): sqlite3", err)
	} else {
		DBPool, err = sqlx.Connect("postgres", params.connectStr)
		printIf("getDBFromPool(): postgres", err)
	}
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
	conn, err := getDB()
	if err != nil {
		fmt.Println(err)
		return false
	}
	if !UsePool {
		defer conn.Close()
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
	conn, err := getDB()
	panicIf(err)
	if !UsePool {
		defer conn.Close()
	}
	return conn.Exec(sqlText, args...)
}

// mapValuesToStrings фикс драйвера sqlite3. Преобразует []uint8 -> string
func mapValuesToStrings(m map[string]interface{}) map[string]interface{} {
	mm := make(map[string]interface{})
	for k, v := range m {
		bytes, ok := v.([]byte)
		if ok {
			mm[k] = string(bytes)
		} else {
			mm[k] = v
		}
	}
	return mm
}

// QuerySliceMap возвращает результат запроса заданного sqlText, как срез отображений ключ - значение.
// Применяется для запросов SELECT возвращающих набор записей.
func QuerySliceMap(sqlText string, args ...interface{}) ([]map[string]interface{}, error) {
	conn, err := getDB()
	panicIf(err)
	if !UsePool {
		defer conn.Close()
	}

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
		if SQLite {
			results = append(results, mapValuesToStrings(row))
		} else {
			results = append(results, row)
		}
	}

	return results, nil
}

// QueryRowMap возвращает результат запроса заданного sqlText, с возможными параметрами args.
// Применяется для исполнения запросов , INSERT, SELECT.
func QueryRowMap(sqlText string, args ...interface{}) (map[string]interface{}, error) {
	conn, err := getDB()
	panicIf(err)
	if !UsePool {
		defer conn.Close()
	}
	result := make(map[string]interface{})
	err = conn.QueryRowx(sqlText, args...).MapScan(result)
	printIf("QueryRowMap() sqlText="+sqlText, err)
	if SQLite {
		return mapValuesToStrings(result), err
	}
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
		tableName, strings.Join(keys, ", "), strings.Join(dollars, ", "), id)
	log.Println("sqlText=", sqlText, values)
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
			// чистим текст от выражений специфичных для postgresql
			sqlText = removeLinesContaining(sqlText, "postgresql-specific")
			_, err = QueryExec(sqlText)
			if err != nil {
				log.Printf("Query Execution error: %s.\t Error:  %v ", fileName, err)
				continue
			}
			log.Printf("Executed: %s \n", fileName)
		}
	}
}

// removeLinesContaining - удаляет из текста строки содержащие данную подстроку
// Функция используется для чистки SQL текстов от выражений специфичных для Postgresql.
func removeLinesContaining(str, substr string) string {
	res := str
	if SQLite {
		re := regexp.MustCompile("(?m)^.*" + substr + ".*$")
		res = re.ReplaceAllString(str, "")
	}
	return res
}
