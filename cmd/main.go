package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/yourusername/slack-ai-bot/internal/config"
	"github.com/yourusername/slack-ai-bot/internal/conversation"
	"github.com/yourusername/slack-ai-bot/internal/handlers"
	"github.com/yourusername/slack-ai-bot/internal/openai"
	"github.com/yourusername/slack-ai-bot/internal/slackclient"
)

func main() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file, reading from environment variables.")
	}

	// Setup zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zlog.With().Str("service", "slack-ai-bot").Logger()

	// Load config
	cfg := config.LoadConfigFromEnv()

	// Initialize Slack & OpenAI
	slackclient.InitSlackClient(cfg.SlackBotToken)
	openai.InitOpenAI(cfg.OpenAIApiKey)

	// Initialize in-memory conversation store
	conversation.InitConversationStore()

	r := mux.NewRouter()

	// Slack Events Endpoint
	r.HandleFunc("/slack/events", handlers.SlackEventsHandler(cfg)).Methods("POST")

	// Slash Commands Endpoint (for /askgpt)
	r.HandleFunc("/slack/slash", handlers.SlashCommandHandler(cfg)).Methods("POST")

	// Health Check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	logger.Info().Msgf("Slack bot listening on port %s", port)
	http.ListenAndServe(":"+port, r)
}
