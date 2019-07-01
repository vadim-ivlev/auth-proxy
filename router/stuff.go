package router

import (
	"io/ioutil"
	"log"

	//blank import
	_ "github.com/golang-migrate/migrate/database/postgres"
	//blank import
	_ "github.com/golang-migrate/migrate/source/file"
	yaml "gopkg.in/yaml.v2"
)

// Вспомогательные функции /////////////////////////////////////////////////////

type Applications map[string]string

var apps Applications

// ReadConfig reads YAML file
func ReadConfig(fileName string, env string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}

	envParams := make(map[string]Applications)
	err = yaml.Unmarshal(yamlFile, &envParams)
	if err != nil {
		log.Println("ERROR: router.ReadConfig():", err.Error())
	}
	apps = envParams[env]
}
