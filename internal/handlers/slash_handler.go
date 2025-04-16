package handlers

import (
    "fmt"
    "net/http"

    "github.com/rs/zerolog/log"
    "github.com/yourusername/slack-ai-bot/internal/config"
    "github.com/yourusername/slack-ai-bot/internal/conversation"
    "github.com/yourusername/slack-ai-bot/internal/openai"
    "github.com/yourusername/slack-ai-bot/internal/slackclient"
    "github.com/yourusername/slack-ai-bot/internal/utils"
)

func SlashCommandHandler(cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Verify Slack signature
        if !utils.VerifySlackSignature(r, cfg.SlackSigningSecret) {
            log.Warn().Msg("Invalid Slack signature for /slack/slash")
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Slack slash command sends form data
        if err := r.ParseForm(); err != nil {
            log.Error().Err(err).Msg("Failed to parse slash command form")
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }

        command := r.FormValue("command") // "/askgpt"
        text := r.FormValue("text")       // the user's message
        userID := r.FormValue("user_id")
        channelID := r.FormValue("channel_id")

        log.Info().Msgf("Slash command: %s from user=%s, text=%s, channel=%s", command, userID, text, channelID)

        if command != "/askgpt" {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("Unknown slash command"))
            return
        }

        // store user message in conversation
        conversation.AddMessage(userID, "user", text)

        aiResp, err := openai.GetAIResponseWithContext(userID, text)
        if err != nil {
            // ephemeral error reply
            ephemeralErr := slackclient.PostEphemeral(channelID, userID, fmt.Sprintf("OpenAI error: %v", err))
            if ephemeralErr != nil {
                log.Error().Err(ephemeralErr).Msg("Failed to post ephemeral error")
            }
            // Return 200 to Slack
            w.WriteHeader(http.StatusOK)
            return
        }

        // Post ephemeral message
        err = slackclient.PostEphemeral(channelID, userID, aiResp)
        if err != nil {
            log.Error().Err(err).Msg("Failed to post ephemeral slash response")
        }

        // Slack requires an immediate 200, no body needed (or you can supply a text)
        w.WriteHeader(http.StatusOK)
    }
}
