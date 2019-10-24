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
	fmt.Println("████████████████████████ revision 6 ████████████████████████")
	// считать параметры командной строки
	servePort, env /*, sqlite, tls, selfreg, secure, captcha, max_attempts, reset_time, admin_url*/ := readCommandLineParams()
	// db.SQLite = sqlite
	// server.SelfRegistrationAllowed = selfreg
	// server.SecureCookie = secure
	// server.IsCaptchaRequired = captcha
	// server.AdminUrl = admin_url
	// counter.MAX_ATTEMPTS = max_attempts
	// counter.RESET_TIME = time.Duration(reset_time)

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
	server.IsCaptchaRequired = app.Params.Captcha
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
func readCommandLineParams() (serverPort int, env string /*, sqlite bool, tls bool, selfreg bool, secure bool, captcha bool, max_attempts int64, reset_time int64, admin_url string*/) {
	flag.IntVar(&serverPort, "serve", 0, "Запустить приложение на порту с номером > 0 ")
	flag.StringVar(&env, "env", "prod", "Окружение. Возможные значения: dev - разработка, front - в докере для фронтэнд разработчиков. prod - по умолчанию для продакшн.")
	// flag.StringVar(&admin_url, "admin_url", "https://auth-admin.now.sh", "URL сайта где размещена админка сервиса. Туда перенаправляется браузер когда пользователь вводит путь /admin")
	// flag.BoolVar(&sqlite, "sqlite", false, "Использовать SQLite")
	// flag.BoolVar(&tls, "tls", false, "Использовать https вместо http")
	// flag.BoolVar(&selfreg, "selfreg", false, "Пользователи могут регистрироваться самостоятельно")
	// flag.BoolVar(&secure, "secure", false, "Установить флаг secure на куки браузера. Работает для https протокола.")
	// flag.BoolVar(&captcha, "captcha", false, "Нужно ли вводить капчу при входе в систему")

	// flag.Int64Var(&max_attempts, "max_attempts", 5, "Максимально допустимое число ошибок ввода пароля")
	// flag.Int64Var(&reset_time, "reset_time", 60, "Время сброса счетчика ошибок пароля в минутах")

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
Environment: %v
Database:%v 
TLS:%v

%v://localhost:%v/testapp
%v://localhost:%v/admin
		
CTRL-C to interrupt.
`
	fmt.Printf(msg, env, database, tls, protocol, serverPort, protocol, serverPort)
}
