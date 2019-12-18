package signature

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/spacemonkeygo/httpsig.v0"
)

var verifier *httpsig.Verifier
var PublicKeyText string

// Verify верифицирует цифровую подпись к запросу как определено
// в спецификации RFC    <https://tools.ietf.org/html/draft-cavage-http-signatures-06>.
// Цифровая подпись это HTTP заголовок вида:
// Authorization: Signature keyId="auth-proxy",algorithm="rsa-sha256",headers="(request-target) host date",signature="ZKNCbJ67zB..."
// Возвращает ошибку.
func Verify(req *http.Request) error {
	if verifier == nil {
		return errors.New("No verifier")
	}
	err := verifier.Verify(req)
	return err
}

// loadPublicKeyFromFile load public RSA key from file
func loadPublicKeyFromFile() {
	publicKey, err := loadPublicKey(params.PublicKeyFile)
	if err != nil {
		fmt.Println("publicKey error:", err)
	} else {
		keystore := httpsig.NewMemoryKeyStore()
		keystore.SetKey(params.KeyID, publicKey)
		verifier = httpsig.NewVerifier(keystore)
	}
}

// helpers --------------------------------------------------------

func loadPublicKey(path string) (interface{}, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	PublicKeyText = string(bytes)
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
