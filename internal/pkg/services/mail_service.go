package services

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/gomail.v2"
)

type MailServiceInterface interface {
	SendOneTimePasswordMail(templateName string, data map[string]interface{}) error
}

type MailService struct {
	basePath  string
	templates map[string]*template.Template
}

func NewMailService() MailServiceInterface {
	pwdPath, err := os.Getwd()
	if err != nil {
		return nil
	}

	// Replace backslashes with forward slashes for Windows compatibility
	pwdPath = strings.ReplaceAll(pwdPath, "\\", "/")

	// Create a new instance of JsonSchemaValidator
	mailService := &MailService{
		basePath:  pwdPath,
		templates: make(map[string]*template.Template),
	}

	err = mailService.loadDirTemplates(os.Getenv("TEMPLATE_MAIL_PATH"))
	if err != nil {
		return nil
	}

	return mailService
}

/**
 * loadDirTemplates loads all the templates in the specified directory path.
 * It walks through all files and directories in the path, checks if the file name ends with .html,
 * parses the file as a template, and adds it to the templates map using the file name as the key.
 *
 * Parameters:
 * - path (string): The directory path to load the templates from.
 *
 * Returns:
 * - error: An error if any occurred during the loading process, otherwise nil.
 */
func (m *MailService) loadDirTemplates(path string) error {
	// Walk through all files and directories in the specified directory path
	err := filepath.Walk(m.basePath+path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current file or directory is a directory
		if f.IsDir() {
			return nil
		}

		// Check if the file name ends with .html
		if !strings.HasSuffix(f.Name(), ".html") {
			return nil
		}

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return err
		}
		// Add the schema to the schemas map using the file name as the key
		m.templates[f.Name()] = tmpl
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *MailService) SendOneTimePasswordMail(templateName string, data map[string]interface{}) error {
	mail := gomail.NewMessage()

	var body bytes.Buffer
	err := m.templates[templateName].Execute(&body, data)
	if err != nil {
		return err
	}

	mail.SetHeader("From", os.Getenv("EMAIL_ACCOUNT"))
	mail.SetHeader("To", data["to"].(string))
	mail.SetHeader("Subject", "Password Change Notification")
	mail.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_ACCOUNT"), os.Getenv("EMAIL_PASS"))
	if err := d.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}
