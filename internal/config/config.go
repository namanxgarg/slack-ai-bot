package config

import "os"

type Config struct {
    SlackSigningSecret string
    SlackBotToken      string
    OpenAIApiKey       string
    Port               string
}

func LoadConfigFromEnv() *Config {
    return &Config{
        SlackSigningSecret: os.Getenv("SLACK_SIGNING_SECRET"),
        SlackBotToken:      os.Getenv("SLACK_BOT_TOKEN"),
        OpenAIApiKey:       os.Getenv("OPENAI_API_KEY"),
        Port:               os.Getenv("PORT"),
    }
}
