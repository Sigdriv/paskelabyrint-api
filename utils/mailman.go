package utils

import (
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

func SendEmail(to string, from string, subject string, body string) (err error) {
	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(from, "Kor ska oss reis"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer("smtp.mail.me.com", 587, "sigdriv@icloud.com", "fkum-tfhk-kalb-smhy")

	err = dialer.DialAndSend(m)
	if err != nil {
		err = fmt.Errorf("Error sending email >> %v", err)
		return
	}

	return
}
