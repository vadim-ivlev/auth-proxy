package server

import (
	"github.com/gin-gonic/gin"
)

// ForkWriter - выполняет немедленный Flush() gin.ResponseWriter
type ForkWriter struct {
	gin.ResponseWriter
}

func (fw ForkWriter) Write(b []byte) (int, error) {
	fw.Flush()
	return fw.ResponseWriter.Write(b)
}

// FlusherMiddleware - выполняет немедленный Flush() gin.ResponseWriter.
// Добавлено для проксирования Stream Events openai/freeai, в режиме "stream: true .
// Issue #112
func FlusherMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fWriter := &ForkWriter{
			ResponseWriter: c.Writer,
		}
		c.Writer = fWriter
		c.Next()
	}
}
