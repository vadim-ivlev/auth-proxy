package app

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type appParams struct {
	// Имя Cookie хранимых на компьютере  пользователя
	CookieName string `env:"cookie_name" envDefault:"auth-proxy"`
	// Имя приложения. Используется для генерации PIN Google authenticator
	AppName string `env:"app_name" envDefault:"auth-proxy-dev"`
	// Использовать https вместо http
	Tls bool `env:"tls" envDefault:"true"`
	// Установить флаг secure на куки браузера. Работает только для https протокола.
	Secure bool `env:"secure" envDefault:"true"`
	// Пользователи могут регистрироваться самостоятельно
	Selfreg bool `env:"selfreg" envDefault:"true"`
	// Нужно ли вводить капчу при входе в систему
	UseCaptcha bool `env:"use_captcha" envDefault:"true"`
	// Нужно ли вводить PIN при входе в систему
	UsePin bool `env:"use_pin" envDefault:"true"`
	// Максимально допустимое число ошибок ввода пароля
	MaxAttempts int64 `env:"max_attempts" envDefault:"5"`
	// Время сброса счетчика ошибок пароля в минутах
	ResetTime int64 `env:"reset_time" envDefault:"60"`

	// url страницы подтверждения email
	ConfirmEmailUrl string `env:"confirm_email_url" envDefault:"https://localhost:4400/confirm-email/"`
	// url куда пренаправляется браузер после подтвержедения email
	EntryPoint string `env:"entry_point" envDefault:"https://www.rg.ru"`
	// адрес почтового сервера SMTP
	SmtpAddress string `env:"smtp_address" envDefault:"localhost:1025"`
	// email от которого посылаются письма пользователям
	From string `env:"from" envDefault:"noreply@rg.ru"`
	// админка сервиса
	AdminUrl string `env:"admin_url" envDefault:"https://auth-admin.vercel.app/?url=https://localhost:4400"`
	// тестовая страница сервиса
	GraphqlTestUrl string `env:"graphql_test_url" envDefault:"https://graphql-test.vercel.app/?end_point=https://localhost:4400/schema&tab_name=auth-proxy4400"`
}

// var EnvParams AppEnvParams
var Params appParams

// ReadEnvConfig reads env file and fill EnvParams with environment variables values
func ReadEnvConfig(fileName string) {
	if err := godotenv.Load(fileName); err != nil {
		log.Println(err.Error())
	}

	if err := env.Parse(&Params); err != nil {
		log.Println(err.Error())
	}
}
