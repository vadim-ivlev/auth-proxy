package signature

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/spacemonkeygo/httpsig.v0"
)

// var privateKey interface{}
var signer *httpsig.Signer

// Sign добавляет цифровую подпись к запросу как определено
// в спецификации RFC    <https://tools.ietf.org/html/draft-cavage-http-signatures-06>.
// Цифровая подпись это HTTP заголовок вида:
// Authorization: Signature keyId="auth-proxy",algorithm="rsa-sha256",headers="(request-target) host date",signature="ZKNCbJ67zB..."
// Возвращает ошибку.
func Sign(req *http.Request) error {
	if signer == nil {
		return errors.New("No signer")
	}
	//add Date header
	t := time.Now()
	req.Header.Add("Date", t.Format(time.RFC1123))

	// Signing
	err := signer.Sign(req)
	if err != nil {
		return err
	}
	return nil
}

// loadPrivateKeyFromFile load private RSA key
func loadPrivateKeyFromFile() {
	privateKey, err := loadPrivateKey(params.PrivateKeyFile)
	if err != nil {
		fmt.Println("privateKey error:", err)
	} else {
		signer = httpsig.NewSigner(params.KeyID, privateKey, httpsig.RSASHA256, params.Headers)
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
