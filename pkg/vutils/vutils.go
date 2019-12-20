package vutils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func PrintRequestHeaders(r *http.Request) {
	fmt.Println("--------------HTTP HEADERS------------------------")
	for name, value := range r.Header {
		fmt.Println(name, ":", value[0])
	}
	fmt.Println("--------------END HTTP HEADERS--------------------")
}

// https://stackoverflow.com/questions/43021058/golang-read-request-body
// https://stackoverflow.com/questions/29746123/convert-byte-slice-to-io-reader

func PrintRequestBody(req *http.Request) *http.Request {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return req
	}
	s := string(body)
	defer req.Body.Close() //FIXME: where to put it ??????
	fmt.Println("--------------HTTP BODY------------------------")
	l := len(s)
	if l > 80 {
		fmt.Println(s[0:80])
	} else {
		fmt.Println(s)
	}
	fmt.Println("--------------END HTTP BODY--------------------")
	// And now set a new body, which will simulate the same data we read:
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return req
}
