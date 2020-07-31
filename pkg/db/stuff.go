package db

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	yaml "gopkg.in/yaml.v2"
)

// UsePool использовать ли пул соединений для подключения к БД.
var UsePool bool = true

// DBPool пул соединений
var DBPool *sqlx.DB = nil

// SQLite Использовать ли SQLite
var SQLite bool = false

// параметры подсоединения к Postgres
type postgresConnectParams struct {
	Host       string
	Port       string
	User       string
	Password   string
	Dbname     string `yaml:"database"`
	Sslmode    string
	SearchPath string `yaml:"search_path"`
	connectStr string
}

// параметры подсоединения к SQLite
type sqliteConnectParams struct {
	Sqlitefile string
}

var params postgresConnectParams
var sqliteParams sqliteConnectParams

// var connectURL string

// ReadConfig reads YAML with Postgres params
func ReadConfig(fileName string, env string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}

	envParams := make(map[string]postgresConnectParams)
	err = yaml.Unmarshal(yamlFile, &envParams)
	printIf("ReadConfig()", err)
	params = envParams[env]
	params.connectStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s", params.Host, params.Port, params.User, params.Password, params.Dbname, params.Sslmode, params.SearchPath)
}

// ReadSQLiteConfig reads YAML with SQLite params
func ReadSQLiteConfig(fileName string, env string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}

	envParams := make(map[string]sqliteConnectParams)
	err = yaml.Unmarshal(yamlFile, &envParams)
	printIf("ReadSQLiteConfig()", err)
	sqliteParams = envParams[env]
}

// PrintConfig prints DB connection parameters.
func PrintConfig() {
	s, _ := yaml.Marshal(params)
	fmt.Printf("\nDB connection parameters:\n%s\n", s)
	fmt.Printf("Postgres connection string: %s\n", params.connectStr)
	fmt.Printf("SQLite file: %s\n", sqliteParams.Sqlitefile)
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
