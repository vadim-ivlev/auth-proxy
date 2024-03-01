package app

import (
	"log"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type appParams struct {
	// Использовать кэш для ускорения запросов
	UseCache bool `json:"use_cache" env:"use_cache" envDefault:"true"`
	// Имя Cookie хранимых на компьютере  пользователя
	CookieName string `json:"cookie_name" env:"cookie_name" envDefault:"auth-proxy"`
	// Имя приложения. Используется для генерации PIN Google authenticator
	AppName string `json:"app_name" env:"app_name" envDefault:"auth-proxy-dev"`
	// Использовать https вместо http
	Tls bool `json:"tls" env:"tls" envDefault:"true"`
	// Установить флаг secure на куки браузера. Работает только для https протокола.
	Secure bool `json:"secure" env:"secure" envDefault:"true"`
	// Пользователи могут регистрироваться самостоятельно
	Selfreg bool `json:"selfreg" env:"selfreg" envDefault:"true"`
	// Нужно ли вводить капчу при входе в систему
	UseCaptcha bool `json:"use_captcha" env:"use_captcha" envDefault:"true"`
	// Нужно ли вводить PIN при входе в систему
	UsePin bool `json:"use_pin" env:"use_pin" envDefault:"true"`
	// Максимально допустимое число ошибок ввода пароля
	MaxAttempts int64 `json:"max_attempts" env:"max_attempts" envDefault:"5"`
	// Время сброса счетчика ошибок пароля в минутах
	ResetTime int64 `json:"reset_time" env:"reset_time" envDefault:"60"`
	// Разрешить авторизацию пользователей не подтвердивших email
	LoginNotConfirmedEmail bool `json:"login_not_confirmed_email" env:"login_not_confirmed_email" envDefault:"true"`
	// Адрес API
	AdminAPI string `json:"admin_api" env:"admin_api" envDefault:"https://localhost:4400"`
	// url куда пренаправляется браузер после подтвержедения email
	EntryPoint string `json:"entry_point" env:"entry_point" envDefault:"https://rg.ru/account/profile?email_success=true"`
	// адрес почтового сервера SMTP
	SmtpAddress string `json:"smtp_address" env:"smtp_address" envDefault:"localhost:1025"`
	// email от которого посылаются письма пользователям
	From string `json:"from" env:"from" envDefault:"noreply@rg.ru"`
	// админка сервиса
	AdminUrl string `json:"admin_url" env:"admin_url" envDefault:"https://localhost:4400/admin/?url=https://localhost:4400"`
	// host сайта
	SiteHost string `json:"site_host" env:"site_host"`

	MailTmplPath string `json:"mail_tmpl_path" env:"mail_tmpl_path" envDefault:"./templates/mail"`
	// подавлять чтение схемы GraphQL
	NoSchema bool `json:"no_schema" env:"no_schema" envDefault:"false"`
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

	// подправляем HTTP протоколы в соответствие с tls
	s0 := "http://localhost"
	s1 := "https://localhost"
	if Params.Tls {
		// Params.ConfirmEmailUrl = strings.Replace(Params.ConfirmEmailUrl, s0, s1, -1)
		Params.AdminUrl = strings.Replace(Params.AdminUrl, s0, s1, -1)
	} else {
		// Params.ConfirmEmailUrl = strings.Replace(Params.ConfirmEmailUrl, s1, s0, -1)
		Params.AdminUrl = strings.Replace(Params.AdminUrl, s1, s0, -1)
	}
}
