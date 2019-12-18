package signature

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type paramsType struct {
	PrivateKeyFile string   `yaml:"private_key_file"`
	PublicKeyFile  string   `yaml:"public_key_file"`
	Headers        []string `yaml:"headers"`
	KeyID          string   `yaml:"keyid"`
}

var params paramsType

// **********************************************************************************

// ReadConfig reads YAML file
func ReadConfig(fileName string, env string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}

	envParams := make(map[string]paramsType)
	err = yaml.Unmarshal(yamlFile, &envParams)
	if err != nil {
		log.Println("Signature ReadConfig() error:", err)
	}
	params = envParams[env]

	loadPrivateKeyFromFile()
	loadPublicKeyFromFile()
}
