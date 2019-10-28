package main

import (
	"auth-proxy/pkg/app"
	"auth-proxy/pkg/counter"
	"auth-proxy/pkg/db"
	"auth-proxy/pkg/mail"
	"auth-proxy/server"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Println("████████████████████████ revision 7 ████████████████████████")
	// считать параметры командной строки
	servePort, env := readCommandLineParams()

	// читаем конфиг Postgres.
	db.ReadConfig("./configs/db.yaml", env)
	// читаем конфиг SQLite.
	db.ReadSQLiteConfig("./configs/sqlite.yaml", env)
	// читаем конфиг mail.
	mail.ReadConfig("./configs/mail.yaml", env)
	// читаем шаблоны писем
	mail.ReadMailTemplate("./configs/mail-templates.yaml")

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
	fmt.Println(app.Params)

	// Ждем готовности базы данных
	db.PrintConfig()
	db.WaitForDbOrExit(10)

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


━━━━━━━━━━ Some parameters ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Auth-Proxy started. 
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
