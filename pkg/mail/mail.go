package mail

import (
	"auth-proxy/pkg/app"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"net/url"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Header struct {
	From    string
	To      string
	Subject string
}

type MailData struct {
	Header Header
	TMPL   *template.Template
	Data   interface{}
}

type NewUserData struct {
	UserName string
	Link     string
}

// var Params connectionParams
var mailTemplates map[string]string

var mailHtmlTemplates = make(map[string]*template.Template)

// ReadMailTemplate reads YAML file with mail templates
func ReadMailTemplate(fileName string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("mail.ReadMailTemplate() ReadFile error:", err)
		return
	}

	mailTemplates = make(map[string]string)
	err = yaml.Unmarshal(yamlFile, &mailTemplates)
	if err != nil {
		log.Println("mail.ReadMailTemplate() Unmarshal error:", err)
	}
	// переопределяем получение шаблонов писем из файлов
	mailHtmlTemplates["new_user"] = getMailTmpl(app.Params.MailTmplPath + "/new_user.html")
	mailHtmlTemplates["reset_password"] = getMailTmpl(app.Params.MailTmplPath + "/reset_password.html")
}

func (m *MailData) ComposeTmpl() (string, error) {

	var buffer bytes.Buffer

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	buffer.WriteString(fmt.Sprintf(
		"Subject: %s\nFrom: %s\nTo: %s\n%s",
		m.Header.Subject,
		m.Header.From,
		m.Header.To,
		mime,
	))
	err := m.TMPL.Execute(&buffer, m.Data)
	if err != nil {
		fmt.Printf("Failed execute template: %s\n", err)
		return "", err
	}
	return buffer.String(), nil
}

func getMailTmpl(mailTemplateFilePath string) *template.Template {

	mailTemplate, err := ioutil.ReadFile(mailTemplateFilePath)
	if err != nil {
		fmt.Printf("Failed read template file: %s\n", err)
	}

	tmpl, err := template.New("NewUser").Parse(string(mailTemplate))
	if err != nil {
		fmt.Printf("Failed parsing template %s\n", err)
	}

	return tmpl
}

func SendNewUserEmail(toEmail, userName, emailhash string) error {
	if tmpl, ok := mailHtmlTemplates["new_user"]; ok {
		return sendNewUserEmail(toEmail, userName, emailhash, tmpl)
	}
	return errors.New("tmpl for new user not found")
}

func sendNewUserEmail(toEmail, userName, emailhash string, tmpl *template.Template) error {
	entryPoint := fmt.Sprintf("%s&email=%s", app.Params.EntryPoint, toEmail)
	urlParams := fmt.Sprintf("emailhash=%s&email=%s&entry_point=%s", emailhash, url.QueryEscape(toEmail), entryPoint)

	mailData := MailData{
		Header: Header{
			Subject: convertSubject("Регистрация нового пользователя на портале \"Российской газеты\""),
			From:    app.Params.From,
			To:      toEmail,
		},
		TMPL: tmpl,
		Data: NewUserData{
			Link:     app.Params.AdminAPI + "/confirm-email" + "?" + urlParams,
			UserName: userName,
		},
	}
	msg, err := mailData.ComposeTmpl()
	if err != nil {
		return err
	}

	if msg == "" {
		return errors.New("new user registration: failed to send mail, message is empty")
	}
	return sendMail(app.Params.From, toEmail, msg)
}

func SendResetPasswordEmail(toEmail, pageAddress string) error {
	if tmpl, ok := mailHtmlTemplates["reset_password"]; ok {
		return sendResetPasswordEmail(toEmail, pageAddress, tmpl)
	}
	return errors.New("tmpl for reset password not found")
}

func sendResetPasswordEmail(toEmail, pageAddress string, tmpl *template.Template) error {
	// msg := fmt.Sprintf(mailTemplates["reset_password"], app.Params.From, toEmail, pageAddress)

	// entryPoint := fmt.Sprintf("%s&email=%s", app.Params.EntryPoint, toEmail)
	// urlParams := fmt.Sprintf("emailhash=%s&email=%s&entry_point=%s", emailhash, url.QueryEscape(toEmail), entryPoint)

	mailData := MailData{
		Header: Header{
			Subject: convertSubject("Восстановление пароля на портале \"Российской газеты\""),
			From:    app.Params.From,
			To:      toEmail,
		},
		TMPL: tmpl,
		Data: NewUserData{
			Link: pageAddress,
			// UserName: userName,
		},
	}
	msg, err := mailData.ComposeTmpl()
	if err != nil {
		return err
	}

	if msg == "" {
		return errors.New("reset password: failed to send mail, message is empty")
	}

	return sendMail(app.Params.From, toEmail, msg)
}

func SendAuthenticatorEmail(toEmail, pageAddress string) error {
	msg := fmt.Sprintf(mailTemplates["reset_authenticator"], app.Params.From, toEmail, pageAddress)
	return sendMail(app.Params.From, toEmail, msg)
}

// Так как наши почторые сервера не поддерживают кириллицу, после обновления go 1.16 отправка по-умолчанию не работает
// https://go.dev/doc/go1.16#net/smtp
// https://ncona.com/2011/06/using-utf-8-characters-on-an-e-mail-subject/
func convertSubject(subject string) string {
	return "=?utf-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?="
}

// sendMail toEmail может содержать несколько адресов через запятую.
// msg должно быть отформатировано специальным образом.
func sendMail(fromEmail, toEmail, msg string) error {
	if toEmail == "" {
		return errors.New("email address is required")
	}

	// Connect to the remote SMTP server.
	c, err := smtp.Dial(app.Params.SmtpAddress)
	if err != nil {
		return err
	}
	defer c.Close()

	// Set the sender and recipient first
	if err := c.Mail(fromEmail); err != nil {
		return err
	}
	if err := c.Rcpt(toEmail); err != nil {
		return err
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	// _, err = fmt.Fprint(wc, msg)
	// if err != nil {
	// 	return err
	// }

	if _, err = wc.Write([]byte(msg)); err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		return err
	}

	return nil
}

// SendEmailTo отправляет письмо
// от имени fromEmail на адрес toEmail, c указанным темой и телом.
func SendEmailTo(fromEmail, toEmail, subject, body string) error {

	from := removeBreaks(fromEmail)
	to := removeBreaks(toEmail)
	subj := removeBreaks(subject)

	msg := fmt.Sprintf("From: %s \nTo: %s \nSubject: %s \n\n%s\n",
		from, to, subj, body)

	return sendMail(from, to, msg)
}

// removeBreaks очищает строку от переводов строк и пробелов в начале и конце
func removeBreaks(email string) string {
	return strings.TrimSpace(strings.ReplaceAll(email, "\n", ""))
}
