package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	//blank import
	_ "github.com/lib/pq"
)

// dbAvailable проверяет, доступна ли база данных
func dbAvailable() bool {
	conn, err := sql.Open("postgres", connectStr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer conn.Close()
	err1 := conn.Ping()
	// _, err1 := conn.Exec("select 1;")
	if err1 != nil {
		fmt.Println(err1.Error())
		return false
	}
	return true
}

// WaitForDbOrExit ожидает доступности базы данных
// делая несколько попыток. Если все попытки неудачны
// завершает программу. Нужна для запуска программы в докерах,
// когда запуск базы данных может быть произойти позже.
func WaitForDbOrExit(attempts int) {
	for i := 0; i < attempts; i++ {
		if dbAvailable() {
			return
		}
		fmt.Println("\nОжидание готовности базы данных...")
		fmt.Printf("Попытка %d/%d. CTRL-C для прерывания.\n", i+1, attempts)
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Не удалось подключиться к базе данных.")
	os.Exit(7777)
}

// QuerySliceMap возвращает результат запроса заданного sqlText, как срез отображений ключ - значение.
// Применяется для запросов SELECT возвращающих набор записей.
func QuerySliceMap(sqlText string, args ...interface{}) ([]map[string]interface{}, error) {
	conn, err := sqlx.Open("postgres", connectStr)
	panicIf(err)
	defer conn.Close()

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
	conn, err := sqlx.Open("postgres", connectStr)
	panicIf(err)
	defer conn.Close()
	result := make(map[string]interface{})
	err = conn.QueryRowx(sqlText, args...).MapScan(result)
	printIf("QueryRowMap() sqlText="+sqlText, err)
	return result, err
}

// CreateRow Вставляет запись в таблицу tableName.
// fieldValues задает имена и значения полей таблицы.
// Возвращает map[string]interface{}, error новой записи таблицы.
func CreateRow(tableName string, fieldValues map[string]interface{}) (map[string]interface{}, error) {
	keys, values, dollars := getKeysAndValues(fieldValues)
	sqlText := fmt.Sprintf(`INSERT INTO "%s" ( %s ) VALUES ( %s ) RETURNING * ;`,
		tableName, strings.Join(keys, ", "), strings.Join(dollars, ", "))
	return QueryRowMap(sqlText, values...)
}

// // GetRowByID возвращает запись в таблице tableName по ее id.
// // Возвращает map[string]interface{} записи таблицы.
// func GetRowByID(tableName string, id int) (map[string]interface{}, error) {
// 	sqlText := "SELECT * FROM " + tableName + " WHERE id = $1 ;"
// 	return QueryRowMap(sqlText, id)
// }

// UpdateRowByID обновляет запись в таблице tableName по ее id.
// map fieldValues задает имена и значения полей таблицы.
// Возвращает map[string]interface{}, error обновленной записи таблицы.
func UpdateRowByID(keyFieldName string, tableName string, id interface{}, fieldValues map[string]interface{}) (map[string]interface{}, error) {
	keys, values, dollars := getKeysAndValues(fieldValues)
	sqlText := fmt.Sprintf(`UPDATE "%s" SET ( %s ) = ( %s ) WHERE `+keyFieldName+` = '%v' RETURNING * ;`,
		tableName, strings.Join(keys, ", "), strings.Join(dollars, ", "), id)
	return QueryRowMap(sqlText, values...)
}

// DeleteRowByID удаляет запись в таблице tableName по ее id.
// Возвращает map[string]interface{}, error удаленной записи таблицы.
func DeleteRowByID(keyFieldName string, tableName string, id interface{}) (map[string]interface{}, error) {
	sqlText := fmt.Sprintf(`DELETE FROM "%s" WHERE `+keyFieldName+` = '%v' RETURNING * ;`, tableName, id)
	return QueryRowMap(sqlText)
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
	keys := []string{}
	values := make([]interface{}, 0)
	dollars := []string{}
	n := 1
	for key, val := range vars {
		vv := SerializeIfArray(val)
		values = append(values, vv)
		keys = append(keys, key)
		dollars = append(dollars, fmt.Sprintf("$%v", n))
		n++
	}
	return keys, values, dollars
}
