// usage go run main.go
// http://localhost:12345/https://rg.ru
// https://gist.github.com/fabrizioc1/4327250

package primitiveproxy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type PrimitiveProxy struct {
	url     string
	appname string
	rebase  string
}

func NewPrimitiveProxy(url, appname, rebase string) *PrimitiveProxy {
	return &PrimitiveProxy{url, appname, rebase}
}

func (p *PrimitiveProxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	// строим url из маршрута и строки запроса
	url := p.url + r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	// Создаем новый запрос с методом и телом исходного запроса
	req, err := http.NewRequest(r.Method, url, r.Body)

	// Копируем заголовки из исходного запроса в новый
	for name, value := range r.Header {
		if p.rebase == "Y" && name == "Accept-Encoding" {
			continue
		}
		req.Header.Set(name, value[0])
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

	if p.rebase == "Y" {
		replaceAbsoluteLinks(wr, r, resp, p)
	}

	// копируем тело ответа в выходной поток
	wr.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(wr, resp.Body)
	resp.Body.Close()

}

// replaceAbsoluteLinks - Попытка изменить абсолютные ссылки HTML документа на относительные и <base>.
// Делается для правильного отображения страниц. Правильный способ -  не использовать HTML ссылок наxинающихся с /
// и не использовать эту функцию.
func replaceAbsoluteLinks(wr http.ResponseWriter, r *http.Request, resp *http.Response, p *PrimitiveProxy) {
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		// устанавливаем base
		// wr.Header().Set("Content-Base", "http://localhost:4000/apps/_pravo/")
		// wr.Header().Set("Content-Location", "http://localhost:4000/apps/_pravo/")

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)

		}

		base := r.Host + "/apps/" + p.appname
		if strings.HasPrefix(r.Proto, "HTTP/") {
			base = "http://" + base
		} else {
			base = "https://" + base
		}
		bodyString := string(bodyBytes)
		bodyString = strings.Replace(bodyString, "<head>", `<head><base href="`+base+`/">`, 1)
		bodyString = strings.Replace(bodyString, `href="//`, `href="||`, -1)
		bodyString = strings.Replace(bodyString, `href='//`, `href='||`, -1)
		bodyString = strings.Replace(bodyString, `src="//`, `src="||`, -1)
		bodyString = strings.Replace(bodyString, `src='//`, `src='||`, -1)

		bodyString = strings.Replace(bodyString, `href="/`, `href="`, -1)
		bodyString = strings.Replace(bodyString, `href='/`, `href='`, -1)
		bodyString = strings.Replace(bodyString, `src="/`, `src="`, -1)
		bodyString = strings.Replace(bodyString, `src='/`, `src='`, -1)

		bodyString = strings.Replace(bodyString, `href="||`, `href="//`, -1)
		bodyString = strings.Replace(bodyString, `href='||`, `href='//`, -1)
		bodyString = strings.Replace(bodyString, `src="||`, `src="//`, -1)
		bodyString = strings.Replace(bodyString, `src='||`, `src='//`, -1)
		fmt.Fprint(wr, bodyString)
	}
}
