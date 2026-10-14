package utils

import (
	"maxpark_opp_registration/config"

	"gopkg.in/gomail.v2"
)

func SendEmail(from, fromName, to, subject, htmlContent string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, fromName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlContent)

	d := config.InitializeMailer()

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
