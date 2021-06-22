package mailer

import (
	"bytes"
	"embed"
	"github.com/go-mail/mail/v2"
	"github.com/knadh/koanf"
	"html/template"
	"log"
	"time"
)

//go:embed "templates"
var templateFS embed.FS

type mailerConf struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Sender   string `koanf:"sender"`
}

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(cfg *koanf.Koanf) Mailer {
	var mailerCfg mailerConf
	if err := cfg.Unmarshal("smtp", &mailerCfg); err != nil {
		log.Fatalf("error loading mailer config: %v\n", err)
	}
	log.Print("initializing mailer")
	dialer := mail.NewDialer(mailerCfg.Host, mailerCfg.Port, mailerCfg.Username, mailerCfg.Password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: mailerCfg.Sender,
	}
}

func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
