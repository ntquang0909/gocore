package twilio

import (
	"net/url"

	"github.com/kevinburke/rest"
	"github.com/kevinburke/twilio-go"
)

// Twilio instance
type Twilio struct {
	*twilio.Client
	config *Config
}

// Config config
type Config struct {
	AccountSID  string
	AuthToken   string
	SenderPhone string
}

// New init
func New(config *Config) *Twilio {
	return &Twilio{
		Client: twilio.NewClient(config.AccountSID, config.AuthToken, nil),
		config: config,
	}
}

// SendSMS send otp
func (client *Twilio) SendSMS(to string, message string, mediaURLs ...url.URL) error {
	var urls []*url.URL
	for _, url := range mediaURLs {
		urls = append(urls, &url)
	}

	_, err := client.Messages.SendMessage(client.config.SenderPhone, to, message, urls)
	if err != nil {
		switch v := err.(type) {
		case *rest.Error:
			return &Error{
				Detail:     v.Detail,
				Message:    v.Title,
				StatusCode: v.Status,
			}
		}

		return err
	}

	return nil

}
