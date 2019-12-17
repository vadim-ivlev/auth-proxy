package signature

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/spacemonkeygo/httpsig.v0"
	"gopkg.in/yaml.v2"
)

type paramsType struct {
	PrivateKeyFile string `yaml:"private_key_file"`
	PublicKeyFile  string `yaml:"public_key_file"`
}

var params paramsType
var keyID = "auth-proxy"
var headers = []string{"(request-target)", "host", "date"}

var privateKey interface{}
var signer *httpsig.Signer

// var publicKey interface{}
// var keystore *httpsig.MemoryKeyStore
var verifier *httpsig.Verifier

// Sign добавляет цифровую подпись к запросу как определено
// в спецификации RFC    <https://tools.ietf.org/html/draft-cavage-http-signatures-06>.
// Цифровая подпись это HTTP заголовок вида:
// Authorization: Signature keyId="auth-proxy",algorithm="rsa-sha256",headers="(request-target) host date",signature="ZKNCbJ67zB..."
func Sign(req *http.Request) error {
	if signer == nil {
		log.Println("No signer")
		return errors.New("No signer")
	}
	//add Date header
	t := time.Now()
	req.Header.Add("Date", t.Format(time.RFC1123))

	// Signing
	err := signer.Sign(req)
	if err != nil {
		log.Println("SIGNING ERROR:", err)
		return err
	}
	fmt.Println("-------------------------------------------")
	fmt.Println("AUTHORIZATION:", req.Header.Get("Authorization"))
	fmt.Println("-------------------------------------------")
	return nil
}

// Verify верифицирует цифровую подпись к запросу как определено
// в спецификации RFC    <https://tools.ietf.org/html/draft-cavage-http-signatures-06>.
// Цифровая подпись это HTTP заголовок вида:
// Authorization: Signature keyId="auth-proxy",algorithm="rsa-sha256",headers="(request-target) host date",signature="ZKNCbJ67zB..."
func Verify(req *http.Request) error {
	if verifier == nil {
		log.Println("No verifier")
		return errors.New("No verifier")
	}
	err := verifier.Verify(req)
	if err == nil {
		fmt.Println("Verified")
	} else {
		fmt.Println("Verification error:", err)
	}
	return err
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

	loadPrivateKeyFromFile()
	loadPublicKeyFromFile()
}

// loadPrivateKeyFromFile load private RSA key
func loadPrivateKeyFromFile() {
	privateKey, err := loadPrivateKey(params.PrivateKeyFile)
	if err != nil {
		fmt.Println("privateKey error:", err)
	} else {
		signer = httpsig.NewSigner(keyID, privateKey, httpsig.RSASHA256, headers)
	}
}

// loadPublicKeyFromFile load public RSA key from file
func loadPublicKeyFromFile() {
	publicKey, err := loadPublicKey(params.PublicKeyFile)
	if err != nil {
		fmt.Println("publicKey error:", err)
	} else {
		// fmt.Printf("Publ key=%#v\n", publicKey)
		keystore := httpsig.NewMemoryKeyStore()
		keystore.SetKey(keyID, publicKey)
		verifier = httpsig.NewVerifier(keystore)
	}
}

// helpers --------------------------------------------------------

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

func loadPublicKey(path string) (interface{}, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parsePublicKey(bytes)
}

func parsePublicKey(pemBytes []byte) (interface{}, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "PUBLIC KEY":
		rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}

	return rawkey, nil
}
