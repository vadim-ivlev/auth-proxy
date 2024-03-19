package requestthrottler

import (
	"auth-proxy/pkg/db"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	db.ReadEnvConfig("../../configs/db.env")
	db.WaitForDbConnection()
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func Test_dbAvailable(t *testing.T) {
	if !db.DbAvailable() {
		t.Errorf("dbAvailable() = false")
	}
}

func Test_SelectCount(t *testing.T) {
	res, err := db.QueryRowMap("SELECT COUNT(*) FROM requests Where email = $1", "aa@bb")
	if err != nil {
		fmt.Println("Error in SelectCount: ", err)
	}
	fmt.Printf("SelectCount: %v\n", res)
}

func printRequestTimes() {
	for i, request := range Requests {
		fmt.Printf("  #%d %v\n", i, request.PinRequestTime)
	}
}

func TestRequestCounter_TryToAddRequest(t *testing.T) {
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
		// произвольный интервал ожидания
		waitingTime := rand.Intn(500)
		fmt.Printf("Ждем %v миллисекунд\n", waitingTime)
		time.Sleep(time.Duration(waitingTime) * time.Millisecond)
	}
}

func printRequestTimesDB() {
	res, err := db.QuerySliceMap("SELECT * FROM requests")
	if err != nil {
		fmt.Println("Error in printRequestTimesDB: ", err)
		return
	}
	for i, request := range res {
		time := time.UnixMilli(request["pin_request_time"].(int64))

		fmt.Printf("  #%d %v\n", i, time)
	}
}

func TestRequestCounter_TryToAddRequestDB(t *testing.T) {
	fmt.Printf("Params: %v\n", Params)

	// устанавливаем интервал времени для запросов
	Params.PinRequestsTimeInterval = 12 * time.Second
	// устанавливаем максимальное количество запросов
	Params.MaxPinRequestsPerInterval = 2

	fmt.Printf("Params: %v\n", Params)

	// делаем несколько попыток добавить запрос
	// с произвольным интервалом ожидания между попытками

	for i := 0; i < 100; i++ {
		requestsNumber, timeToWait, ok := TryToAddRequestDB("aaa@bbb.ccc")
		fmt.Printf("requestsNumber=%v, timeToWait=%v ok=%v\n", requestsNumber, timeToWait, ok)
		if ok {
			printRequestTimesDB()
		}
		// произвольный интервал ожидания
		waitingTime := rand.Intn(500)
		fmt.Printf("Ждем %v миллисекунд\n", waitingTime)
		time.Sleep(time.Duration(waitingTime) * time.Millisecond)
	}
}
