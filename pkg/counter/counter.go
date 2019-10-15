/*
Package counter хранит счетчики неудачных попыток пользователей войти в систему.

- После неудачного входа счетчик инкрементируется.
- После удачного входа счетчик сбрасывается.
- После истечения определенного времени счетчик сбрасывается.
*/

package counter

import (
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

// MAX_ATTEMPTS - максимальное количество неудачных попыток входа
var MAX_ATTEMPTS int64 = 5

// RESET_TIME - время по истечении которого счетчик сбрасывается
var RESET_TIME time.Duration = 60

//
var theCache = cache.New(RESET_TIME*time.Minute, RESET_TIME*2*time.Minute)

// IncrementCounter инкрементирует счетчик пользователя
func IncrementCounter(username string) {
	err := theCache.Add(username, int64(0), cache.DefaultExpiration)
	if err != nil {
		log.Println("counter.IncrementCounter.Add:", err)
	}
	err = theCache.Increment(username, 1)
	if err != nil {
		log.Println("counter.IncrementCounter.Increment:", err)
		return
	}
}

// ResetCounter сбрасывает счетчик пользователя
func ResetCounter(username string) {
	theCache.Delete(username)
}

// IsCounterTooBig возвращает true если разрешенное количество неудачных попыток еще не достигнуто
func IsTooBig(username string) bool {
	n, ok := theCache.Get(username)
	if !ok {
		return false
	}
	if n.(int64) >= MAX_ATTEMPTS {
		return true
	}
	return false
}
