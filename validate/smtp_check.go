package validate

import (
	"fmt"
	"net"
	"net/smtp"
	"time"

	"go-mail-validate/config"
	"go-mail-validate/log"
)

func trySmtp(hostStr string, localStr string, domainStr string) error {
	host := fmt.Sprintf("%s:%d", hostStr, 25)
	log.Debugf("trying to connect to %s", host)

	smtpTimeout := time.Duration(config.GetSmtpTimeout()) * time.Millisecond
	conn, err := net.DialTimeout("tcp", host, smtpTimeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()

	log.Debugf("sending HELLO to %s", host)
	err = client.Hello(domainStr)
	if err != nil {
		return err
	}

	log.Debugf("sending MAIL to %s", host)
	fakeFrom := fmt.Sprintf("hello@%s", domainStr)
	err = client.Mail(fakeFrom)
	if err != nil {
		return err
	}

	log.Debugf("sending RCPT to %s", host)
	fullEmail := fmt.Sprintf("%s@%s", localStr, domainStr)
	err = client.Rcpt(fullEmail)
	if err != nil {
		return err
	}
	client.Reset() // #nosec
	client.Quit()  // #nosec
	return nil
}
