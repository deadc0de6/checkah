// Copyright (c) 2021 deadc0de6

package alert

import (
	"fmt"
	"net/smtp"

	log "github.com/sirupsen/logrus"
)

// Email alert file struct
type Email struct {
	host     string
	port     string
	mailfrom string
	mailto   string
	user     string
	password string
	options  map[string]string
}

func (a *Email) sendAuth(body []byte) error {
	addr := fmt.Sprintf("%s:%s", a.host, a.port)
	log.Debugf("email with auth to %s", addr)
	auth := smtp.PlainAuth("", a.user, a.password, a.host)
	return smtp.SendMail(addr, auth, a.mailfrom, []string{a.mailto}, body)
}

func (a *Email) sendNoAuth(content string) error {
	addr := fmt.Sprintf("%s:%s", a.host, a.port)
	log.Debugf("email plain to %s", addr)
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}

	// send from
	err = c.Mail(a.mailfrom)
	if err != nil {
		return err
	}

	// send to
	err = c.Rcpt(a.mailto)
	if err != nil {
		return err
	}

	// send body
	wc, err := c.Data()
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(wc, content)
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	// and quit
	err = c.Quit()
	return err
}

// Notify notifies
func (a *Email) Notify(content string) error {
	// body
	var body string
	body += fmt.Sprintf("From: %s <%s>\r\n", a.mailfrom, a.mailfrom)
	body += fmt.Sprintf("To: %s\r\n", a.mailto)
	body += fmt.Sprintf("Subject: checkah %s\r\n", "alert")
	body += "\r\n"
	body += content
	body += "\r\n"

	// send email
	if len(a.user) > 0 && len(a.password) > 0 {
		return a.sendAuth([]byte(body))
	}

	return a.sendNoAuth(body)
}

// GetOptions returns this alert options
func (a *Email) GetOptions() map[string]string {
	return a.options
}

// GetDescription returns a description for this alert
func (a *Email) GetDescription() string {
	return fmt.Sprintf("alert to email %s", a.mailto)
}

// NewAlertEmail creates a new file alert instance
func NewAlertEmail(options map[string]string) (*Email, error) {
	host, ok := options["host"]
	if !ok {
		return nil, fmt.Errorf("\"host\" option required")
	}

	port, ok := options["port"]
	if !ok {
		return nil, fmt.Errorf("\"port\" option required")
	}

	mailfrom, ok := options["mailfrom"]
	if !ok {
		return nil, fmt.Errorf("\"mailfrom\" option required")
	}

	mailto, ok := options["mailto"]
	if !ok {
		return nil, fmt.Errorf("\"mailto\" option required")
	}

	user, ok := options["user"]
	if !ok {
		user = ""
	}

	password, ok := options["password"]
	if !ok {
		password = ""
	}

	a := &Email{
		host:     host,
		port:     port,
		mailto:   mailto,
		mailfrom: mailfrom,
		user:     user,
		password: password,
		options:  options,
	}
	return a, nil
}
