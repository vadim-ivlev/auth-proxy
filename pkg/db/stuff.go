package db

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	// yaml "gopkg.in/yaml.v2"
)

// DBPool пул соединений
var DBPool *sqlx.DB = nil

// параметры подсоединения к Postgres
type postgresConnectParams struct {
	Host       string `env:"PG_HOST"`
	Port       string `env:"PG_PORT"`
	User       string `env:"PG_USER"`
	Password   string `env:"PG_PASSWORD"`
	Dbname     string `env:"PG_DATABASE"`
	Sslmode    string `env:"PG_SSLMODE"`
	SearchPath string `env:"PG_SEARCH_PATH" envDefault:"auth,extensions"`
	connectStr string
}

var params postgresConnectParams

// ReadConfig reads YAML with Postgres params
func ReadEnvConfig(fileName string) {
	if err := godotenv.Load(fileName); err != nil {
		log.Println("ОШИБКА чтения env файла:", err.Error())
	}
	if err := env.Parse(&params); err != nil {
		fmt.Printf("%+v\n", err)
	}
	params.SearchPath = strings.Replace(params.SearchPath, " ", "", -1)
	params.connectStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s", params.Host, params.Port, params.User, params.Password, params.Dbname, params.Sslmode, params.SearchPath)
}

// PrintConfig prints DB connection parameters.
func PrintConfig() {
	fmt.Printf("Строка соединения Postgres: %s\n", params.connectStr)
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func printIf(msg string, err error) {
	if err != nil {
		log.Println(msg, err.Error())
	}
}

// WaitForDbConnection - Ожидает соединения с базой данных
func WaitForDbConnection() {
	for {
		fmt.Println("Пытаемся соединиться с базой.")
		if DbAvailable() {
			return
		}
		time.Sleep(5 * time.Second)
	}
}
