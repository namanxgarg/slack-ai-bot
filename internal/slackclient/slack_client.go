package slackclient

import (
	"github.com/slack-go/slack"
)

var api *slack.Client

func InitSlackClient(botToken string) {
	api = slack.New(botToken)
}

func PostMessage(channel string, text string) error {
	_, _, err := api.PostMessage(channel, slack.MsgOptionText(text, false))
	return err
}

// Ephemeral reply to a slash command or any ephemeral message
func PostEphemeral(channel, user, text string) error {
	_, err := api.PostEphemeral(channel, user, slack.MsgOptionText(text, false))
	return err
}

// (Optional) If you need the underlying client
func GetClient() *slack.Client {
	return api
}
