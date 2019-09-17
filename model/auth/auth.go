package auth

import (
	"auth-proxy/model/db"
	"auth-proxy/model/session"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/patrickmn/go-cache"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// var Store = gsessions.NewCookieStore([]byte("secret"))

var Cache = cache.New(2*time.Minute, 4*time.Minute)

// CheckUserPassword проверяет пароль пользователя. Возвращает true если проверка прошла успешно.
func CheckUserPassword(username, password string) bool {
	_, err := db.QueryRowMap(`
		SELECT * FROM "user" WHERE username=$1 AND "password"=$2 AND disabled = 0
		UNION
		SELECT * FROM "user" WHERE email=$1    AND "password"=$2 AND disabled = 0
		`, username, GetHash(password))

	if err != nil {
		return false
	}
	return true
}

// UlidNum - возвращает случайную строку числа в диапазоне [min,max)
func UlidNum(min, max int) string {
	t := time.Now().UnixNano()
	rand.Seed(t)
	return strconv.Itoa(rand.Intn(max-min) + min)
}

// GenerateNewPassword Генерирует новый пароль для пользователя,
// Сохраняет его в базе данных.
func GenerateNewPassword(usernameOrEmail string) (string, string, string, error) {
	rec, err := db.QueryRowMap(`
		SELECT * FROM "user" WHERE username=$1 AND disabled = 0
		UNION
		SELECT * FROM "user" WHERE email=$1    AND disabled = 0
		`, usernameOrEmail)

	if err != nil {
		return "", "", "", err
	}
	foundUsername := rec["username"].(string)
	foundEmail := rec["email"].(string)
	if foundEmail == "" {
		return "", "", "", errors.New("No email address for user " + foundUsername)
	}
	// generate new password
	newPassword := UlidNum(100000, 999999)
	// newPassword := "123456"
	newHash := GetHash(newPassword)
	sqlText := `UPDATE "user" SET password = $1 WHERE username = $2;`
	_, err = db.QueryExec(sqlText, newHash, foundUsername)
	if err != nil {
		return "", "", "", err
	}
	return foundUsername, foundEmail, newPassword, nil
}

// GetAppURLs Возвращает url-ы приложений.
func GetAppURLs() (map[string][]string, error) {
	records, err := db.QuerySliceMap(`SELECT appname,url,rebase FROM app WHERE url IS NOT NULL;`)
	if err != nil {
		return nil, err
	}
	m := make(map[string][]string)
	for _, rec := range records {
		app, _ := rec["appname"].(string)
		url, _ := rec["url"].(string)
		rebase, _ := rec["rebase"].(string)
		if url == "" || app == "" {
			continue
		}
		m[app] = []string{url, rebase}
	}
	return m, nil
}

// GetUserRoles возвращает строку с сериализованным масссивом ролей пользователя в заданном приложении.
func GetUserRoles(user, app string) string {
	cacheKey := user + "-" + app + "-roles"
	cachedValue, found := Cache.Get(cacheKey)
	if found {
		fmt.Println("cached ", cacheKey, "=", cachedValue)
		return cachedValue.(string)
	}

	records, err := db.QuerySliceMap(`SELECT rolename FROM app_user_role WHERE  appname  = $1 AND username = $2 `, app, user)
	if err != nil {
		return ""
	}

	bytes, _ := json.Marshal(records)
	roles := string(bytes)

	Cache.Set(cacheKey, roles, cache.DefaultExpiration)
	return roles

}

// GetUserInfo возвращает сериализованную информацию о пользователе
func GetUserInfo(user string) string {
	cacheKey := user + "-info"
	cachedValue, found := Cache.Get(cacheKey)
	if found {
		fmt.Println("cached ", cacheKey, "=", cachedValue)
		return cachedValue.(string)
	}

	record, err := db.QueryRowMap(`SELECT username, email, fullname, description 
		FROM "user" 
		WHERE username = $1;`, user)
	if err != nil {
		return ""
	}
	jsonBytes, _ := json.Marshal(record)
	jsonString := string(jsonBytes)

	Cache.Set(cacheKey, jsonString, cache.DefaultExpiration)
	return jsonString
}

// AppUserRoleExist проверят наличие связки appname-username-rolename в таблице app_user_role.
func AppUserRoleExist(appname, username, rolename string) bool {
	_, err := db.QueryRowMap(`SELECT * FROM  app_user_role  
		WHERE appname = $1 AND username = $2 AND rolename = $3 ;`,
		appname, username, rolename)
	if err == nil {
		return true
	}
	return false
}

// Logout разлогинить текущего пользователя
func Logout(c *gin.Context) {
	session.DeleteSession(c)
}

// Login залогинить пользователя
func Login(c *gin.Context, username, password string) error {
	if CheckUserPassword(username, password) {
		return session.SetVariable(c, "user", username)
	} else {
		return errors.New("Authentication failed")
	}
}

// IsUserEnabled  Если false, пользователь отключен
func IsUserEnabled(user string) bool {
	cacheKey := user + "-enabled"
	cachedValue, found := Cache.Get(cacheKey)
	if found {
		fmt.Println("cached ", cacheKey, "=", cachedValue)
		return cachedValue.(bool)
	}

	_, err := db.QueryRowMap(`SELECT * FROM "user" WHERE username = $1 AND disabled = 0 ;`, user)
	if err == nil {
		Cache.Set(cacheKey, true, cache.DefaultExpiration)
		return true
	}
	Cache.Set(cacheKey, false, cache.DefaultExpiration)
	return false
}

// IsAppPublic  true если приложение доступно для пользователей без роли
func IsAppPublic(appName string) bool {
	cacheKey := "is-" + appName + "-public"
	cachedValue, found := Cache.Get(cacheKey)
	if found {
		fmt.Println("cached ", cacheKey, "=", cachedValue)
		return cachedValue.(bool)
	}

	_, err := db.QueryRowMap(`SELECT * FROM "app" WHERE appname = $1 AND public = 'Y' ;`, appName)
	if err == nil {
		Cache.Set(cacheKey, true, cache.DefaultExpiration)
		return true
	}
	Cache.Set(cacheKey, false, cache.DefaultExpiration)
	return false
}

func GetHash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}
