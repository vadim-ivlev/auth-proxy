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

// GetUserRoles возвращает массив ролей пользователя в заданном приложении.
func GetUserRoles(user, app string) []string {
	records, err := db.QuerySliceMap(`SELECT rolename FROM app_user_role WHERE  appname  = $1 AND username = $2 `, app, user)
	if err != nil {
		return []string{}
	}

	roles := []string{}
	for _, rec := range records {
		role, _ := rec["rolename"].(string)
		roles = append(roles, role)
	}

	return roles
}

// GetUserRolesString возвращает строку с сериализованным масссивом ролей пользователя в заданном приложении.
func GetUserRolesString(user, app string) string {
	cacheKey := user + "-" + app + "-roles"
	cachedValue, found := Cache.Get(cacheKey)
	if found {
		fmt.Println("cached ", cacheKey, "=", cachedValue)
		return cachedValue.(string)
	}
	roles := GetUserRoles(user, app)

	bytes, _ := json.Marshal(roles)
	rolesString := string(bytes)

	Cache.Set(cacheKey, rolesString, cache.DefaultExpiration)
	return rolesString

}

// addRolesToUserInfoString Добавляет поле со списком ролей пользователя к информации о пользователе.
func addRolesToUserInfoString(userInfoString, rolesString string) string {
	info := make(map[string]interface{})
	roles := make([]string, 0)

	_ = json.UnmarshalFromString(userInfoString, &info)
	_ = json.UnmarshalFromString(rolesString, &roles)

	info["roles"] = roles
	s, _ := json.MarshalToString(info)
	return s
}

// GetUserInfoString возвращает сериализованную информацию о пользователе, включая его роли в приложении.
func GetUserInfoString(user, app string) string {

	userInfoString := ""
	rolesString := ""
	userInfoStringWithRoles := ""

	// читаем из кэша
	cacheKey := user + "-info"
	cachedValue, found := Cache.Get(cacheKey)
	if found {
		fmt.Println("Cached:", cacheKey, "=", cachedValue)
		// return cachedValue.(string)
		userInfoString = cachedValue.(string)
	} else {
		fmt.Println("NOT cached:", cacheKey)
		// обновляем кэш
		record, err := db.QueryRowMap(`SELECT id, username, email, fullname, description FROM "user" WHERE username = $1;`, user)
		if err != nil {
			return ""
		}
		jsonBytes, _ := json.Marshal(record)
		userInfoString = string(jsonBytes)
		Cache.Set(cacheKey, userInfoString, cache.DefaultExpiration)
	}

	// FIXME: из соображений производительности нужно возвращать информацию без ролей
	// возвращаем информацию о пользователе
	rolesString = GetUserRolesString(user, app)
	userInfoStringWithRoles = addRolesToUserInfoString(userInfoString, rolesString)
	fmt.Println("userInfoStringWithRoles=", userInfoStringWithRoles)
	fmt.Println("------------------------------------------------------")
	return userInfoStringWithRoles
}

// AppUserRoleExist проверят наличие связки appname-username-rolename в таблице app_user_role.
func AppUserRoleExist(appname, username, rolename string) bool {
	_, err := db.QueryRowMap(`SELECT * FROM  app_user_role  
		WHERE appname = $1 AND username = $2 AND rolename = $3 ;`,
		appname, username, rolename)
	return err == nil
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

// IsRequestToAppSigned  true если приложение доступно для пользователей без роли
func IsRequestToAppSigned(appName string) bool {
	cacheKey := "is-request-to-" + appName + "-signed"
	cachedValue, found := Cache.Get(cacheKey)
	if found {
		fmt.Println("cached ", cacheKey, "=", cachedValue)
		return cachedValue.(bool)
	}

	_, err := db.QueryRowMap(`SELECT * FROM "app" WHERE appname = $1 AND sign = 'Y' ;`, appName)
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

func GetUserNameByEmail(email string) string {
	rec, err := db.QueryRowMap(` SELECT username FROM "user" WHERE email=$1`, email)
	if err != nil {
		return ""
	}
	return rec["username"].(string)
}

// ListPublicApps Возвращает список публичных приложений не требующих авторизации пользователя.
func ListPublicApps() []map[string]interface{} {
	records, err := db.QuerySliceMap(`SELECT appname,description FROM app WHERE public='Y'`)
	if err != nil {
		fmt.Println(err)
	}
	return records
}
