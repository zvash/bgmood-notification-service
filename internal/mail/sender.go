package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	email "github.com/xhit/go-simple-mail/v2"
	"github.com/zvash/bgmood-notification-service/internal/util"
	"html/template"
	"log"
	"time"
)

type EmailSender interface {
	SendEmail(
		templateName string,
		variables interface{},
		to []string,
		attachments []string,
	) error
}

type GeneralEmailSender struct {
	Title             string
	From              string
	Host              string
	Port              int
	Password          string
	templateToSubject map[string]string
	disabled          bool
}

func NewGeneralEmailSender(config util.Config) EmailSender {
	sender := &GeneralEmailSender{
		Title:    config.MailSenderTitle,
		From:     config.MailFromAddress,
		Host:     config.MailSMTPServer,
		Port:     config.MailSMTPServerPort,
		Password: config.MailSMTPServerPassword,
		disabled: config.MailDisableSend,
	}
	sender.templateToSubject = map[string]string{
		"verify":         fmt.Sprintf("%s - Verify Your Email Address", config.AppName),
		"reset-password": fmt.Sprintf("%s - Here is Your Password Reset Code", config.AppName),
	}
	return sender
}

func (sender *GeneralEmailSender) SendEmail(templateName string, variables interface{}, to []string, attachments []string) error {
	if sender.disabled {
		log.Printf("Sending %s email to %v is disabled", templateName, to)
		return nil
	}
	path := fmt.Sprintf("template/%s.html", templateName)
	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	var filled bytes.Buffer
	if err := t.Execute(&filled, variables); err != nil {
		return err
	}
	body := filled.String()

	client, err := sender.prepareSMTPClient()
	if err != nil {
		fmt.Println(err)
		return err
	}

	message := email.NewMSG()
	message.SetFrom(fmt.Sprintf("FROM %s <%s>", sender.Title, sender.From)).
		AddTo(to...).
		SetSubject(sender.templateToSubject[templateName]).
		SetBody(email.TextHTML, body)

	message.SetDSN([]email.DSN{email.SUCCESS, email.FAILURE}, false)

	if message.Error != nil {
		return message.Error
	}

	if err := message.Send(client); err != nil {
		return err
	}
	return nil
}

func (sender *GeneralEmailSender) prepareSMTPClient() (*email.SMTPClient, error) {
	server := email.NewSMTPClient()
	server.Host = sender.Host
	server.Port = sender.Port
	server.Username = sender.From
	server.Password = sender.Password

	server.Authentication = email.AuthAuto

	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 30 * time.Second
	server.Encryption = email.EncryptionSTARTTLS
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	smtpClient, err := server.Connect()
	if err != nil {
		return nil, err
	}
	return smtpClient, nil
}
