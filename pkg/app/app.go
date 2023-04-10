/*
Package app хранит значения глобальных переменных приложения
*/

package app

import (
	"auth-proxy/pkg/reqcounter"
	"encoding/json"
	"log"
	"runtime"

	// log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// GetParams returns some app parameters
func GetParams() map[string]interface{} {
	return map[string]interface{}{
		"app_name":                  Params.AppName,
		"selfreg":                   Params.Selfreg,
		"use_captcha":               Params.UseCaptcha,
		"use_pin":                   Params.UsePin,
		"login_not_confirmed_email": Params.LoginNotConfirmedEmail,
		"max_attempts":              Params.MaxAttempts,
		"reset_time":                Params.ResetTime,
		"no_schema":                 Params.NoSchema,
	}
}

// GetStat returns some stat about running app
func GetStat() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	alloc := m.Alloc / 1024 / 1024
	totalAlloc := m.TotalAlloc / 1024 / 1024
	sys := m.Sys / 1024 / 1024
	day, hour, min, sec := reqcounter.GetCounters()
	return map[string]interface{}{
		"alloc":               alloc,
		"total_alloc":         totalAlloc,
		"sys":                 sys,
		"requests_per_day":    day,
		"requests_per_hour":   hour,
		"requests_per_minute": min,
		"requests_per_second": sec,
	}
}

// Stat A controller wrapper for REST
func Stat(c *gin.Context) {
	c.JSON(200, GetStat())
}

// LogMessage message controller
func LogMessage(c *gin.Context) {
	message := c.Param("message")
	log.Println("LogMessage", message)
	c.String(200, "Logger test. Search log file or Elastic for: %s", message)
}

func Serialize(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
