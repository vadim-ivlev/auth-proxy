// usage go run main.go
// http://localhost:12345/https://rg.ru
// https://gist.github.com/fabrizioc1/4327250

package primitiveproxy

import (
	"auth-proxy/pkg/auth"
	"auth-proxy/pkg/signature"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// PrimitiveProxy Структура для хранения параметров прокси.
// auth-proxy хранит map таких структур индексированную по имени приложения/ маршруту проксирования
type PrimitiveProxy struct {
	// куда проксировать запросы
	Url string
	// префикс проксирования и дновременно имя приложения в терминах auth-proxy
	appname string
	Rebase  string
	XToken  string
}

// NewPrimitiveProxy Возвращает указатель на PrimitiveProxy
func NewPrimitiveProxy(url, appname, rebase, xtoken string) *PrimitiveProxy {
	return &PrimitiveProxy{url, appname, rebase, xtoken}
}

func (p *PrimitiveProxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	// строим url из маршрута и строки запроса
	url := p.Url + r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	// Создаем новый запрос с методом и телом исходного запроса
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		fmt.Println("Ошибка создания http.NewRequest.", err.Error())
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
	// Устанавливаем длину контента
	req.ContentLength = r.ContentLength

	// Копируем заголовки из исходного запроса в новый
	for name, value := range r.Header {
		if p.Rebase == "Y" && name == "Accept-Encoding" {
			continue
		}
		req.Header.Set(name, value[0])
		if strings.TrimSpace(p.XToken) != "" {
			req.Header.Set("X-AuthProxy-Token", p.XToken)
		}
	}

	// Подписываем запрос
	if auth.IsRequestToAppSigned(p.appname) {
		err := signature.Sign(req)
		if err != nil {
			fmt.Println("Signing error:", err)
		}
	}

	// Исполняем новый запрос
	resp, err := client.Do(req)
	r.Body.Close()

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	// копируем заголовки ответа в выходной поток
	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}

	if p.Rebase == "Y" {
		replaceAbsoluteLinks(wr, r, resp, p)
	}

	// копируем тело ответа в выходной поток
	wr.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(wr, resp.Body)
	resp.Body.Close()

}

// replaceAbsoluteLinks - Попытка изменить абсолютные ссылки в HTML, CSS, JS  на относительные и <base>.
// Делается для корректного отображения страниц на некоторых сайтах.
// Правильный способ -  не использовать ссылок наxинающихся с /
// и не использовать эту функцию.
func replaceAbsoluteLinks(wr http.ResponseWriter, r *http.Request, resp *http.Response, p *PrimitiveProxy) {
	contentType := resp.Header.Get("Content-Type")
	isHTML := strings.Contains(contentType, "text/html")
	isCSS := strings.Contains(contentType, "text/css")
	isJS := strings.Contains(contentType, "application/javascript")

	if isHTML || isCSS || isJS {

		// вычисляем path & base
		path := "/apps/" + p.appname + "/"
		base := r.Host + path
		if strings.HasPrefix(r.Proto, "HTTP/") {
			base = "http://" + base
		} else {
			base = "https://" + base
		}

		// устанавливаем HTTP base. Похожене оказывает влияния в chrome.
		wr.Header().Set("Content-Base", base)
		wr.Header().Set("Content-Location", base)

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		bodyString := string(bodyBytes)

		if isHTML {
			// bodyString = strings.ReplaceAll(bodyString, "<head>", `<head><base href="`+base+`">)
			bodyString = strings.ReplaceAll(bodyString, `href="//`, `href="||`)
			bodyString = strings.ReplaceAll(bodyString, `href='//`, `href='||`)
			bodyString = strings.ReplaceAll(bodyString, `src="//`, `src="||`)
			bodyString = strings.ReplaceAll(bodyString, `src='//`, `src='||`)

			bodyString = strings.ReplaceAll(bodyString, `href="/`, `href="`+path)
			bodyString = strings.ReplaceAll(bodyString, `href='/`, `href='`+path)
			bodyString = strings.ReplaceAll(bodyString, `src="/`, `src="`+path)
			bodyString = strings.ReplaceAll(bodyString, `src='/`, `src='`+path)

			bodyString = strings.ReplaceAll(bodyString, `href="||`, `href="//`)
			bodyString = strings.ReplaceAll(bodyString, `href='||`, `href='//`)
			bodyString = strings.ReplaceAll(bodyString, `src="||`, `src="//`)
			bodyString = strings.ReplaceAll(bodyString, `src='||`, `src='//`)
		}
		if isCSS {
			bodyString = strings.ReplaceAll(bodyString, `url(/`, `url(`+path)
			bodyString = strings.ReplaceAll(bodyString, `url("/`, `url("`+path)
			bodyString = strings.ReplaceAll(bodyString, `url('/`, `url('`+path)
		}
		if isJS {
			bodyString = strings.ReplaceAll(bodyString, `sourceMappingURL=`, `sourceMappingURL=`+path)
		}
		fmt.Fprint(wr, bodyString)
	}
}
