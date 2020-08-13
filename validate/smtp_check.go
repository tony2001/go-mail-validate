package validate

import (
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"time"

	"github.com/tony2001/go-mail-validate/config"
	"github.com/tony2001/go-mail-validate/log"
)

func trySmtp(hostStr string, localStr string, domainStr string, strictCheck bool) (bool, error) {
	host := fmt.Sprintf("%s:%d", hostStr, 25)
	log.Debugf("trying to connect to %s", host)

	smtpTimeout := time.Duration(config.GetSmtpTimeout()) * time.Millisecond
	conn, err := net.DialTimeout("tcp", host, smtpTimeout)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return false, err
	}
	defer client.Close()

	log.Debugf("sending HELLO to %s", host)
	err = client.Hello(domainStr)
	if err != nil {
		return false, err
	}

	defer client.Reset()
	defer client.Quit()

	log.Debugf("sending MAIL to %s", host)
	fullEmail := fmt.Sprintf("%s@%s", localStr, domainStr)

	//try to send email from the same mail address
	err = client.Mail(fullEmail)
	if err != nil {
		return true, err
	}

	log.Debugf("sending RCPT to %s", host)
	err = client.Rcpt(fullEmail)
	if err != nil {
		errStr := err.Error()
		if len(errStr) < 3 {
			return true, err
		}

		//trying to be smart here
		errCode, errAtoi := strconv.Atoi(errStr[:3])
		if errAtoi != nil {
			return true, err
		}

		if errCode >= 300 {
			//server responded with a valid permanent error
			return false, err
		}

		return true, err
	}
	return true, nil
}
