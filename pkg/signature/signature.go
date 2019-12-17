package signature

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"gopkg.in/spacemonkeygo/httpsig.v0"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type paramsType struct {
	PrivateKeyFile string `yaml:"private_key_file"`
	PublicKeyFile  string `yaml:"public_key_file"`
	privateKey     interface{}
	publicKey      interface{}
	signer         *httpsig.Signer
}

var params paramsType

// Sign добавляет цифровую подпись к запросу как определено
// в спецификации RFC    <https://tools.ietf.org/html/draft-cavage-http-signatures-06>.
// Цифровая подпись это HTTP заголовок вида:
// Authorization: Signature keyId="auth-proxy",algorithm="rsa-sha256",headers="(request-target) host date",signature="ZKNCbJ67zB..."
func Sign(req *http.Request) {
	if params.signer == nil {
		log.Panicln("No signer")
		return
	}
	//add Date header
	t := time.Now()
	req.Header.Add("Date", t.Format(time.RFC1123))

	// Signing
	err := params.signer.Sign(req)
	if err != nil {
		log.Panicln("SIGNING ERROR:", err)
		return
	}
	fmt.Println("-------------------------------------------")
	fmt.Println("AUTHORIZATION:", req.Header.Get("Authorization"))
	fmt.Println("-------------------------------------------")
}

// Verify верифицирует цифровую подпись к запросу как определено
// в спецификации RFC    <https://tools.ietf.org/html/draft-cavage-http-signatures-06>.
// Цифровая подпись это HTTP заголовок вида:
// Authorization: Signature keyId="auth-proxy",algorithm="rsa-sha256",headers="(request-target) host date",signature="ZKNCbJ67zB..."
func Verify(req *http.Request) {

}

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

	// load RSA keys
	params.privateKey, err = loadPrivateKey(params.PrivateKeyFile)
	if err != nil {
		fmt.Println("privateKey error:", err)
	}

	// params.signer = httpsig.NewSigner("auth-proxy", params.privateKey, httpsig.RSASHA256, nil)
	params.signer = httpsig.NewSigner("auth-proxy", params.privateKey, httpsig.RSASHA256, []string{"(request-target)", "host", "date"})
}

func loadPrivateKey(path string) (interface{}, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parsePrivateKey(bytes)
}

func parsePrivateKey(pemBytes []byte) (interface{}, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return rawkey, nil
}
