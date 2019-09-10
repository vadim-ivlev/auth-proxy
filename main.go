package main

import (
	"auth-proxy/model/db"
	"auth-proxy/model/mail"
	"auth-proxy/router"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	fmt.Println("████████████████████████ revision 2 ████████████████████████")
	// считать параметры командной строки
	servePort, env, sqlite, tls := readCommandLineParams()
	db.SQLite = sqlite

	// читаем конфиг Postgres.
	db.ReadConfig("./configs/db.yaml", env)
	// читаем конфиг SQLite.
	db.ReadSQLiteConfig("./configs/sqlite.yaml", env)
	// читаем конфиг mail.
	mail.ReadConfig("./configs/mail.yaml", env)

	// Ждем готовности базы данных
	db.PrintConfig()
	db.WaitForDbOrExit(10)

	// порождаем базу данных если ее нет
	db.CreateDatabaseIfNotExists()

	// если servePort > 0, печатаем приветствие и запускаем сервер
	if servePort > 0 {
		printGreetings(servePort, env, sqlite, tls)
		router.Serve(":"+strconv.Itoa(servePort), tls)
	}

	db.DBPool.Close()
	fmt.Println("Bye")
}

// Вспомогательные функции =========================================

// readCommandLineParams читает параметры командной строки
func readCommandLineParams() (serverPort int, env string, sqlite bool, tls bool) {
	flag.IntVar(&serverPort, "serve", 0, "Запустить приложение на порту с номером > 0 ")
	flag.StringVar(&env, "env", "prod", "Окружение. Возможные значения: dev - разработка, front - в докере для фронтэнд разработчиков. prod - по умолчанию для продакшн.")
	flag.BoolVar(&sqlite, "sqlite", false, "Использовать SQLite")
	flag.BoolVar(&tls, "tls", false, "Использовать https вместо http")
	flag.Parse()
	fmt.Println("\nПример запуска: ./auth-proxy -serve 4000 -env=dev\n ")
	flag.Usage()
	if serverPort == 0 {
		os.Exit(0)
	}
	return
}

// printGreetings печатаем приветственное сообщение
func printGreetings(serverPort int, env string, sqlite bool, tls bool) {
	greetings, _ := ioutil.ReadFile("./templates/greetings.txt")
	fmt.Printf(string(greetings))
	protocol := "http"
	if tls {
		protocol = "https"
	}

	database := "Postgres"
	if sqlite {
		database = "SQLite"
	}

	msg := `
Auth-Proxy started. 
Environment: %v
Database:%v 

%v://localhost:%v/testapp
		
CTRL-C to interrupt.
`
	fmt.Printf(msg, env, database, protocol, serverPort)
}
