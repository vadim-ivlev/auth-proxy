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
	"time"
)

// Build версия сборки из gitlab-ci,
// используется флаг со значением переменнай CI_PIPELINE_ID (-ldflags="-X 'main.Build=${CI_PIPELINE_ID}'")
// если не установлено по умолчанию равно development
var Build = "env variables"

func main() {

	fmt.Println("Build number:\t", Build)

	// считать параметры командной строки
	servePort, env, pgconfig, pgParamsFromOS := readCommandLineParams()

	// Считать конфиги и установить параметры
	tls := readConfigsAndSetParams(env, pgconfig, pgParamsFromOS)
	db.PrintConfig()

	// ждем готовности базы данных
	db.WaitForDbConnection()

	// порождаем базу данных если ее нет
	db.CreateDatabaseIfNotExists()

	// печатаем приветствие и запускаем сервер
	printGreetings(servePort, env, tls)
	server.Up(servePort, tls, Build)

	db.DBPool.Close()
	fmt.Println("Bye")
}

// Вспомогательные функции =========================================

// readConfigsAndSetParams читаем конфиги, устанавливаем параметры,
// возвращаем true если требуется соединение по https.
// env - конфигурация {dev|front|prod}
// pgconfig - env файл с параметрами подсоединения к Postgres
// pgParamsFromOS - Брать параметры Postgres из переменных окружения OS.
func readConfigsAndSetParams(env, pgconfig string, pgParamsFromOS bool) bool {
	// Если приказано читать параметры Postgres из операционной системы, читаем из OS.
	// Иначе читаем параметры из файлов. Причем:
	// - читаем файл если если он задан в pgconfig.
	// - В противном случае читаем конфиг соответствующий заданной конфигурации env
	if pgParamsFromOS {
		db.ReadEnvConfig("")
	} else {
		if pgconfig != "" {
			db.ReadEnvConfig(pgconfig)
		} else {
			db.ReadEnvConfig("./configs/db.env." + env)
		}
	}
	// читаем конфиг mail.
	mail.ReadConfig("./configs/mail.yaml", env)
	// читаем шаблоны писем
	mail.ReadMailTemplate("./configs/mail-templates.yaml")
	// читаем конфиг Oauth2
	server.ReadOauth2Config("./configs/oauth2.yaml", env)
	// читаем конфиг signature
	signature.ReadConfig("./configs/signature.yaml", env)

	// читаем конфиг общих параметров приложения.
	app.ReadConfig("./configs/app.yaml", env)
	// устанавливаем параметры пакетов
	tls := app.Params.Tls
	server.SelfRegistrationAllowed = app.Params.Selfreg
	server.SecureCookie = app.Params.Secure
	server.UseCaptcha = app.Params.UseCaptcha
	server.Redirects = app.Params.Redirects
	counter.MAX_ATTEMPTS = app.Params.MaxAttempts
	counter.RESET_TIME = time.Duration(app.Params.ResetTime)

	return tls
}

// readCommandLineParams читает параметры командной строки
func readCommandLineParams() (serverPort, env, pgconfig string, pgParamsFromOS bool) {
	flag.StringVar(&serverPort, "port", "4400", "Запустить приложение на указанном порту.")
	flag.StringVar(&env, "env", "dev", "Окружение. Возможные значения: dev - разработка, front - в докере для фронтэнд разработчиков. prod - продакшн.")
	flag.StringVar(&pgconfig, "pgconfig", "", "Конфигурационный файл Postgres.")
	flag.BoolVar(&pgParamsFromOS, "pg-params-from-os", false, "Брать параметры Postgres из переменных окружения OS.")

	flag.Parse()
	flag.Usage()
	return
}

// printGreetings печатаем приветственное сообщение
func printGreetings(serverPort string, env string, tls bool) {
	protocol := "http"
	if tls {
		protocol = "https"
	}
	database := "Postgres"
	msg := `
	
 █████  ██    ██ ████████ ██   ██
██   ██ ██    ██    ██    ██   ██
███████ ██    ██    ██    ███████
██   ██ ██    ██    ██    ██   ██
██   ██  ██████     ██    ██   ██


Auth-Proxy started. 

━━━━━━━━━━ Some parameters ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Environment: %v
Database:%v 
TLS:%v

━━━━━━━━━━ GraphQL endpoints ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

%v://localhost:%v/schema
%v://localhost:%v/graphql

━━━━━━━━━━ Redirects ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

`
	fmt.Printf(msg, env, database, tls, protocol, serverPort, protocol, serverPort)

	for path, url := range server.Redirects {
		fmt.Printf("%v://localhost:%v%v\t-> %v\n", protocol, serverPort, path, url)
	}

	if env == "dev" || env == "front" {
		fmt.Println("\n━━━━━━━━━━ Login credentials for 'dev' or 'front' evironments ━━━━━━━━")
		fmt.Println("username = admin , password = rosgas2011")
	}
}
