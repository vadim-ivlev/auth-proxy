/*
Пакет requestthrottler предоставляет структуру ,
которая содержит интервал времени и максимальное количество запросов, разрешенных в этом интервале
*/

package requestthrottler

import (
	"time"
)

// RequestThrottler - содержит интервал времени и максимальное количество запросов, разрешенных в этом интервале
type RequestThrottler struct {
	// Интервал времени
	TimeInterval time.Duration
	// Максимальное количество запросов, разрешенных в этом интервале времени
	MaxRequests int64
	// Список запросов в интервале времени.
	// Каждый элемент списка содержит время, когда был добавлен запрос
	Requests []time.Time
}

// NewRequestThrottler создает новую структуру RequestTrottler
func NewRequestThrottler(timeInterval time.Duration, maxRequests int64) *RequestThrottler {
	return &RequestThrottler{
		TimeInterval: timeInterval,
		MaxRequests:  maxRequests,
		Requests:     make([]time.Time, 0),
	}
}

// RemoveOldRequests удаляет старые запросы из RequestTrottler
func (rc *RequestThrottler) RemoveOldRequests() {
	for len(rc.Requests) > 0 && time.Since(rc.Requests[0]) > rc.TimeInterval {
		rc.Requests = rc.Requests[1:]
	}
}

// CheckRequests проверяет, не превышает ли количество запросов в RequestTrottler разрешенного максимума.
// Возращает количество запросов в RequestTrottler,
// время, которое нужно подождать, прежде чем можно будет добавить новый запрос
// и флаг, разрешено ли добавление нового запроса
func (rc *RequestThrottler) CheckRequests() (requestsNumber int64, timeToWait time.Duration, ok bool) {
	// Удалить старые запросы
	rc.RemoveOldRequests()
	// Не превышает ли количество запросов разрешенный максимум ?
	if int64(len(rc.Requests)) >= rc.MaxRequests {
		timeToWait = rc.TimeInterval - time.Since(rc.Requests[0])
		return int64(len(rc.Requests)), timeToWait, false
	}
	return int64(len(rc.Requests)), 0, true
}

// TryToAddRequest пытается добавить запрос в  RequestTrottler.
// Возвращает количество запросов в RequestTrottler,
// время, которое нужно подождать, прежде чем можно будет добавить новый запрос
// и флаг, удалось ли добавить новый запрос
func (rc *RequestThrottler) TryToAddRequest() (requestsNumber int64, timeToWait time.Duration, ok bool) {
	// проверяем, не превышает ли количество запросов разрешенного максимума
	requestsNumber, timeToWait, ok = rc.CheckRequests()
	if !ok {
		return requestsNumber, timeToWait, false
	}
	// Добавляем запрос в список запросов
	rc.Requests = append(rc.Requests, time.Now())
	return int64(len(rc.Requests)), timeToWait, true
}
