module github.com/yourusername/slack-ai-bot

go 1.20

require (
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.4.0
	github.com/rs/zerolog v1.29.0
	github.com/slack-go/slack v0.10.1
// If you want to do a direct REST call to OpenAI, no additional package needed.
// Or you can add github.com/sashabaranov/go-openai if you prefer.
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.6.0 // indirect
)
