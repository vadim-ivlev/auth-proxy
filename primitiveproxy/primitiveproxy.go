// usage go run main.go
// http://localhost:12345/https://rg.ru
// https://gist.github.com/fabrizioc1/4327250

package primitiveproxy

import (
	"io"
	"net/http"
	"strings"
)

type PrimitiveProxy struct {
	url string
}

func NewPrimitiveProxy(url string) *PrimitiveProxy { return &PrimitiveProxy{url} }

func (p *PrimitiveProxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	// url := p.url
	// url := p.url + r.RequestURI
	url := p.url + r.URL.Path //+ "?" +r.URL.RawQuery
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}
	if p.url == "" {
		url = strings.TrimPrefix(r.RequestURI, "/")
	}

	req, err := http.NewRequest(r.Method, url, r.Body)

	for name, value := range r.Header {
		req.Header.Set(name, value[0])
	}

	resp, err := client.Do(req)
	r.Body.Close()

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}

	wr.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(wr, resp.Body)
	resp.Body.Close()

}

// func main() {
// 	proxy := &PrimitiveProxy{}
// 	err := http.ListenAndServe(":12345", proxy)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err.Error())
// 	}
// }
