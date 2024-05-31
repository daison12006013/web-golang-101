package mailer

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
)

type MockMailgun struct {
	mailgun.MailgunImpl
	SendFunc func(ctx context.Context, m *mailgun.Message) (string, string, error)
}

func (m *MockMailgun) Send(ctx context.Context, message *mailgun.Message) (string, string, error) {
	return m.SendFunc(ctx, message)
}

func TestNew(t *testing.T) {
	// Set environment variables
	os.Setenv("MAILGUN_DOMAIN", "example.com")
	os.Setenv("MAILGUN_API_KEY", "testkey")

	mailer := New()

	// Type assertion to access the underlying type of the Mailgun interface
	mgImpl, ok := mailer.mg.(*mailgun.MailgunImpl)
	if !ok {
		t.Fatal("Expected MailgunImpl")
	}

	assert.Equal(t, "example.com", mgImpl.Domain())
	assert.Equal(t, "testkey", mgImpl.APIKey())
	assert.Equal(t, 5*time.Second, mailer.Timeout)
}

func TestSendEmail(t *testing.T) {
	mockMg := &MockMailgun{
		MailgunImpl: *mailgun.NewMailgun("example.com", "testkey"),
		SendFunc: func(ctx context.Context, m *mailgun.Message) (string, string, error) {
			return "OK", "TestID", nil
		},
	}

	mailer := &Mailer{
		mg:      mockMg,
		Timeout: 5 * time.Second,
	}

	to := "test@example.com"
	subject := "Test Subject"
	htmlBody := "<h1>Test Body</h1>"

	resp, err := mailer.SendEmail(to, subject, htmlBody)

	assert.NoError(t, err)
	assert.Equal(t, "ID: TestID, Response: OK", resp)
}
