package mailer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/mailgun/mailgun-go/v5"
	"go.uber.org/zap"
)

type MailGunMailer struct {
	domain    string
	fromEmail string
	apiKey    string
	client    *mailgun.Client
	logger    *zap.SugaredLogger
}

func NewMailGunMailer(logger *zap.SugaredLogger, domain string, fromEmail string, apiKey string) *MailGunMailer {
	if apiKey == "" {
		return nil
	}

	client := mailgun.NewMailgun(apiKey)
	mailer := MailGunMailer{
		domain:    domain,
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
		logger:    logger,
	}
	return &mailer
}

func (mailer *MailGunMailer) Send(templateFile string, username string, email string, data any, isSandbox bool) error {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	msg := mailgun.NewMessage(mailer.domain, mailer.fromEmail, subject.String(), "", email)
	msg.SetHTML(body.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10-second timeout
	for i := 0; i < maxRetries; i++ {
		resp, err := mailer.client.Send(ctx, msg)
		if err != nil {
			mailer.logger.Warnw(fmt.Sprintf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetries), "error", err.Error())

			// exponentional backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		mailer.logger.Infow("Sent email", "ID", resp.ID, "Resp:", resp.Message)

		return nil
	}

	return fmt.Errorf("Failed to send email after %d attempts", maxRetries)
}
