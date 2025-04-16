package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/slack-ai-bot/internal/conversation"
)

var openAIKey string

func InitOpenAI(apiKey string) {
	openAIKey = apiKey
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float32       `json:"temperature,omitempty"`
}
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func GetAIResponseWithContext(userID, newUserMessage string) (string, error) {
	if openAIKey == "" {
		return "", fmt.Errorf("missing OPENAI_API_KEY")
	}

	// 1) Build the conversation messages
	var messages []ChatMessage
	// Provide a system prompt at the start
	messages = append(messages, ChatMessage{
		Role:    "system",
		Content: "You are a helpful AI assistant for Slack users.",
	})

	// Add prior conversation
	conv := conversation.GetConversation(userID)
	for _, msg := range conv {
		role := "user"
		if msg.Role == "assistant" {
			role = "assistant"
		}
		messages = append(messages, ChatMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Add this new user message
	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: newUserMessage,
	})

	// 2) Call OpenAI
	reqBody := ChatRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		MaxTokens:   200,
		Temperature: 0.7,
	}
	data, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAIKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI returned status %d", resp.StatusCode)
	}

	var cr ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return "", err
	}
	if len(cr.Choices) == 0 {
		return "", fmt.Errorf("no choices from OpenAI")
	}
	aiReply := cr.Choices[0].Message.Content

	// 3) Update conversation store
	conversation.AddMessage(userID, "assistant", aiReply)

	return aiReply, nil
}
