package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"github.com/Zhoangp/Auth-Service/config"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"path/filepath"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

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

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}
func SendToken(cf *config.Config, destMail string, data, name string, path string) error {
	from := cf.Email.AppEmail
	password := cf.Email.AppPassword

	dataEmail := EmailData{
		URL:       path + data,
		FirstName: name,
		Subject:   "Verify Account",
	}

	var body bytes.Buffer
	template, err := ParseTemplateDir("pkg/templates")
	if err != nil {
		return err
	}

	if err := template.ExecuteTemplate(&body, "verify.html", &dataEmail); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", destMail)
	m.SetHeader("Subject", dataEmail.Subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil

}
func GenerateOTP(length int) (otp string, err error) {
	b := make([]byte, length)
	if _, err = rand.Read(b); err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		otp += string(otpChars[int(b[i])%10])
	}
	return
}

const otpChars = "1234567890"
