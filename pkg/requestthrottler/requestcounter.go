/*
Пакет requestthrottler предоставляет структуру RequestCounter,
которая содержит интервал времени и максимальное количество запросов, разрешенных в этом интервале
*/
// request throttler
package requestthrottler

import (
	"time"
)

// RequestCounter - содержит интервал времени и максимальное количество запросов, разрешенных в этом интервале
type RequestCounter struct {
	// Интервал времени в секундах
	TimeInterval int64
	// Максимальное количество запросов, разрешенных в этом интервале времени
	MaxRequests int64
	// Количество запросов, полученных за последние 60 секунд
	Requests []int64
}

// NewRequestCounter создает новую структуру RequestCounter
func NewRequestCounter(timeInterval int64, maxRequests int64) *RequestCounter {
	return &RequestCounter{
		TimeInterval: timeInterval,
		MaxRequests:  maxRequests,
		Requests:     make([]int64, 0),
	}
}

// AddRequest добавляет запрос в RequestCounter
// Возвращает количество запросов в RequestCounter,
// время, которое нужно подождать, прежде чем можно будет добавить новый запрос
// и флаг, удалось ли добавить новый запрос
func (rc *RequestCounter) AddRequest() (requestsNumber, timeToWait int64, ok bool) {
	// проверяем, не превышает ли количество запросов разрешенного максимума
	requestsNumber, timeToWait, ok = rc.CheckRequests()
	if !ok {
		return requestsNumber, timeToWait, false
	}
	// Добавляем запрос в список запросов
	rc.Requests = append(rc.Requests, time.Now().Unix())
	return int64(len(rc.Requests)), timeToWait, true
}

// RemoveOldRequests удаляет старые запросы из RequestCounter
// и возвращает время, которое нужно подождать, прежде чем можно будет добавить новый запрос
func (rc *RequestCounter) RemoveOldRequests() (timeToWait int64) {
	for len(rc.Requests) > 0 && time.Now().Unix()-rc.Requests[0] > rc.TimeInterval {
		rc.Requests = rc.Requests[1:]
	}
	if len(rc.Requests) > 0 {
		return rc.TimeInterval - (time.Now().Unix() - rc.Requests[0])
	}
	return 0
}

// CheckRequests проверяет, не превышает ли количество запросов разрешенного максимума.
// Возращает количество запросов в RequestCounter,
// время, которое нужно подождать, прежде чем можно будет добавить новый запрос
// и флаг, разрешено ли добавление нового запроса
func (rc *RequestCounter) CheckRequests() (requestsNumber, timeToWait int64, ok bool) {
	// Удалить старые запросы
	timeToWait = rc.RemoveOldRequests()
	// Не превышает ли количество запросов разрешенный максимум ?
	if int64(len(rc.Requests)) >= rc.MaxRequests {
		return int64(len(rc.Requests)), timeToWait, false
	}
	return int64(len(rc.Requests)), timeToWait, true
}
