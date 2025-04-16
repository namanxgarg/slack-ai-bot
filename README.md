Slack AI Bot
A production-style Slack bot in Go that integrates with OpenAI’s ChatGPT to respond to messages and slash commands with AI-generated text. Includes:

Conversation Context: Maintains short-term chat history for each user to provide more relevant AI responses.

Slash Commands (/askgpt): Ephemeral replies to the user’s command for a private AI conversation.

Ephemeral Messages: Ensures slash replies are only visible to the requesting user.

Kubernetes Manifests: For production deployment.

Docker / Docker Compose: For local development.

Zerolog: Structured logging for better observability.

Features
Responds to Channel Messages: Users can mention or talk in a channel; the bot replies with an AI response.

Slash Command (/askgpt) for ephemeral Q&A.

In-Memory Conversation: Short conversation log per user, so the AI can provide some continuity.

Slack Events API: Uses Slack’s signing secret to verify requests.

OpenAI Integration: Calls GPT-3.5 for advanced, natural language replies.

Docker + Kubernetes for easy container-based deployment.

Project Structure
bash
Copy
Edit
slack-ai-bot/
├── cmd/
│   └── main.go                    # Entry point
├── internal/
│   ├── config/
│   │   └── config.go             # Loads env vars
│   ├── conversation/
│   │   └── context.go            # In-memory conversation store
│   ├── handlers/
│   │   ├── events_handler.go     # Slack Events API
│   │   └── slash_handler.go      # /askgpt slash command
│   ├── openai/
│   │   └── openai_client.go      # Calls OpenAI's chat/completions
│   ├── slackclient/
│   │   └── slack_client.go       # Thin wrapper around slack-go
│   └── utils/
│       └── verify_signature.go   # Slack signature verification
├── k8s/
│   ├── slack-bot-deployment.yaml # K8s Deployment & Secret example
│   └── slack-bot-service.yaml    # K8s Service (LoadBalancer)
├── .env                          # Local dev environment variables
├── Dockerfile
├── docker-compose.yml            # Optional for local dev
├── go.mod
└── go.sum
Requirements
Go 1.20 or higher

Slack Bot with Bot Token and Signing Secret

OpenAI API Key (ChatGPT, GPT-3.5, etc.)

(Optional) Docker and/or Kubernetes to run in containers

Installation
Clone this repo:

bash
Copy
Edit
git clone https://github.com/yourusername/slack-ai-bot.git
cd slack-ai-bot
Set up environment variables:

SLACK_SIGNING_SECRET: Your Slack App’s signing secret.

SLACK_BOT_TOKEN: Your Slack Bot User OAuth token (starts with xoxb-).

OPENAI_API_KEY: Your OpenAI key (sk-...).

PORT: (optional) default 3000.

For local development, create a .env file:

bash
Copy
Edit
SLACK_SIGNING_SECRET=xxx
SLACK_BOT_TOKEN=xoxb-xxx
OPENAI_API_KEY=sk-xxx
PORT=3000
Local build & run without Docker:

bash
Copy
Edit
go mod tidy
go run ./cmd/main.go
The bot listens on port 3000 by default.

Set Slack App Event Subscription:

In Slack’s App config, go to Event Subscriptions → Enable.

Request URL: https://YOUR_HOST/slack/events.

Slack will verify by sending a challenge request.

Subscribe to Bot Events → add message.channels (and others if needed).

Add Slash Command /askgpt with URL: https://YOUR_HOST/slack/slash.

Invite your bot to relevant channels (e.g., #general).

Usage
Channel Messages: When a user posts a message, the bot calls OpenAI with stored context (last 5 messages from that user) and replies publicly.

Slash Command: Type /askgpt [your question]. The bot replies ephemerally (visible only to you).

Docker Compose (Local Dev)
Edit .env with your Slack & OpenAI keys.

Build & run:

bash
Copy
Edit
docker-compose build
docker-compose up
The bot is on http://localhost:3000. Expose that port to the public if Slack can’t reach localhost.

Kubernetes Deployment
Build and push your Docker image to a registry:

bash
Copy
Edit
docker build -t your-registry/slack-ai-bot:latest .
docker push your-registry/slack-ai-bot:latest
Create a secret in K8s with your Slack & OpenAI credentials:

yaml
Copy
Edit
# example snippet in k8s/slack-bot-deployment.yaml
apiVersion: v1
kind: Secret
metadata:
  name: slack-bot-secrets
type: Opaque
data:
  signingSecret: BASE64_ENCODED_SLACK_SECRET
  botToken: BASE64_ENCODED_BOT_TOKEN
  openaiKey: BASE64_ENCODED_OPENAI_KEY
Apply the manifests:

bash
Copy
Edit
kubectl apply -f k8s/slack-bot-deployment.yaml
kubectl apply -f k8s/slack-bot-service.yaml
Expose the service (LoadBalancer or Ingress) so Slack can reach your /slack/events & /slack/slash endpoints.

Advanced Notes
Conversation Context: Currently in memory. For production or horizontal scaling, store in a DB or Redis.

Ephemeral vs Public:

Normal channel messages → public reply.

Slash commands → ephemeral reply only to the user.

Logging: Uses zerolog for structured logs. You can add distributed tracing or advanced log configs.

Scaling: If you run multiple replicas, ensure shared state for conversation logs, or a consistent partitioning scheme per user.

