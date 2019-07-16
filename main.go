package main

import (
	"auth-proxy/model/db"
	"auth-proxy/router"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
)

func main() {

	// считать параметры командной строки
	servePort, env, sqlite := readCommandLineParams()
	db.SQLite = sqlite

	// читаем конфиги Postgres, и роутера.
	db.ReadConfig("./configs/db.yaml", env)
	router.ReadConfig("./configs/routes.yaml", env)

	// Ждем готовности базы данных
	db.WaitForDbOrExit(10)

	// порождаем базу данных если ее нет
	db.CreateDatabaseIfNotExists()

	// если servePort > 0, печатаем приветствие и запускаем сервер
	if servePort > 0 {
		greetings, _ := ioutil.ReadFile("./templates/greetings.txt")
		fmt.Printf(string(greetings), sqlite, servePort)
		router.Serve(":" + strconv.Itoa(servePort))
	}

	fmt.Println("Bye")
}

// Вспомогательные функции =========================================

// readCommandLineParams читает параметры командной строки
func readCommandLineParams() (serverPort int, env string, sqlite bool) {
	flag.IntVar(&serverPort, "serve", 0, "Запустить приложение на порту с номером > 0 ")
	flag.StringVar(&env, "env", "prod", "Окружение. Возможные значения: dev - разработка, docker - в докере для фронтэнд разработчиков. prod - по умолчанию для продакшн.")
	flag.BoolVar(&sqlite, "sqlite", false, "Использовать SQLite")
	flag.Parse()
	fmt.Println("\nПример запуска: ./auth-proxy -serve 4000 -env=dev")
	flag.Usage()
	return
}
