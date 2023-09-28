// usage:
// go test -v ./primitiveproxy
// and open http://localhost:12345/https://rg.ru

package primitiveproxy

import (
	"log"
	"net/http"
	"testing"
)

func TestSetGetDelete(t *testing.T) {
	url := "https://rg.ru"
	if url == "" {
		log.Println("Open http://localhost:12345/https://rg.ru")
	} else {
		log.Println("Open http://localhost:12345/")
	}

	proxy := NewPrimitiveProxy(url, "rg", "", "xxx")
	err := http.ListenAndServe(":12345", proxy)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
