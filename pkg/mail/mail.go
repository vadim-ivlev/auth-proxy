package mail

import (
	"auth-proxy/pkg/app"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"

	"gopkg.in/yaml.v2"
)

// type connectionParams struct {
// 	SmtpAddress string `yaml:"smtp_address"`
// 	From        string `yaml:"from"`
// }

// var Params connectionParams
var mailTemplates map[string]string

// // ReadConfig reads YAML file
// func ReadConfig(fileName string, env string) {
// 	yamlFile, err := ioutil.ReadFile(fileName)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}

// 	envParams := make(map[string]connectionParams)
// 	err = yaml.Unmarshal(yamlFile, &envParams)
// 	if err != nil {
// 		log.Println("Mail ReadConfig() error:", err)
// 	}
// 	Params = envParams[env]
// }

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

func SendNewUserEmail(toEmail, emailhash string) error {
	urlParams := fmt.Sprintf("emailhash=%s&email=%s&entry_point=%s", emailhash, toEmail, app.Params.EntryPoint)
	msg := fmt.Sprintf(mailTemplates["new_user"], app.Params.From, toEmail, app.Params.ConfirmEmailUrl, urlParams)
	return sendMail(app.Params.From, toEmail, msg)
}

func SendResetPasswordEmail(toEmail, pageAddress string) error {
	msg := fmt.Sprintf(mailTemplates["reset_password"], app.Params.From, toEmail, pageAddress)
	return sendMail(app.Params.From, toEmail, msg)
}

func SendAuthenticatorEmail(toEmail, pageAddress string) error {
	msg := fmt.Sprintf(mailTemplates["reset_authenticator"], app.Params.From, toEmail, pageAddress)
	return sendMail(app.Params.From, toEmail, msg)
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
