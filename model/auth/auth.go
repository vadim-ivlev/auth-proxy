package auth

import (
	"auth-proxy/model/db"
	"encoding/json"
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

func GetUserName(c *gin.Context) string {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		return ""
	} else {
		return user.(string)
	}
}

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

func CheckUserPassword(username, password string) bool {
	record, err := db.QueryRowMap("SELECT password FROM public.user WHERE username = $1;", username)
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

func GetUserRoles(user, app string) string {
	cacheKey := user + "-" + app + "-roles"
	cachedValue, found := Cache.Get(cacheKey)
	if found {
		fmt.Println("cached roles=", cachedValue)
		return cachedValue.(string)
	}

	// record, err := db.QueryRowMap(`SELECT public.get_app_user_roles($1,$2) AS roles;`, app, user)
	record, err := db.QueryRowMap(`SELECT json_agg(rolename) AS roles FROM app_user_role WHERE  appname  = $1 AND username = $2 `, app, user)
	if err != nil {
		return ""
	}
	roles := string(record["roles"].([]byte))

	Cache.Set(cacheKey, roles, cache.DefaultExpiration)
	return roles
}

func GetUserInfo(user string) string {
	cacheKey := user + "-info"
	cachedValue, found := Cache.Get(cacheKey)
	if found {
		fmt.Println("cached info=", cachedValue)
		return cachedValue.(string)
	}

	record, err := db.QueryRowMap(`SELECT username, email, fullname, description 
		FROM public.user 
		WHERE username = $1;`, user)
	if err != nil {
		return ""
	}
	jsonBytes, _ := json.Marshal(record)
	jsonString := string(jsonBytes)

	Cache.Set(cacheKey, jsonString, cache.DefaultExpiration)
	return jsonString
}
