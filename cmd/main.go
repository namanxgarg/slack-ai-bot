package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/yourusername/slack-ai-bot/internal/config"
	"github.com/yourusername/slack-ai-bot/internal/conversation"
	"github.com/yourusername/slack-ai-bot/internal/handlers"
	"github.com/yourusername/slack-ai-bot/internal/openai"
	"github.com/yourusername/slack-ai-bot/internal/slackclient"
	"github.com/yourusername/slack-ai-bot/internal/utils"
)

func main() {
	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Info().Msg("No .env file, reading from environment variables")
	}

	// Initialize structured logging
	utils.InitLogger()

	// Load config
	cfg := config.LoadConfigFromEnv()

	// Initialize Slack & OpenAI
	slackclient.InitSlackClient(cfg.SlackBotToken)
	openai.InitOpenAI(cfg.OpenAIApiKey)

	// Initialize conversation cache with 1-hour TTL
	conversation.InitCache(1 * time.Hour)

	r := mux.NewRouter()

	// Middleware for request logging
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			utils.LogRequest(r.Method, r.URL.Path, r.Header.Get("X-Slack-User-ID"), duration, 200)
		})
	})

	// Slack Events Endpoint with signature verification
	r.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.VerifySlackRequest(r); err != nil {
			log.Error().Err(err).Msg("Invalid Slack request signature")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handlers.SlackEventsHandler(cfg)(w, r)
	}).Methods("POST")

	// Slash Commands Endpoint with signature verification
	r.HandleFunc("/slack/slash", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.VerifySlackRequest(r); err != nil {
			log.Error().Err(err).Msg("Invalid Slack request signature")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handlers.SlashCommandHandler(cfg)(w, r)
	}).Methods("POST")

	// Health Check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	log.Info().Str("port", port).Msg("Slack bot listening")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
