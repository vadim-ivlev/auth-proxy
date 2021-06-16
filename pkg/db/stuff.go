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
	Host       string `env:"PG_HOST" envDefault:"localhost"`
	Port       string `env:"PG_PORT" envDefault:"5432"`
	User       string `env:"PG_USER" envDefault:"pgadmin"`
	Password   string `env:"PG_PASSWORD" envDefault:"159753"`
	Database   string `env:"PG_DATABASE" envDefault:"rgru"`
	Sslmode    string `env:"PG_SSLMODE" envDefault:"disable"`
	SearchPath string `env:"PG_SEARCH_PATH" envDefault:"auth,extensions"`
	connectStr string
}

var params postgresConnectParams

// ReadConfig reads YAML with Postgres params
func ReadEnvConfig(fileName string) {
	if fileName != "" {
		if err := godotenv.Load(fileName); err != nil {
			log.Println("ОШИБКА чтения env файла:", err.Error())
		}
	} else {
		fmt.Println("Параметры Postgres берутся из операционной системы.")
	}

	if err := env.Parse(&params); err != nil {
		fmt.Printf("%+v\n", err)
	}
	params.SearchPath = strings.Replace(params.SearchPath, " ", "", -1)
	params.connectStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s", params.Host, params.Port, params.User, params.Password, params.Database, params.Sslmode, params.SearchPath)
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
