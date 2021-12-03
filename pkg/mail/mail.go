package mail

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"

	"gopkg.in/yaml.v2"
)

type connectionParams struct {
	Addr string
	From string
	Link string
}

var envName string
var params connectionParams
var mailTemplates map[string]string

// ReadConfig reads YAML file
func ReadConfig(fileName string, env string) {
	envName = env
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}

	envParams := make(map[string]connectionParams)
	err = yaml.Unmarshal(yamlFile, &envParams)
	if err != nil {
		log.Println("Mail ReadConfig() error:", err)
	}
	params = envParams[env]
}

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
}

// SendMessage send mail
// TODO: https://stackoverflow.com/questions/46805579/send-smtp-email-to-multiple-receivers
// https://github.com/golang/go/wiki/SendingMail
func SendMessage(templateName string, username, toMail, password string) error {
	if toMail == "" {
		return errors.New("email address is required")
	}

	// Connect to the remote SMTP server.
	c, err := smtp.Dial(params.Addr)
	if err != nil {
		return err
	}
	defer c.Close()

	// Set the sender and recipient first
	if err := c.Mail(params.From); err != nil {
		return err
	}
	if err := c.Rcpt(toMail); err != nil {
		return err
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	msg := fmt.Sprintf(mailTemplates[templateName], params.From, toMail, params.Link, username, password)
	_, err = fmt.Fprintf(wc, msg)
	if err != nil {
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

// ------------- new code -----------------------------------------------------------------------------

func SendNewUserEmail(username, toEmail, password string) error {
	msg := fmt.Sprintf(mailTemplates["new_user"], params.From, toEmail, params.Link, username, password)
	return sendMail(params.From, toEmail, msg)
}

func SendNewPasswordEmail(username, toEmail, password string) error {
	msg := fmt.Sprintf(mailTemplates["new_password"], params.From, toEmail, params.Link, username, password)
	return sendMail(params.From, toEmail, msg)
}

func SendResetPasswordEmail(toEmail, pageAddress string) error {
	msg := fmt.Sprintf(mailTemplates["reset_password"], params.From, toEmail, pageAddress)
	return sendMail(params.From, toEmail, msg)
}

func SendAuthenticatorEmail(toEmail, pageAddress string) error {
	msg := fmt.Sprintf(mailTemplates["reset_authenticator"], params.From, toEmail, pageAddress)
	return sendMail(params.From, toEmail, msg)
}

// func sendMail(fromEmail, toEmail, msg string) error {
// 	// if envName == "dev" {
// 	// 	return sendSecureMail(fromEmail, toEmail, msg)
// 	// } else {
// 	// 	return sendInsecureMail(fromEmail, toEmail, msg)
// 	// }

// 	return sendInsecureMail(fromEmail, toEmail, msg)
// }

// // sendSecureMail toEmail может содержать несколько адресов через запятую.
// // msg должно быть отформатировано специальным образом.
// func sendSecureMail(fromEmail, toEmail, msg string) error {
// 	if toEmail == "" {
// 		return errors.New("email должен быть не пустым")
// 	}
// 	username := "vadim.ivlev@gmail.com"
// 	password := os.Getenv("GMAIL_PASSWORD")
// 	host := "smtp.gmail.com"
// 	port := ":587"

// 	toEmailsArray := strings.Split(toEmail, ",")

// 	fmt.Println("MAIL_PASSWORD=", password, fromEmail, toEmailsArray, msg)
// 	err := smtp.SendMail(host+port, smtp.PlainAuth("", username, password, host), fromEmail, toEmailsArray, []byte(msg))
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return err
// }

// sendMail toEmail может содержать несколько адресов через запятую.
// msg должно быть отформатировано специальным образом.
func sendMail(fromEmail, toEmail, msg string) error {
	if toEmail == "" {
		return errors.New("email address is required")
	}

	// Connect to the remote SMTP server.
	c, err := smtp.Dial(params.Addr)
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

	_, err = fmt.Fprint(wc, msg)
	if err != nil {
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
