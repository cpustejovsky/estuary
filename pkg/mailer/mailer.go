package mailer

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func formatHtml(token string, m *mailgun.Message) {
	//TODO: make use of html/templates for templating
	var message bytes.Buffer
	message.WriteString("<h1>Password Reset</h1>")
	message.WriteString("<p>Please visit the following link to reset your password</p>")
	link := fmt.Sprintf("'http://localhost:3000/new-password?token=%v'", token)
	message.WriteString("<a target='_blank' rel='noopener noreferrer' href=" + link + ">Reset Your Password</a>")

	m.SetHtml(message.String())
}

func SendEmail(recipient, token string, mg *mailgun.MailgunImpl) error {
	sender := "password-reset@estuaryapp.com"
	subject := "Password Reset"
	html := ""

	m := mg.NewMessage(sender, subject, html, recipient)

	formatHtml(token, m)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, m)

	if err != nil {
		return err
	}

	fmt.Printf("MailGun API:\nID: %s\nResp: %s\n", id, resp)
	return nil
}
