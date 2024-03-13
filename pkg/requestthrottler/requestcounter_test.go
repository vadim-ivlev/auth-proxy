package requestthrottler

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestRequestCounter_AddRequest(t *testing.T) {
	// Создаем новый RequestCounter
	// с интервалом времени 11 секунд и максимальным количеством запросов 3
	// и добавляем два запроса 10 и 5 секунд назад
	rc := &RequestCounter{
		TimeInterval: 11,
		MaxRequests:  3,
		Requests:     []int64{time.Now().Unix() - 10, time.Now().Unix() - 5},
	}

	// делаем десять попыток добавить запрос
	// с интервалом ожидания
	var interval int64 = 1

	for i := 0; i < 10; i++ {
		requestsNumber, timeToWait, ok := rc.AddRequest()
		fmt.Printf("requestsNumber=%d, timeToWait=%d ok=%v\n", requestsNumber, timeToWait, ok)
		fmt.Println("Ждем")
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
