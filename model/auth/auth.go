package auth

import (
	"auth-proxy/model/db"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	// "net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	// gsessions "github.com/gorilla/sessions"
)

// var Store = gsessions.NewCookieStore([]byte("secret"))

var Cache = cache.New(2*time.Minute, 4*time.Minute)

// GetUserName возвращает имя текущего пользователя или пустую строку
func GetUserName(c *gin.Context) string {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		return ""
	} else {
		return user.(string)
	}
}

// DeleteSession удаляет текущую сессию (стирает куки на стороне клиента)
func DeleteSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
}

// func DeleteSessionR(w http.ResponseWriter, r *http.Request) {
// 	session, err := auth.Store.Get(c.Request, "auth-proxy")
// 	delete(session.Values, "user")
// 	session.Options = &gsessions.Options{MaxAge: -1}
// 	session.Save(r, w)
// }

// func GetUserNameR(r *http.Request) (userName string) {
// 	session, err := Store.Get(r, "auth-proxy")
// 	if err != nil {
// 		return
// 	} else {
// 		userName, _ := session.Values["user"].(string)
// 	}
// }

// CheckUserPassword проверяет пароль пользователя. Возвращает true если проверка прошла успешно.
func CheckUserPassword(username, password string) bool {
	record, err := db.QueryRowMap(`SELECT password FROM "user" WHERE username = $1;`, username)
	if err != nil {
		return false
	}
	dbPassword, ok := record["password"].(string)
	if !ok {
		log.Println("Cannot get password")
		return false
	}
	if dbPassword == password {
		return true
	}
	return false
}

// GetUserRoles возвращает строку с сериализованным масссивом ролей пользователя в заданном приложении.
func GetUserRoles(user, app string) string {
	cacheKey := user + "-" + app + "-roles"
	cachedValue, found := Cache.Get(cacheKey)
	if found {
		fmt.Println("cached roles=", cachedValue)
		return cachedValue.(string)
	}

	// record, err := db.QueryRowMap(`SELECT get_app_user_roles($1,$2) AS roles;`, app, user)
	// SELECT DISTINCT roles FROM user_roles WHERE username='vadim' AND appname = 'app1';
	record, err := db.QueryRowMap(`SELECT json_agg(rolename) AS roles FROM app_user_role WHERE  appname  = $1 AND username = $2 `, app, user)
	if err != nil {
		return ""
	}
	bytes, _ := record["roles"].([]byte)
	roles := string(bytes)

	Cache.Set(cacheKey, roles, cache.DefaultExpiration)
	return roles
}

// GetUserInfo возвращает сериализованную информацию о пользователе
func GetUserInfo(user string) string {
	cacheKey := user + "-info"
	cachedValue, found := Cache.Get(cacheKey)
	if found {
		fmt.Println("cached info=", cachedValue)
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
	return (err != nil)
}

// Logout разлогинить текущего пользователя
func Logout(c *gin.Context) {
	DeleteSession(c)
}

// Login залогинить пользователя
func Login(c *gin.Context, username, password string) error {
	session := sessions.Default(c)
	if CheckUserPassword(username, password) {
		session.Set("user", username)
		// session.Options(sessions.Options{MaxAge: 0})
		err := session.Save()
		if err != nil {
			log.Println("ERROR: Login():session.Set()", err)
		}
		return err
	} else {
		return errors.New("Authentication failed")
	}
}
