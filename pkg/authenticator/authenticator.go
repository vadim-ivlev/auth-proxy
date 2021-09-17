package authenticator

import (
	"auth-proxy/pkg/db"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var authenticatorURL = "https://www.authenticatorapi.com"
var AppName = "auth-proxy"
var secret = "supersecret"
var barcodeRe = regexp.MustCompile(`src='(.*?)'`)

// IsPinGood верен ли PIN введенный пользователем
func IsPinGood(username, pin string) error {
	if len(pin) < 2 {
		return errors.New("PIN  должен быть длиннее")
	}
	// validate pin
	text, err := getResponseText(fmt.Sprintf(`%v/Validate.aspx?Pin=%v&SecretCode=%v%v`, authenticatorURL, pin, username, secret))
	if err != nil {
		return err
	}
	if text == "True" {
		return nil
	}
	return errors.New("PIN is not correct")
}

// SetAuthenticator если пин правильный устанавливает поле pinset=TRUE для пользователя в базе данных
func SetAuthenticator(c *gin.Context) {
	username := c.Param("username")
	pin := c.Param("pin")
	// log.Printf(`SetAuthenticator username=%v pin=%v`, username, pin)
	err := IsPinGood(username, pin)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	// update field pinset in the database
	_, err = db.QueryExec(`UPDATE "user" SET pinset = TRUE WHERE username = $1;`, username)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"result": true, "error": nil})
}

// SetAuthenticator возвращает изображение Barcode для установки Google Authenticator
// на телефоне клиента
func AuthenticatorBarcode(c *gin.Context) {
	username := c.Param("username")
	// log.Printf(`AuthenticatorBarcode username=%v`, username)

	//  get barcode image url
	pairUrl := fmt.Sprintf(`%v/pair.aspx?AppName=%v&AppInfo=%v&SecretCode=%v%v`, authenticatorURL, AppName, username, username, secret)
	text, err := getResponseText(pairUrl)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	match := barcodeRe.FindStringSubmatch(text)
	if len(match) < 2 {
		c.JSON(200, gin.H{"error": "barcodeUrl не найден в ответе " + authenticatorURL})
		return
	}
	barcodeUrl := match[1]
	// log.Printf(`barcodeUrl=%v`, barcodeUrl)

	//  proxy to barcode url
	// https://hackernoon.com/writing-a-reverse-proxy-in-just-one-line-with-go-c1edfa78c84b
	// https://stackoverflow.com/questions/21270945/how-to-read-the-response-from-a-newsinglehostreverseproxy
	// https://www.sanarias.com/blog/1214PlayingwithimagesinHTTPresponseingolang

	// get the image
	body, err := getResponseBody(barcodeUrl)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	// return the image
	c.Data(200, "image/png", body)

}

func getResponseText(url string) (string, error) {
	body, err := getResponseBody(url)
	return string(body), err
}

func getResponseBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
