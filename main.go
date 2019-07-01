package main

import (
	"auth-proxy/model/db"
	"auth-proxy/router"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"videos/model/redis"
)

func main() {

	// считать параметры командной строки
	servePort, env := readCommandLineParams()

	// читаем конфиги Postgres, и роутера.
	db.ReadConfig("./configs/db.yaml", env)
	router.ReadConfig("./configs/routes.yaml", env)

	// Инициализируем Redis
	redis.Init()

	// Ждем готовности базы данных
	db.WaitForDbOrExit(10)

	// порождаем базу данных если ее нет
	db.CreateDatabaseIfNotExists()

	// если servePort > 0, печатаем приветствие и запускаем сервер
	if servePort > 0 {
		greetings, _ := ioutil.ReadFile("./templates/greetings.txt")
		fmt.Printf(string(greetings), servePort)
		router.Serve(":" + strconv.Itoa(servePort))
	}
}

// Вспомогательные функции =========================================

// readCommandLineParams читает параметры командной строки
func readCommandLineParams() (serverPort int, env string) {
	flag.IntVar(&serverPort, "serve", 0, "Запустить приложение на порту с номером > 0 ")
	flag.StringVar(&env, "env", "prod", "Окружение. Возможные значения: dev - разработка, docker - в докере для фронтэнд разработчиков. prod - по умолчанию для продакшн.")
	flag.Parse()
	fmt.Println("\nПример запуска: ./auth-proxy -serve 4000 -env=dev")
	flag.Usage()
	return
}
