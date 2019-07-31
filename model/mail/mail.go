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
	Addr     string
	From     string
	Body     string
}

var params connectionParams

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

func SendPassword(toMail, password string) error {
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

	msg := fmt.Sprintf(params.Body, params.From, toMail, password)
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
