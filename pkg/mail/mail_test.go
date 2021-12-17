package mail

import (
	"fmt"
	"net/url"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// ReadConfig("../../configs/mail.yaml", "dev")
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func Test_url_scaping(t *testing.T) {
	fmt.Println("Original   :", "vadim.ivlev@rg.ru & fff=привет ?")
	fmt.Println("QueryEscape :", url.QueryEscape("vadim.ivlev@rg.ru & fff=привет ?")) // use it for query params
	fmt.Println("PathEscape :", url.PathEscape("vadim.ivlev@rg.ru & fff=привет ?"))
}
