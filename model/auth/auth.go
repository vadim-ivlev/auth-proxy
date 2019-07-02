package auth

import (
	"auth-proxy/model/db"
	"encoding/json"
	"log"
)

// import (
// 	gsessions "github.com/gorilla/sessions"
// )

// var Store = gsessions.NewCookieStore([]byte("secret"))

func CheckUserPassword(user, password string) bool {
	record, err := db.QueryRowMap("SELECT password FROM public.user WHERE username = $1;", user)
	if err != nil {
		return false
	}
	savedPassword, ok := record["password"].(string)
	if !ok {
		log.Println("Cannot get password")
		return false
	}
	if savedPassword == password {
		return true
	}
	return false
}

func GetUserRoles(user, app string) string {
	record, err := db.QueryRowMap(`SELECT public.get_app_user_roles($1,$2) AS roles;`, app, user)
	if err != nil {
		return ""
	}
	roles, _ := record["roles"].(string)
	return roles

	// jsonBytes, _ := json.Marshal(record)
	// jsonString := string(jsonBytes)
	// return jsonString
}

func GetUserInfo(user string) string {
	record, err := db.QueryRowMap(`SELECT username, email, fullname, description 
		FROM public.user 
		WHERE username = $1;`, user)
	if err != nil {
		return ""
	}
	jsonBytes, _ := json.Marshal(record)
	jsonString := string(jsonBytes)
	return jsonString
}
