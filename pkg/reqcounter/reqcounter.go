/*
Package reqcounter Cчетчики запросов за час, минуту, секунду.

*/

package reqcounter

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var theCache = cache.New(24*time.Hour, 48*time.Hour)

// IncrementCounters инкрементирует счетчики
func IncrementCounters() {
	addIncCounter("day", 24*time.Hour)
	addIncCounter("hour", time.Hour)
	addIncCounter("min", time.Minute)
	addIncCounter("sec", time.Second)
}

// GetCounters возвращает значения счетчиков
func GetCounters() (int64, int64, int64, int64) {
	return getCounter("day"),
		getCounter("hour"),
		getCounter("min"),
		getCounter("sec")
}

func addIncCounter(name string, duration time.Duration) {
	_ = theCache.Add(name, int64(0), duration)
	// if err != nil {
	// 	log.Println("reqcounter.IncrementCounters.Add "+name+":", err)
	// }
	_ = theCache.Increment(name, 1)
	// if err != nil {
	// 	log.Println("counter.IncrementCounter.Increment"+name+":", err)
	// }
}

func getCounter(name string) int64 {
	n, ok := theCache.Get(name)
	if !ok {
		return 0
	}
	return n.(int64)
}
