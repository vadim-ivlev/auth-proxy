package db

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	yaml "gopkg.in/yaml.v2"
)

// Вспомогательные функции /////////////////////////////////////////////////////

type connectionParams struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

var params connectionParams
var connectStr string
var connectURL string

// ReadConfig reads YAML file
func ReadConfig(fileName string, env string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}

	envParams := make(map[string]connectionParams)
	err = yaml.Unmarshal(yamlFile, &envParams)
	printIf("ReadConfig()", err)
	params = envParams[env]

	connectStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=public", params.Host, params.Port, params.User, params.Password, params.Dbname, params.Sslmode)
	connectURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", params.User, params.Password, params.Host, params.Port, params.Dbname, params.Sslmode)

}

// PrintConfig prints DB connection parameters.
func PrintConfig() {
	s, _ := yaml.Marshal(params)
	fmt.Printf("\nDB connection parameters:\n%s\n", s)
	fmt.Printf("DB connection string: %s\n", connectStr)
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

// CreateDatabaseIfNotExists порождает объекты базы данных и наполняет базу тестовыми данными
func CreateDatabaseIfNotExists() {
	fmt.Println("Миграция ...")
	m, err := migrate.New("file://migrations/", connectURL)
	panicIf(err)
	printIf("CreateDatabaseIfNotExists()", m.Up())
}
