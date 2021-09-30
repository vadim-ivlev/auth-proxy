/*
Package app хранит значения глобальных переменных приложения
*/

package app

import (
	"auth-proxy/pkg/reqcounter"
	"io/ioutil"
	"log"
	"runtime"

	// log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	yaml "gopkg.in/yaml.v2"
)

type appParams struct {
	// Имя приложения. Используется для генерации PIN Google authenticator
	AppName string `yaml:"app_name"`
	// Использовать https вместо http
	Tls bool
	// Установить флаг secure на куки браузера. Работает только для https протокола.
	Secure bool
	// Пользователи могут регистрироваться самостоятельно
	Selfreg bool
	// Нужно ли вводить капчу при входе в систему
	UseCaptcha bool `yaml:"use_captcha"`
	// Нужно ли вводить PIN при входе в систему
	UsePin bool `yaml:"use_pin"`
	// Максимально допустимое число ошибок ввода пароля
	MaxAttempts int64 `yaml:"max_attempts"`
	// Время сброса счетчика ошибок пароля в минутах
	ResetTime int64 `yaml:"reset_time"`
	// перенаправления браузера для предоставления различных GUI
	Redirects map[string]string
}

// Общие параметры приложения
var Params appParams

// ReadConfig reads YAML with Postgres params
func ReadConfig(fileName string, env string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}

	envParams := make(map[string]appParams)
	err = yaml.Unmarshal(yamlFile, &envParams)
	if err != nil {
		panic(err.Error())
	}
	Params = envParams[env]
}

// GetParams returns some app parameters
func GetParams() map[string]interface{} {
	return map[string]interface{}{
		"app_name":     Params.AppName,
		"selfreg":      Params.Selfreg,
		"use_captcha":  Params.UseCaptcha,
		"use_pin":      Params.UsePin,
		"max_attempts": Params.MaxAttempts,
		"reset_time":   Params.ResetTime,
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
