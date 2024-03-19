/*
Пакет requestthrottler предоставляет структуру ,
которая содержит интервал времени и максимальное количество запросов, разрешенных в этом интервале

*/

package requestthrottler

import (
	"fmt"
	"time"

	"github.com/caarlos0/env"
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

// init инициализирует параметры ограничителя запросов
func init() {
	// Считать параметры ограничителя запросов из переменных окружения
	if err := env.Parse(&Params); err != nil {
		fmt.Println(err.Error())
	}
	Requests = make([]Request, 0)
}
