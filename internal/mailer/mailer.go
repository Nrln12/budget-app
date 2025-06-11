package mailer

import (
	"bytes"
	"embed"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"path"
	"strconv"
)

//go:embed templates/*
var templateFS embed.FS

type Mailer struct {
	dialer *gomail.Dialer
	sender string
	Logger echo.Logger
}

type EmailData struct {
	AppName string
	Subject string
	Meta    interface{}
}

func NewMailer(logger echo.Logger) Mailer {
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		logger.Fatal(err)
	}
	mailHost := os.Getenv("MAIL_HOST")
	mailUser := os.Getenv("MAIL_USERNAME")
	mailPass := os.Getenv("MAIL_PASSWORD")

	dialer := gomail.NewDialer(mailHost, mailPort, mailUser, mailPass)
	return Mailer{
		dialer: dialer,
		sender: os.Getenv("MAIL_SENDER"),
		Logger: logger,
	}
}

func (mailer *Mailer) Send(recipient string, templateFile string, data EmailData) error {
	absolutePath := path.Join("templates", templateFile)
	tmpl, err := template.ParseFS(templateFS, absolutePath)
	if err != nil {
		mailer.Logger.Error("Template parse error: ", err)
		return err
	}

	data.AppName = os.Getenv("APP_NAME")

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		mailer.Logger.Error("Template subject error: ", err)
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		mailer.Logger.Error("Template htmlBody error: ", err)
		return err
	}

	gomailMessage := gomail.NewMessage()
	gomailMessage.SetHeader("To", recipient)
	gomailMessage.SetHeader("From", mailer.sender)
	gomailMessage.SetHeader("Subject", subject.String())
	gomailMessage.SetBody("text/html", htmlBody.String())

	err = mailer.dialer.DialAndSend(gomailMessage)
	if err != nil {
		mailer.Logger.Error("Email sending error: ", err)
		return err
	}
	return nil
}
