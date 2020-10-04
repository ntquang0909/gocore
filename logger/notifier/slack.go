package notifier

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

func redirectPolicyFunc(req gorequest.Request, via []gorequest.Request) error {
	return fmt.Errorf("Incorrect token (redirection)")
}

// SlackNotifier instance
type SlackNotifier struct {
	WebhookURL string
	ProxyURL   string
	Channel    string
}

// New instance
func New(webhookURL, proxyURL, channel string) *SlackNotifier {
	return &SlackNotifier{
		WebhookURL: webhookURL,
		ProxyURL:   proxyURL,
		Channel:    channel,
	}
}

// Send send msg
func (slack *SlackNotifier) Send(title, body string) []error {
	var payload = map[string]interface{}{
		"channel": slack.Channel,
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]interface{}{

					"type": "mrkdwn",
					"text": fmt.Sprintf("```%s\n%s```", title, body),
				},
			},
		},
	}
	request := gorequest.New().Proxy(slack.ProxyURL)
	resp, _, err := request.
		Post(slack.WebhookURL).
		RedirectPolicy(redirectPolicyFunc).
		Send(payload).
		End()

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return []error{fmt.Errorf("Error sending msg. Status: %v", resp.Status)}
	}

	return nil
}
