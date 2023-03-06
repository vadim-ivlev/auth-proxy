package mail

import (
	"auth-proxy/pkg/app"
	"fmt"
	"net/url"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// ReadConfig("../../configs/mail.yaml", "dev")
	// call flag.Parse() here if TestMain uses flags

	// конфиг общих параметров приложения.
	app.ReadEnvConfig("../../configs/app.env")
	os.Exit(m.Run())
}

func Test_url_scaping(t *testing.T) {
	fmt.Println("Original   :", "vadim.ivlev@rg.ru & fff=привет ?")
	fmt.Println("QueryEscape :", url.QueryEscape("vadim.ivlev@rg.ru & fff=привет ?")) // use it for query params
	fmt.Println("PathEscape :", url.PathEscape("vadim.ivlev@rg.ru & fff=привет ?"))
}

func TestComposeTmpl(t *testing.T) {

	tmpl := getMailTmpl("../../templates/mail_lk/new_user.html")

	data := NewUserData{
		Link: "https://rg.ru/login?test",
	}

	mailData := MailData{
		Header: Header{
			Subject: "Регистрация нового пользователя на портале RG.RU",
			From:    app.Params.From,
			To:      "dev@rg.ru",
		},
		TMPL: tmpl,
		Data: data,
	}
	msg, err := mailData.ComposeTmpl()
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(msg)
}

func TestSendNewUserEmail(t *testing.T) {

	tmpl := getMailTmpl("../../templates/mail/new_user.html")

	err := sendNewUserEmail("dev@rg.ru", "Иван Иванов", "test", tmpl)
	if err != nil {
		t.Errorf("%s", err)
	}

}
