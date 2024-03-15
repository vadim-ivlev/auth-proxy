package requestthrottler

import (
	"auth-proxy/pkg/db"
	"fmt"
	"time"
)

// RemoveOldPinRequestsDB удаляет старые запросы из таблицы requests
func RemoveOldPinRequestsDB() {

	res, err := db.QueryExec("DELETE FROM requests WHERE pin_request_time < $1", time.Now().Add(-Params.PinRequestsTimeInterval).UnixMilli())
	if err != nil {
		fmt.Println("Error in RemoveOldPinRequestsDB: ", err)
		return
	}
	_ = res
	// rowsAffected, _ := res.RowsAffected()
	// fmt.Printf("RemoveOldPinRequestsDB: %d rows affected\n", rowsAffected)
}

// getPinRequestsNumberDB возвращает количество запросов в RequestTrottler для указанного email из таблицы requests
func getPinRequestsNumberDB(email string) (pinRequestsNumber int64, err error) {
	err = db.DBPool.QueryRow("SELECT COUNT(*) FROM requests WHERE email = $1", email).Scan(&pinRequestsNumber)
	if err != nil {
		fmt.Println("Error in getPinRequestsNumberDB: ", err)
	}
	return
}

// getEarliestPinRequestTimeDB возвращает время самого раннего запроса в Requests для указанного email из таблицы requests
func getEarliestPinRequestTimeDB(email string) (earliestPinRequestTime time.Time, err error) {
	earliestPinRequestTimeInt64 := int64(0)
	err = db.DBPool.QueryRow("SELECT MIN(pin_request_time) FROM requests WHERE email = $1", email).Scan(&earliestPinRequestTimeInt64)
	if err != nil {
		fmt.Println("Error in getEarliestPinRequestTimeDB: ", err)
	}
	earliestPinRequestTime = time.UnixMilli(earliestPinRequestTimeInt64)
	return
}

// AddRequestDB - Добавляет запрос в таблицу requests
func AddRequestDB(email string) {
	_, err := db.QueryExec("INSERT INTO requests (email, pin_request_time) VALUES ($1, $2)", email, time.Now().UnixMilli())
	if err != nil {
		fmt.Println("Error in AddRequestDB: ", err)
	}
}

// TryToAddRequestDB пытается добавить запрос в таблицу requests
// Возвращает количество запросов в Requests для указанного email,
// время, которое нужно подождать, прежде чем можно будет добавить новый запрос
// и флаг, удалось ли добавить новый запрос
func TryToAddRequestDB(email string) (pinRequestsNumber int64, timeToWait time.Duration, ok bool) {
	// Удалить старые запросы
	RemoveOldPinRequestsDB()

	// Получаем количество запросов в Requests для указанного email
	pinRequestsNumber, err := getPinRequestsNumberDB(email)
	if err != nil {
		fmt.Println("Error in TryToAddRequestDB: ", err)
		return
	}

	// Проверяем не превышает ли количество запросов разрешенный максимум для указанного email.
	// Если превышает, то возвращаем время, которое нужно подождать,
	// прежде чем можно будет добавить новый запрос
	if pinRequestsNumber >= Params.MaxPinRequestsPerInterval {
		earliestPinRequestTime, err := getEarliestPinRequestTimeDB(email)
		if err != nil {
			fmt.Println("Error in TryToAddRequestDB: ", err)
			return
		}
		// Время, которое нужно подождать, прежде чем можно будет добавить новый запрос
		timeToWait = Params.PinRequestsTimeInterval - time.Since(earliestPinRequestTime)

		return pinRequestsNumber, timeToWait, false
	}

	// Добавляем запрос в Requests
	AddRequestDB(email)

	// получаем количество запросов в Requests для указанного email
	pinRequestsNumber, err = getPinRequestsNumberDB(email)
	if err != nil {
		fmt.Println("Error in TryToAddRequestDB: ", err)
		return
	}
	return pinRequestsNumber, timeToWait, true
}
