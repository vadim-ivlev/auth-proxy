package session

import (
	"log"

	// "github.com/gin-gonic/contrib/sessions"
	// gsessions "github.com/gorilla/sessions"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

// SetVariable устанавливает значение переменной сессии
func SetVariable(c *gin.Context, name, value string) error {
	session := sessions.Default(c)
	session.Set(name, value)
	// session.Options(sessions.Options{MaxAge: 0})
	err := session.Save()
	if err != nil {
		log.Println("ERROR: Login():session.Set()", err)
	}
	return err
}

// GetVariable возвращает значение переменной сессии пользователя или пустую строку
func GetVariable(c *gin.Context, varname string) string {
	session := sessions.Default(c)
	varvalue := session.Get(varname)
	if varvalue == nil {
		return ""
	} else {
		return varvalue.(string)
	}
}

// GetUserName возвращает имя текущего пользователя или пустую строку
func GetUserName(c *gin.Context) string {
	return GetVariable(c, "user")
}

// GetCaptcha возвращает капчу пользователя или пустую строку
func GetCaptcha(c *gin.Context) string {
	return GetVariable(c, "captcha")
}

// DeleteSession удаляет текущую сессию (стирает куки на стороне клиента)
func DeleteSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
}
