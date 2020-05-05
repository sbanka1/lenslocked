package email

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/mailgun/mailgun-go"
)

const (
	welcomeSubject = "Welome to Lenslocked!"
	resetSubject   = "Instructions for resetting your password"
	resetBaseURL   = "https://lenslocked.sbanka.io/reset"

	welcomeText = `Hi there!
	
	Welcome to Lenslocked! We really hope you enjoy using our application!
	
	Best,
	Saurabh
	`

	welcomeHTML = `Hi there!<br/>
	<br/>
	Welcome to <a href="https://lenslocked.sbanka.io">Lenslocked</a>! We really hope you enjoy using our application!<br/>
	<br/>	
	Best,<br/>
	Saurabh
	`

	resetTextImpl = `Hi there!

	It appears that you have requested a password reset. If this was you, please follow the link below to update your password:

	%s

	If you are asked for a token, please use the following value:

	%s

	If you didn't request a password reset you can safely ignore this email and your account will not be changed.

	Best,
	Lenslocked Support
	`

	resetHTMLImpl = `Hi there!<br/>
	<br/>
	It appears that you have requested a password reset. If this was you, please follow the link below to update your password:<br/>
	<br/>
	<a href="%s">%s</a><br/>
	<br/>
	If you are asked for a token, please use the following value:<br/>
	<br/>
	%s<br/>
	<br/>
	If you didn't request a password reset you can safely ignore this email and your account will not be changed.<br/>
	<br/>
	Best,<br/>
	Lenslocked Support<br/>
	`
)

func WithMailgun(domain, apiKey string) ClientConfig {
	return func(c *Client) {
		c.mg = mailgun.NewMailgun(domain, apiKey)
	}
}

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		// Set a default from email address...
		from: buildEmail("Saurabh Banka", "sbanka@lenslocked.sbanka.io"),
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func (c *Client) Welcome(toName, toEmail string) error {
	message := c.mg.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(toName, toEmail))
	message.SetHtml(welcomeHTML)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, _, err := c.mg.Send(ctx, message)
	return err
}

func (c *Client) ResetPw(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetURL := resetBaseURL + "?" + v.Encode()
	resetText := fmt.Sprintf(resetTextImpl, resetURL, token)
	message := c.mg.NewMessage(c.from, resetSubject, resetText, toEmail)
	resetHTML := fmt.Sprintf(resetHTMLImpl, resetURL, resetURL, token)
	message.SetHtml(resetHTML)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, _, err := c.mg.Send(ctx, message)
	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
