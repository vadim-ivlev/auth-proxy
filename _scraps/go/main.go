package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	// Логинимся
	loginQuery := `
	query { 
		login(
			username: "max",
			password: ""
		)
	}	
	`
	cookie := login(loginQuery)
	log.Println("COOCKIE=", cookie)

	// GraphQL запрос на получение всех пользователей
	listUsersQuery := `
	query {
		list_user {
			length
			list {
			  description
			  disabled
			  email
			  fullname
			  username
			}
		  }
		}
	`
	request, _ := json.Marshal(map[string]string{"query": listUsersQuery})

	// отправляем запрос
	req, _ := http.NewRequest("POST", "https://auth-proxy.rg.ru/schema", bytes.NewBuffer(request))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", cookie) //  <---   вставляем куки
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// печатаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))

}

// login Логинится и в случае успеха возвращает непустую строку кук
func login(q string) (cookie string) {
	// отправляем запрос на сервер
	request, _ := json.Marshal(map[string]string{"query": q})
	resp, err := http.Post("https://auth-proxy.rg.ru/schema", "application/json", bytes.NewBuffer(request))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// печатаем ответ сервера
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))

	// возвращаем куки
	cookie = resp.Header.Get("Set-Cookie")

	return cookie
}
