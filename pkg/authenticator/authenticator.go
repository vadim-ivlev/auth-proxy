package authenticator

import (
	"auth-proxy/pkg/app"
	"auth-proxy/pkg/db"
	"auth-proxy/pkg/mail"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var authenticatorURL = "https://www.authenticatorapi.com"
var AppName = "auth-proxy"
var secret = "supersecret"
var barcodeRe = regexp.MustCompile(`src='(.*?)'`)
var manualcodeRe = regexp.MustCompile(`secret%3D(.*?)%`)

// IsPinGood верен ли PIN введенный пользователем
func IsPinGood(username, pin string) error {
	if len(pin) < 2 {
		return errors.New("PIN  должен быть длиннее")
	}

	// validate pin
	secretCode, err := getSecretCode(username)
	if err != nil {
		return err
	}
	text, err := getResponseText(fmt.Sprintf(`%v/Validate.aspx?Pin=%v&SecretCode=%v`, authenticatorURL, pin, secretCode))
	log.Printf(`%v/Validate.aspx?  Pin=%v  SecretCode=%v rand=%v --> text=%v`, authenticatorURL, pin, secretCode, text)
	if err != nil {
		return err
	}
	if text == "True" {
		return nil
	}
	return errors.New("PIN неверен")
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
	_, err = db.QueryExec(`UPDATE "user" SET pinset = TRUE WHERE username = $1 OR email = $1 `, username)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"result": true, "error": nil})
}

// ResetAuthenticator если пин правильный устанавливает поле pinset=TRUE для пользователя в базе данных
func ResetAuthenticator(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(200, gin.H{"result": false, "error": "username is required"})
		return
	}
	// - утанавливаем поля pinset и pinhash в базе
	pinhash := uuid.New().String()
	log.Println("uuid=", pinhash)

	_, err := db.QueryExec(`UPDATE "user" SET ( pinset, pinhash ) = ( FALSE, $1 ) WHERE username = $2 OR email = $2 ;`, pinhash, username)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	// - генерируем ссылку на страничку

	link := fmt.Sprintf(`%v/set-authenticator.html#url=%v&username=%v&pinhash=%v`, app.Params.AuthAdminUrl, app.Params.AuthProxyUrl, username, pinhash)
	user, err := db.QueryRowMap(`SELECT * FROM "user" WHERE username=$1 OR email=$1`, username)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	email, _ := user["email"].(string)
	err = mail.SendAuthenticatorEmail(email, link)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"result": "Письмо с инструкциями выслано по электронной почте ", "error": nil})
}

// AuthenticatorBarcode возвращает изображение Barcode
// Google Authenticator на телефоне клиента.
func AuthenticatorBarcode(c *gin.Context) {
	AuthenticatorCode(c, "barcode")
}

// AuthenticatorManualCode возвращает код для ручной установки
// Google Authenticator на телефоне клиента.
func AuthenticatorManualCode(c *gin.Context) {
	AuthenticatorCode(c, "manualcode")
}

// AuthenticatorCode возвращает изображение Barcode
// или код для ручной установки Google Authenticator на телефоне клиента
// в зависимости от параметра:
// codetype string = {barcode|manualcode}
func AuthenticatorCode(c *gin.Context, codetype string) {
	username := c.Param("username")
	pinhash := c.Param("pinhash")
	// для обновления картинок в браузерах
	c.Header("Cache-control", "no-cache")
	// log.Printf(`AuthenticatorBarcode username=%v`, username)

	// проверяем есть ли пользователь в базе данных и установлен ли уже аутентификатор
	_, dbPinSet, dbPinHash, err := GetUserPinFields(username)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	if dbPinSet {
		c.JSON(200, gin.H{"error": "Аутентификатор уже установлен"})
		return
	}
	if pinhash != dbPinHash {
		c.JSON(200, gin.H{"error": "pinhash не совпадает"})
		return
	}

	//  get barcode image url
	secretCode, err := getSecretCode(username)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}

	pairUrl := fmt.Sprintf(`%v/pair.aspx?AppName=%v&AppInfo=%v&SecretCode=%v`, authenticatorURL, AppName, username, secretCode)
	text, err := getResponseText(pairUrl)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}

	if codetype == "manualcode" {
		match := manualcodeRe.FindStringSubmatch(text)
		if len(match) < 2 {
			c.JSON(200, gin.H{"result": nil, "error": "manualcode не найден в ответе " + authenticatorURL})
			return
		}
		c.JSON(200, gin.H{"result": match[1], "error": nil})
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

func GetUserPinFields(username string) (pinRequired, pinSet bool, pinHash string, err error) {
	user, err := db.QueryRowMap(`SELECT pinrequired, pinset, pinhash FROM "user" WHERE username=$1 OR email=$1 `, username)
	if err != nil {
		return
	}
	pinRequired, _ = user["pinrequired"].(bool)
	pinSet, _ = user["pinset"].(bool)
	pinHash, _ = user["pinhash"].(string)
	return
}

func getSecretCode(username string) (secredCode string, err error) {
	_, _, pinHash, err := GetUserPinFields(username)
	if err != nil {
		log.Printf(`ERROR: authenticator.getSecretCode("%v"): %v`, username, err.Error())
		return
	}
	secredCode = username + secret + pinHash
	return
}

func getResponseText(url string) (string, error) {
	body, err := getResponseBody(url)
	return string(body), err
}

func getResponseBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	// client := &http.Client{}
	// req, _ := http.NewRequest("GET", url, nil)
	// req.Header.Set("cache", "no-store")
	// resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
