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
	"strconv"
	"time"
)

func main() {
	fmt.Println("████████████████████████ revision 8 ████████████████████████")
	// считать параметры командной строки
	servePort, env := readCommandLineParams()

	tls := readConfigsAndSetParams(env)

	fmt.Println(app.Params)

	// Ждем готовности базы данных
	db.PrintConfig()
	db.WaitForDbOrExit(20)

	// порождаем базу данных если ее нет
	db.CreateDatabaseIfNotExists()

	// если servePort > 0, печатаем приветствие и запускаем сервер
	if servePort > 0 {
		printGreetings(servePort, env, app.Params.Sqlite, tls)
		server.Serve(":"+strconv.Itoa(servePort), tls)
	}

	db.DBPool.Close()
	fmt.Println("Bye")
}

// Вспомогательные функции =========================================

// readConfigsAndSetParams читаем конфиги, устанавливаем параметры,
// возвращаем true если требуется соединение по https.
func readConfigsAndSetParams(env string) bool {
	// читаем конфиг Postgres.
	db.ReadConfig("./configs/db.yaml", env)
	// читаем конфиг SQLite.
	db.ReadSQLiteConfig("./configs/sqlite.yaml", env)
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
	db.SQLite = app.Params.Sqlite
	server.SelfRegistrationAllowed = app.Params.Selfreg
	server.SecureCookie = app.Params.Secure
	server.UseCaptcha = app.Params.UseCaptcha
	server.Redirects = app.Params.Redirects
	counter.MAX_ATTEMPTS = app.Params.MaxAttempts
	counter.RESET_TIME = time.Duration(app.Params.ResetTime)

	return tls
}

// readCommandLineParams читает параметры командной строки
func readCommandLineParams() (serverPort int, env string) {
	flag.IntVar(&serverPort, "serve", 0, "Запустить приложение на порту с номером > 0 ")
	flag.StringVar(&env, "env", "prod", "Окружение. Возможные значения: dev - разработка, front - в докере для фронтэнд разработчиков. prod - по умолчанию для продакшн.")

	flag.Parse()
	fmt.Println("\nПример запуска: ./auth-proxy -serve 4400 -env=dev\n ")
	flag.Usage()
	if serverPort == 0 {
		os.Exit(0)
	}
	return
}

// printGreetings печатаем приветственное сообщение
func printGreetings(serverPort int, env string, sqlite bool, tls bool) {
	protocol := "http"
	if tls {
		protocol = "https"
	}

	database := "Postgres"
	if sqlite {
		database = "SQLite"
	}

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

	fmt.Println("\n\nCTRL-C to interrupt.\n")
}
