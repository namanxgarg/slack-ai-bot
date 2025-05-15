package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the structured logger
func InitLogger() {
	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Use JSON format in production, pretty console output in development
	if os.Getenv("ENV") == "production" {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		log.Logger = zerolog.New(output).With().Timestamp().Logger()
	}

	// Add global fields
	log.Logger = log.With().
		Str("service", "slack-ai-bot").
		Str("version", "1.0.0").
		Logger()
}

// LogRequest logs an incoming HTTP request
func LogRequest(method, path, userID string, duration time.Duration, status int) {
	log.Info().
		Str("method", method).
		Str("path", path).
		Str("user_id", userID).
		Dur("duration", duration).
		Int("status", status).
		Msg("request completed")
}

// LogError logs an error with context
func LogError(err error, context map[string]interface{}) {
	event := log.Error().Err(err)
	for k, v := range context {
		event.Interface(k, v)
	}
	event.Msg("error occurred")
}

// LogSlackEvent logs a Slack event
func LogSlackEvent(eventType, userID string, metadata map[string]interface{}) {
	event := log.Info().
		Str("event_type", eventType).
		Str("user_id", userID)

	for k, v := range metadata {
		event.Interface(k, v)
	}
	event.Msg("slack event received")
}
