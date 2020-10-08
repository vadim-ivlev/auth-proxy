package server

import (
	"log"
	"net/http"

	// "github.com/gin-gonic/contrib/sessions"
	// gsessions "github.com/gorilla/sessions"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

// SetSessionVariable устанавливает значение переменной сессии
func SetSessionVariable(c *gin.Context, name, value string) error {
	session := sessions.Default(c)
	session.Set(name, value)
	err := session.Save()
	if err != nil {
		log.Println("SetSessionVariable()", err)
	}
	return err
}

// GetSessionVariable возвращает значение переменной сессии пользователя или пустую строку
func GetSessionVariable(c *gin.Context, varname string) string {
	session := sessions.Default(c)
	varvalue := session.Get(varname)
	if varvalue == nil {
		return ""
	}
	return varvalue.(string)
}

// DeleteSession удаляет текущую сессию (стирает куки на стороне клиента)
func DeleteSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{
		MaxAge:   -1,
		Secure:   SecureCookie,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})
	session.Save()
}
