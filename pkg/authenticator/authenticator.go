package authenticator

import (
	"auth-proxy/pkg/app"
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/db"
	"auth-proxy/pkg/mail"
	"crypto/tls"
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

// var AppName = "auth-proxy"
var secret = "supersecret"
var barcodeRe = regexp.MustCompile(`src='(.*?)'`)
var manualcodeRe = regexp.MustCompile(`secret%3D(.*?)%`)

// IsPinGood верен ли PIN введенный пользователем
func IsPinGood(username, pin string, tempCode bool) error {
	if len(pin) < 2 {
		return errors.New("PIN  должен быть длиннее")
	}

	// validate pin
	secretCode, err := getSecretCode(username, tempCode)
	if err != nil {
		return err
	}
	text, err := getResponseText(fmt.Sprintf(`%v/Validate.aspx?Pin=%v&SecretCode=%v`, authenticatorURL, pin, secretCode))
	log.Printf(`%v/Validate.aspx?  Pin=%v  SecretCode=%v --> text=%v`, authenticatorURL, pin, secretCode, text)
	if err != nil {
		return err
	}
	if text == "True" {
		return nil
	}
	return errors.New("PIN неверен")
}

// SetAuthenticator если пин правильный  поле pinhash_temp в поле pashash для пользователя в базе данных
func SetAuthenticator(c *gin.Context) {
	username := c.Query("username")
	pin := c.Query("pin")
	// log.Printf(`SetAuthenticator username=%v pin=%v`, username, pin)
	err := IsPinGood(username, pin, true)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	// update field pinhash in the database
	_, pinHashTemp, _, err := GetUserPinFields(username)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	_, err = db.QueryExec(`UPDATE "user" SET (pinhash, pinhash_temp) = ($1, NULL)  WHERE username = $2 OR email = $2 `, pinHashTemp, username)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"result": true, "error": nil})
}

// SetPassword если пин правильный устанавливает поле pashash='' и пароль для пользователя в базе данных
func SetPassword(c *gin.Context) {
	username := c.Query("username")
	hash := c.Query("hash")
	password := c.Query("password")
	log.Printf(`SetPassword username=%v password=%v hash=%v`, username, password, hash)

	// проверить длину пароля
	if len(password) < 6 {
		c.JSON(200, gin.H{"result": false, "error": "Пароль должен быть не менее 6-ти символов"})
		return
	}
	// зашифровать пароль
	password = auth.GetHash(password)

	// обновить поля в базе
	res, err := db.QueryExec(`UPDATE "user" SET ( password, pashash ) = ( $1, NULL ) WHERE pashash = $2`, password, hash)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": "SetPassword: " + err.Error()})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": "SetPassword: " + err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(200, gin.H{"result": false, "error": "Пароль уже был установлен"})
		return
	}
	c.JSON(200, gin.H{"result": true, "error": nil})
}

// ResetPassword устанавливает поле pashash для пользователя в базе данных
// и посылает ему письмо по email с адресом страницы установки пароля.
func ResetPassword(c *gin.Context) {
	username := c.Query("username")
	adminurl := c.Query("adminurl")
	authurl := c.Query("authurl")

	if username == "" {
		c.JSON(200, gin.H{"result": false, "error": "username is required"})
		return
	}
	// - устанавливаем поле pashash в для пользователя базе
	hash := uuid.New().String()
	log.Println("pashash=", hash)

	_, err := db.QueryExec(`UPDATE "user" SET  pashash  = $1  WHERE username = $2 OR email = $2 ;`, hash, username)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	// - генерируем ссылку на страничку
	link := fmt.Sprintf(`%vset-password.html#username=%v&hash=%v&authurl=%v`, adminurl, username, hash, authurl)
	// - находим email пользователя
	user, err := db.QueryRowMap(`SELECT * FROM "user" WHERE username=$1 OR email=$1`, username)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	email, _ := user["email"].(string)
	// посылаем письмо пользователю
	err = mail.SendResetPasswordEmail(email, link)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"result": "Письмо с инструкциями выслано по электронной почте ", "error": nil})
}

// ResetAuthenticator устанавливает поле pinhash для пользователя в базе данных
// и посылает ему письмо по email с адресом страницы установки пина.
func ResetAuthenticator(c *gin.Context) {
	username := c.Query("username")
	adminurl := c.Query("adminurl")
	authurl := c.Query("authurl")

	if username == "" {
		c.JSON(200, gin.H{"result": false, "error": "username is required"})
		return
	}
	// - устанавливаем  pinhash_temp в базе
	hash := uuid.New().String()
	log.Println("hash=", hash)

	_, err := db.QueryExec(`UPDATE "user" SET pinhash_temp = $1  WHERE username = $2 OR email = $2 ;`, hash, username)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	// - генерируем ссылку на страничку
	link := fmt.Sprintf(`%vset-authenticator.html#username=%v&hash=%v&authurl=%v`, adminurl, username, hash, authurl)
	// - находим email пользователя
	user, err := db.QueryRowMap(`SELECT * FROM "user" WHERE username=$1 OR email=$1`, username)
	if err != nil {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}
	email, _ := user["email"].(string)
	// посылаем письмо пользователю
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
	username := c.Query("username")
	pinhash := c.Query("hash")
	// для обновления картинок в браузерах
	c.Header("Cache-control", "no-cache")
	// log.Printf(`AuthenticatorBarcode username=%v`, username)

	// проверяем есть ли пользователь в базе данных и установлен ли уже аутентификатор
	_, dbPinHashTemp, _, err := GetUserPinFields(username)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	if dbPinHashTemp == "" {
		c.JSON(200, gin.H{"error": "Аутентификатор уже установлен"})
		return
	}
	if pinhash != dbPinHashTemp {
		c.JSON(200, gin.H{"error": "pinhash не совпадает"})
		return
	}

	//  get barcode image url
	secretCode, err := getSecretCode(username, true)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}

	pairUrl := fmt.Sprintf(`%v/pair.aspx?AppName=%v&AppInfo=%v&SecretCode=%v`, authenticatorURL, app.Params.AppName, username, secretCode)
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

func GetUserPinFields(username string) (pinRequired bool, pinHashTemp, pinHash string, err error) {
	user, err := db.QueryRowMap(`SELECT pinrequired, pinhash_temp, pinhash FROM "user" WHERE username=$1 OR email=$1 `, username)
	if err != nil {
		return
	}
	pinRequired, _ = user["pinrequired"].(bool)
	pinHashTemp, _ = user["pinhash_temp"].(string)
	pinHash, _ = user["pinhash"].(string)
	return
}

func getSecretCode(username string, tempCode bool) (secretCode string, err error) {
	_, pinHashTemp, pinHash, err := GetUserPinFields(username)
	if err != nil {
		log.Printf(`ERROR: authenticator.getSecretCode("%v"): %v`, username, err.Error())
		return
	}
	secretCode = username + secret + pinHash
	if tempCode {
		secretCode = username + secret + pinHashTemp
	}
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
		// return nil, err
		// Если произошла ошибка пробуем сделать запрос без проверки сертификата
		// https://stackoverflow.com/questions/12122159/how-to-do-a-https-request-with-bad-certificate
		log.Println("ERROR: Произошла ошибка запроса к ", url)
		log.Println(err.Error())
		log.Println("Пробуем запрос без проверки сертификата ------------------ ")
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err = client.Get(url)
		if err != nil {
			log.Println("ERROR: Запрос без проверки сертификата тоже вернул ошибку")
			log.Println(err.Error())
			return nil, err
		}
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
