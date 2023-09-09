package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"

	"github.com/k3a/html2text"
	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/internal/logger"
	"github.com/s1Sharp/s1-tts-restapi/internal/models"

	"gopkg.in/gomail.v2"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

var (
	log = logger.GetLogger()
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// 👇 Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	log.Println("Parsing templates...")

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

// TODO -> another service
func SendEmail(user *models.DBResponse, data *EmailData, templateName string, config config.Config) error {
	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort, _ := strconv.Atoi(config.SMTPPort)

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template = template.Lookup(templateName)
	err = template.Execute(&body, &data)
	if err != nil {
		return err
	}
	fmt.Println(template.Name())

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
