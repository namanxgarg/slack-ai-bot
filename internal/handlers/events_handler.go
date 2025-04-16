package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/rs/zerolog/log"
    "github.com/yourusername/slack-ai-bot/internal/config"
    "github.com/yourusername/slack-ai-bot/internal/conversation"
    "github.com/yourusername/slack-ai-bot/internal/openai"
    "github.com/yourusername/slack-ai-bot/internal/slackclient"
    "github.com/yourusername/slack-ai-bot/internal/utils"
)

type SlackEvent struct {
    Token     string `json:"token"`
    Challenge string `json:"challenge"`
    Type      string `json:"type"`
    Event     struct {
        Type    string `json:"type"`
        Text    string `json:"text"`
        User    string `json:"user"`
        Channel string `json:"channel"`
        BotID   string `json:"bot_id,omitempty"`
    } `json:"event"`
}

// SlackEventsHandler verifies signature, handles url_verification, message events
func SlackEventsHandler(cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Verify Slack signature
        if !utils.VerifySlackSignature(r, cfg.SlackSigningSecret) {
            log.Warn().Msg("Invalid Slack signature for /slack/events")
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        var se SlackEvent
        if err := json.NewDecoder(r.Body).Decode(&se); err != nil {
            log.Error().Err(err).Msg("Failed to decode SlackEvent")
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }

        // url_verification
        if se.Type == "url_verification" {
            w.Write([]byte(se.Challenge))
            return
        }

        // If it's a standard message event
        if se.Event.Type == "message" && se.Event.BotID == "" && se.Event.User != "" {
            // store user message in conversation
            conversation.AddMessage(se.Event.User, "user", se.Event.Text)

            go func(user, text, channel string) {
                aiResp, err := openai.GetAIResponseWithContext(user, text)
                if err != nil {
                    log.Error().Err(err).Msg("OpenAI error in /slack/events")
                    return
                }
                // Post publicly in channel
                if err := slackclient.PostMessage(channel, aiResp); err != nil {
                    log.Error().Err(err).Msg("Failed to post Slack message from /slack/events")
                }
            }(se.Event.User, se.Event.Text, se.Event.Channel)
        }

        w.WriteHeader(http.StatusOK)
    }
}
