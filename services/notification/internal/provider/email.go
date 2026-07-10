package provider

import (
	"github.com/Eucastan/eucastanpay/services/notification/config"
	"github.com/Eucastan/eucastanpay/services/notification/internal/domain"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailClient struct {
	*sendgrid.Client
	APIKey   string
	AppEmail string
	FromName string
}

func NewEmailProvider(cfg *config.Config) *EmailClient {
	return &EmailClient{
		Client:   sendgrid.NewSendClient(cfg.EmailAPIKey),
		APIKey:   cfg.EmailAPIKey,
		AppEmail: cfg.AppEmail,
		FromName: cfg.FromName,
	}
}

func (e *EmailClient) SendEmail(n *domain.Notification) error {
	from := mail.NewEmail(e.FromName, e.AppEmail)
	subject := n.Title
	to := mail.NewEmail("User", n.Title)
	msg := mail.NewContent("text/html", n.Message)

	message := mail.NewSingleEmail(from, subject, to, msg.Value, msg.Value)
	_, err := e.Client.Send(message)

	return err
}

func (e *EmailClient) SendPush(n *domain.Notification) error {
	from := mail.NewEmail(e.FromName, e.AppEmail)
	subject := n.Title
	to := mail.NewEmail("User", n.Title)
	msg := mail.NewContent("text/html", n.Message)

	message := mail.NewSingleEmail(from, subject, to, msg.Value, msg.Value)
	_, err := e.Client.Send(message)

	return err
}

func (e *EmailClient) SendInApp(n *domain.Notification) error {
	from := mail.NewEmail(e.FromName, e.AppEmail)
	subject := n.Title
	to := mail.NewEmail("User", n.Title)
	msg := mail.NewContent("text/html", n.Message)

	message := mail.NewSingleEmail(from, subject, to, msg.Value, msg.Value)
	_, err := e.Client.Send(message)

	return err
}
