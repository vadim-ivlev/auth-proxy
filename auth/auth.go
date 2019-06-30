package auth

func CheckUserPassword(user, password string) bool {
	if user == "q" && password == "q" {
		return true
	}
	return false
}

func GetUserRoles(user, app string) string {
	println(app)
	return `["admin", "editor"]`
}

func GetUserInfo(user string) string {
	return `{"name":"Petr", "email":"petr@rg.ru"}`
}
