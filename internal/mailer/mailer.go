package mailer

import (
	"context"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type Mailer struct {
	mg      mailgun.Mailgun
	Timeout time.Duration
}

func New() *Mailer {
	mg := mailgun.NewMailgun(
		os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_API_KEY"),
	)
	return &Mailer{
		mg:      mg,
		Timeout: 5 * time.Second,
	}
}

func (m *Mailer) SendEmail(to, subject, htmlBody string) (string, error) {
	message := m.mg.NewMessage(os.Getenv("MAIL_FROM"), subject, "", to)
	message.SetHtml(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()

	resp, id, err := m.mg.Send(ctx, message)
	if err != nil {
		return "", err
	}
	return "ID: " + id + ", Response: " + resp, nil
}
