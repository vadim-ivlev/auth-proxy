/*
Package app хранит значения глобальных переменных приложения
*/

package app

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type appParams struct {
	// Использовать SQLite
	Sqlite bool
	// Использовать https вместо http
	Tls bool
	// Установить флаг secure на куки браузера. Работает только для https протокола.
	Secure bool
	// Пользователи могут регистрироваться самостоятельно
	Selfreg bool
	// Нужно ли вводить капчу при входе в систему
	UseCaptcha bool `yaml:"use_captcha"`
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
