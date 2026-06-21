package usecase

import (
	"github.com/Eucastan/eucastanpay/services/user/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailSender interface {
	SendVerificationEmail(string, string) error
	SendResetPasswordEmail(string, string) error
}

type EmailService struct {
	APIKey   string
	AppEmail string
	FromName string
}

func NewEmailService(cfg *config.Config) EmailSender {
	return &EmailService{
		APIKey:   cfg.EmailAPIKey,
		AppEmail: cfg.AppEmail,
		FromName: cfg.FromName,
	}
}

func (e *EmailService) SendVerificationEmail(to, link string) error {
	from := mail.NewEmail(e.FromName, e.AppEmail)
	subject := "Verify your email"

	html := "<h2>Verify Account</h2><p><a href='" + link + "'>Click here</a></p>"

	toEmail := mail.NewEmail("User", to)
	message := mail.NewSingleEmail(from, subject, toEmail, "", html)

	client := sendgrid.NewSendClient(e.APIKey)
	_, err := client.Send(message)
	return err
}

func (e *EmailService) SendResetPasswordEmail(to, link string) error {
	from := mail.NewEmail(e.FromName, e.AppEmail)
	subject := "Password Reset"

	html := "<h2>Reset Password</h2><p><a href='" + link + "'>Click here</a></p>"

	toEmail := mail.NewEmail("User", to)
	message := mail.NewSingleEmail(from, subject, toEmail, "", html)

	client := sendgrid.NewSendClient(e.APIKey)
	_, err := client.Send(message)
	return err
}
