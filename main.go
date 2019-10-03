package main

import (
	"auth-proxy/pkg/db"
	"auth-proxy/pkg/mail"
	"auth-proxy/server"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("████████████████████████ revision 5 ████████████████████████")
	// считать параметры командной строки
	servePort, env, sqlite, tls, selfreg, secure := readCommandLineParams()
	db.SQLite = sqlite
	server.SelfRegistrationAllowed = selfreg
	server.SecureCookie = secure

	// читаем конфиг Postgres.
	db.ReadConfig("./configs/db.yaml", env)
	// читаем конфиг SQLite.
	db.ReadSQLiteConfig("./configs/sqlite.yaml", env)
	// читаем конфиг mail.
	mail.ReadConfig("./configs/mail.yaml", env)
	// читаем шаблоны писем
	mail.ReadMailTemplate("./configs/mail-templates.yaml")

	// Ждем готовности базы данных
	db.PrintConfig()
	db.WaitForDbOrExit(10)

	// порождаем базу данных если ее нет
	db.CreateDatabaseIfNotExists()

	// если servePort > 0, печатаем приветствие и запускаем сервер
	if servePort > 0 {
		printGreetings(servePort, env, sqlite, tls)
		server.Serve(":"+strconv.Itoa(servePort), tls)
	}

	db.DBPool.Close()
	fmt.Println("Bye")
}

// Вспомогательные функции =========================================

// readCommandLineParams читает параметры командной строки
func readCommandLineParams() (serverPort int, env string, sqlite bool, tls bool, selfreg bool, secure bool) {
	flag.IntVar(&serverPort, "serve", 0, "Запустить приложение на порту с номером > 0 ")
	flag.StringVar(&env, "env", "prod", "Окружение. Возможные значения: dev - разработка, front - в докере для фронтэнд разработчиков. prod - по умолчанию для продакшн.")
	flag.BoolVar(&sqlite, "sqlite", false, "Использовать SQLite")
	flag.BoolVar(&tls, "tls", false, "Использовать https вместо http")
	flag.BoolVar(&selfreg, "selfreg", false, "Пользователи могут регистрироваться самостоятельно")
	flag.BoolVar(&secure, "secure", false, "Установить флаг secure на куки браузера. Работает для https протокола.")
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
	// greetings, _ := ioutil.ReadFile("./templates/greetings.txt")
	// fmt.Printf(string(greetings))
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
Environment: %v
Database:%v 
TLS:%v

%v://localhost:%v/testapp
		
CTRL-C to interrupt.
`
	fmt.Printf(msg, env, database, tls, protocol, serverPort)
}
