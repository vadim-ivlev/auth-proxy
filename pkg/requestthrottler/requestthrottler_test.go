package requestthrottler

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	m.Run()
}

func printRequestTimes() {
	for i, request := range Requests {
		fmt.Printf("  #%d %v\n", i, request.PinRequestTime)
	}
}

func TestRequestCounter_AddRequest(t *testing.T) {
	// устанавливаем интервал времени для запросов
	Params.PinRequestsTimeInterval = 5 * time.Second
	// устанавливаем максимальное количество запросов
	Params.MaxPinRequestsPerInterval = 3

	// делаем несколько попыток добавить запрос
	// с произвольным интервалом ожидания между попытками

	for i := 0; i < 30; i++ {
		requestsNumber, timeToWait, ok := TryToAddRequest("aaa@bbb.ccc")
		fmt.Printf("requestsNumber=%v, timeToWait=%v ok=%v\n", requestsNumber, timeToWait, ok)
		if ok {
			printRequestTimes()
		}
		// произвольный интервал ожидания в промежутке от 0 до 100 миллисекунд
		waitingTime := rand.Intn(500)
		fmt.Printf("Ждем %v миллисекунд\n", waitingTime)
		time.Sleep(time.Duration(waitingTime) * time.Millisecond)
	}
}
