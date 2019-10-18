package auth

import (
	"auth-proxy/pkg/db"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/patrickmn/go-cache"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var Cache = cache.New(5*time.Minute, 10*time.Minute)

const (
	OK             = 0
	NO_USER        = 1
	USER_DISABLED  = 2
	WRONG_PASSWORD = 3
)

// CheckUserPassword2 проверяет пароль пользователя.
// Возвращает:
// OK 			 - проверка прошла успешно.
// NO_USER 		 - пользователя нет.
// USER_DISABLED - пользователь заблокирован.
// WRONG_PASSWORD - пароль не подходит.
func CheckUserPassword2(username, password string) (int, string) {
	rec, err := db.QueryRowMap(`SELECT * FROM "user" WHERE username=$1 OR email=$1`, username)
	if err != nil {
		return NO_USER, ""
	}
	dbUsername := rec["username"].(string)
	disabled := rec["disabled"].(int64)
	if disabled > 0 {
		return USER_DISABLED, dbUsername
	}
	hashedPassword := rec["password"].(string)
	if hashedPassword != GetHash(password) {
		return WRONG_PASSWORD, dbUsername
	}
	return OK, dbUsername
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
