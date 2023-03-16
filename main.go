package main

import (
	"auth-proxy/pkg/app"
	"auth-proxy/pkg/counter"
	"auth-proxy/pkg/db"
	"auth-proxy/pkg/mail"
	"auth-proxy/pkg/signature"
	"auth-proxy/server"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Build версия сборки из gitlab-ci,
// используется флаг со значением переменнай CI_PIPELINE_ID (-ldflags="-X 'main.Build=${CI_PIPELINE_ID}'")
// если не установлено по умолчанию равно development
var Build = "env variables"

func main() {
	fmt.Println("Build number:\t", Build)
	// считать параметры командной строки
	servePort, config, pgconfig, noIntrospection := readCommandLineParams()
	// Скрываем описания GraphQL если нужно
	server.SchemaInit(noIntrospection)
	// Считать конфиги и установить параметры
	readConfigsAndSetParams(config, pgconfig)
	// ждем готовности базы данных
	db.WaitForDbConnection()
	// порождаем базу данных если ее нет
	db.CreateDatabaseIfNotExists()
	// печатаем приветствие
	printGreetings(servePort, app.Params.Tls)
	// запускаем сервер
	server.Up(servePort, app.Params.Tls, Build)

	// Если server.Up запущено в горутине
	// fmt.Println("Bye")
	// db.DBPool.Close()
}

// Вспомогательные функции =========================================

// readConfigsAndSetParams читаем конфиги
// config - env файл с параметрами приложения
// pgconfig - env файл с параметрами подсоединения к Postgres
func readConfigsAndSetParams( /*env,*/ config, pgconfig string) {
	fmt.Println("\nПараметры конфигурационных файлов.")
	// конфиг базы данных
	db.ReadEnvConfig(pgconfig)
	fmt.Printf("\npgconfig=%v \n\ndb.Params: %s\n\n", pgconfig, app.Serialize(db.Params))
	// конфиг Oauth2
	server.ReadOauth2Config("./configs/oauth2.yaml", "front")
	// конфиг signature
	signature.ReadConfig("./configs/signature.yaml", "front")
	// конфиг общих параметров приложения.
	app.ReadEnvConfig(config)
	fmt.Printf("\nconfig=%v \n\napp.Params: %+v\n\n", config, app.Serialize(app.Params))
	// шаблоны писем (важно запустить после ReadEnvConfig)
	mail.ReadMailTemplate("./configs/mail-templates.yaml")
	// устанавливаем параметры пакетов
	server.SelfRegistrationAllowed = app.Params.Selfreg
	server.SecureCookie = app.Params.Secure
	server.UseCaptcha = app.Params.UseCaptcha
	counter.MAX_ATTEMPTS = app.Params.MaxAttempts
	counter.RESET_TIME = time.Duration(app.Params.ResetTime)
}

// readCommandLineParams читает параметры командной строки
func readCommandLineParams() (serverPort, config, pgconfig string, noIntrospection bool) {
	flag.StringVar(&serverPort, "port", "4400", "Запустить приложение на указанном порту.")
	flag.StringVar(&config, "config", "./configs/app.env", "Конфигурационный файл приложения.")
	flag.StringVar(&pgconfig, "pgconfig", "./configs/db.env", "Конфигурационный файл Postgres.")
	flag.BoolVar(&noIntrospection, "no-introspection", false, "Подавлять интроспекцию GraphQL объектов.")

	flag.Parse()
	flag.Usage()
	return
}

// printGreetings печатаем приветственное сообщение
func printGreetings(serverPort string, tls bool) {
	gitpodurl := strings.Replace(os.Getenv("GITPOD_WORKSPACE_URL"), "://", "://"+serverPort+"-", 1)

	protocol := "http"
	if tls {
		protocol = "https"
	}

	msg := `TLS:%v

━━━━━━━━━━ GraphQL endpoints ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

%v://localhost:%v/schema
%v://localhost:%v/graphql

local test
http://localhost:5000/?url=%v://localhost:4400

GitPod URL
%v

━━━━━━━━━━ URLS ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

`
	fmt.Printf(msg, tls, protocol, serverPort, protocol, serverPort, protocol, gitpodurl)
	fmt.Println("Admin Url - >", app.Params.AdminUrl)
	fmt.Println("Test Mail - > http://localhost:8025/")
	fmt.Printf("\n\n")
}
