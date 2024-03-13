package requestthrottler

import "time"

// Значения по умолчанию
var defaultTimeInterval = 1 * time.Second
var defauiltMaxRequests int64 = 10

// throttlers - хранит ограничители запросов для каждого имени
var throttlers map[string]*RequestThrottler

func init() {
	throttlers = make(map[string]*RequestThrottler)
}

// GetRequestThrottler возвращает ограничитель запросов для указанного имени
func GetRequestThrottler(name string) *RequestThrottler {
	if throttlers[name] == nil {
		throttlers[name] = NewRequestThrottler(defaultTimeInterval, defauiltMaxRequests)
	}
	return throttlers[name]
}
