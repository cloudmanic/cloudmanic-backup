package main

import (
	"errors"

	"gopkg.in/gomail.v2"
	"gopkg.in/mailgun/mailgun-go.v1"
)

//
// Pass in everything we need to send an email and we send it.
// If we have a SMTP in our configs we use that if not we use
// Mailgun's library for sending mail.
//
func EmailSend(to string, subject string, html string, text string) error {

	// Are we sending as SMTP or via Mailgun? Typically we
	// send as SMTP for local development so we can use Mailhog
	if cfg.MailDriver == "smtp" {
		return SmtpSend(to, subject, html, text)
	}

	// Send via mailgun
	if cfg.MailDriver == "mailgun" {
		return MailgunSend(to, subject, html, text)
	}

	// We should never get here if we are configured correctly.
	var err = errors.New("No mail driver found.")
	Log("Send() - No mail driver found.")
	return err

}

//
// Send via Mailgun.
//
func MailgunSend(to string, subject string, html string, text string) error {

	// Setup mailgun
	mg := mailgun.NewMailgun(cfg.MailGunDomain, cfg.MailGunApiKey, "")

	// Create message
	message := mailgun.NewMessage("Cloudmanic Backup"+"<"+cfg.MailFromEmail+">", subject, text, to)
	message.SetHtml(html)

	// Send the message
	_, _, err := mg.Send(message)

	if err != nil {
		Log("MailgunSend() - Unable to send email.")
		return err
	}

	// Everything went well!
	return nil

}

//
// Send as SMTP.
//
func SmtpSend(to string, subject string, html string, text string) error {

	// Setup the email to send.
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.MailFromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)
	m.AddAlternative("text/plain", text)

	// Make a SMTP connection
	d := gomail.NewDialer(cfg.MailHost,
		cfg.MailPort,
		cfg.MailUsername,
		cfg.MailPassword)

	// Send Da Email
	if err := d.DialAndSend(m); err != nil {
		Log("SmtpSend() - Unable to send email.")
		return err
	}

	// Everything went well!
	return nil

}

/* End File */
