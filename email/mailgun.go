package email

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go"
)

const (
	welcomeSubject = "Welome to Lenslocked!"

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

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
