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
	Refreshtime int    `json:"PG_REFRESH_TIME" env:"PG_REFRESH_TIME" envDefault:"5"`
	Host        string `json:"PG_HOST" env:"PG_HOST" envDefault:"localhost"`
	Port        string `json:"PG_PORT" env:"PG_PORT" envDefault:"5432"`
	User        string `json:"PG_USER" env:"PG_USER" envDefault:"pgadmin"`
	Password    string `json:"PG_PASSWORD" env:"PG_PASSWORD" envDefault:"159753"`
	Database    string `json:"PG_DATABASE" env:"PG_DATABASE" envDefault:"rgru"`
	Sslmode     string `json:"PG_SSLMODE" env:"PG_SSLMODE" envDefault:"disable"`
	SearchPath  string `json:"PG_SEARCH_PATH" env:"PG_SEARCH_PATH" envDefault:"auth,extensions"`
	ConnectStr  string `json:"PG_CONNECT_STR" env:"PG_CONNECT_STR"`
}

var Params postgresConnectParams

// ReadEnvConfig reads YAML with Postgres params
func ReadEnvConfig(fileName string) {
	if fileName != "" {
		if err := godotenv.Load(fileName); err != nil {
			log.Println("ОШИБКА чтения env файла:", err.Error())
		}
	} else {
		fmt.Println("Параметры Postgres берутся из операционной системы.")
	}

	if err := env.Parse(&Params); err != nil {
		fmt.Printf("%+v\n", err)
	}
	Params.SearchPath = strings.Replace(Params.SearchPath, " ", "", -1)
	Params.ConnectStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s", Params.Host, Params.Port, Params.User, Params.Password, Params.Database, Params.Sslmode, Params.SearchPath)
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
