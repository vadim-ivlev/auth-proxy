/*
Пакет requestthrottler предоставляет структуру ,
которая содержит интервал времени и максимальное количество запросов, разрешенных в этом интервале

TODO: переместисть Requests в таблицу базы данных.
*/

package requestthrottler

import (
	"time"
)

type throttlerParams struct {
	// Интервал времени
	PinRequestsTimeInterval time.Duration `json:"pin_requests_time_interval" env:"pin_requests_time_interval" envDefault:"1m"`
	// Максимальное количество запросов, разрешенных в этом интервале времени
	MaxPinRequestsPerInterval int64 `json:"max_pin_requests_per_interval" env:"max_pin_requests_per_interval" envDefault:"10"`
}

// Request - структура, содержащая имя пользователя и время последнего запроса
type Request struct {
	// Email пользователя
	Email string
	// Время последнего запроса
	PinRequestTime time.Time
}

// throttlerParams - параметры ограничителя запросов
var Params throttlerParams

// Requests - список запросов пользователей
var Requests []Request

// init инициализирует параметры ограничителя запросов
func init() {
	Params = throttlerParams{
		PinRequestsTimeInterval:   2 * time.Minute,
		MaxPinRequestsPerInterval: 5,
	}
	Requests = make([]Request, 0)
}

// RemoveOldPinRequests удаляет старые запросы из RequestTrottler
func RemoveOldPinRequests() {
	for len(Requests) > 0 && time.Since(Requests[0].PinRequestTime) > Params.PinRequestsTimeInterval {
		Requests = Requests[1:]
	}
}

// getPinRequestsNumber возвращает количество запросов в RequestTrottler для указанного email
func getPinRequestsNumber(email string) int64 {
	var pinRequestsNumber int64
	for _, request := range Requests {
		if request.Email == email {
			pinRequestsNumber++
		}
	}
	return pinRequestsNumber
}

// getEarliestPinRequestTime возвращает время самого раннего запроса в Requests для указанного email
func getEarliestPinRequestTime(email string) (earliestPinRequestTime time.Time) {
	// Найдем минимальное значение RequestTime для указанного email

	earliestPinRequestTime = time.Now()
	for _, request := range Requests {
		if request.Email == email && request.PinRequestTime.Before(earliestPinRequestTime) {
			earliestPinRequestTime = request.PinRequestTime
		}
	}
	return
}

// AddRequest - Добавляет запрос в Requests
func AddRequest(email string) {
	Requests = append(Requests, Request{Email: email, PinRequestTime: time.Now()})
}

// TryToAddRequest пытается добавить запрос в  Requests.
// Возвращает количество запросов в Requests для указанного email,
// время, которое нужно подождать, прежде чем можно будет добавить новый запрос
// и флаг, удалось ли добавить новый запрос
func TryToAddRequest(email string) (pinRequestsNumber int64, timeToWait time.Duration, ok bool) {
	// Удалить старые запросы
	RemoveOldPinRequests()

	// Получаем количество запросов в Requests для указанного email
	pinRequestsNumber = getPinRequestsNumber(email)

	// Проверяем не превышает ли количество запросов разрешенный максимум для указанного email.
	// Если превышает, то возвращаем время, которое нужно подождать,
	// прежде чем можно будет добавить новый запрос
	if pinRequestsNumber >= Params.MaxPinRequestsPerInterval {
		minRequestTime := getEarliestPinRequestTime(email)
		// Время, которое нужно подождать, прежде чем можно будет добавить новый запрос
		timeToWait = Params.PinRequestsTimeInterval - time.Since(minRequestTime)

		return pinRequestsNumber, timeToWait, false
	}

	// Добавляем запрос в Requests
	AddRequest(email)

	// получаем количество запросов в Requests для указанного email
	pinRequestsNumber = getPinRequestsNumber(email)

	return pinRequestsNumber, 0, true
}
