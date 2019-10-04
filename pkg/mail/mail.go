// https://github.com/golang/go/wiki/SendingMail
package mail

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"

	"gopkg.in/yaml.v2"
)

type connectionParams struct {
	Addr string
	From string
}

var params connectionParams
var mailTemplates map[string]string

// **********************************************************************************
// ReadConfig reads YAML file
func ReadConfig(fileName string, env string) {
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

// ReadConfig reads YAML file with mail templates
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

// TODO: https://stackoverflow.com/questions/46805579/send-smtp-email-to-multiple-receivers
func SendMessage(templateName string, username, toMail, password string) error {
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

	msg := fmt.Sprintf(mailTemplates[templateName], params.From, toMail, username, password)
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
