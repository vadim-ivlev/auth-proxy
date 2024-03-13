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

func TestRequestCounter_AddRequest(t *testing.T) {

	// создаем новый RequestThrottler
	// с интервалом времени 1 секунда и максимальным количеством запросов 10
	rc := NewRequestThrottler(1*time.Second, 10)

	// делаем несколько попыток добавить запрос
	// с произвольным интервалом ожидания между попытками

	for i := 0; i < 100; i++ {
		requestsNumber, timeToWait, ok := rc.TryToAddRequest()
		fmt.Printf("requestsNumber=%v, timeToWait=%v ok=%v\n", requestsNumber, timeToWait, ok)
		// произвольный интервал ожидания в промежутке от 0 до 100 миллисекунд
		waitingTime := rand.Intn(100)
		fmt.Printf("Ждем %v миллисекунд\n", waitingTime)
		time.Sleep(time.Duration(waitingTime) * time.Millisecond)
	}
}
