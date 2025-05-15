# Slack AI Bot

A production-style Slack bot in Go that integrates with OpenAI’s ChatGPT to respond to messages and slash commands with AI-generated text. Includes:

- **Conversation Context**: Maintains short-term chat history for each user to provide more relevant AI responses.
- **Slash Commands** (`/askgpt`): Ephemeral replies to the user’s command for a private AI conversation.
- **Ephemeral Messages**: Ensures slash replies are only visible to the requesting user.
- **Kubernetes Manifests**: For production deployment.
- **Docker** / **Docker Compose**: For local development.
- **Zerolog**: Structured logging for better observability.

---

## Features

1. **Responds to Channel Messages**: Users can mention or talk in a channel; the bot replies with an AI response.  
2. **Slash Command** (`/askgpt`) for ephemeral Q&A.  
3. **In-Memory Conversation**: Short conversation log per user, so the AI can provide some continuity.  
4. **Slack Events API**: Uses Slack’s signing secret to verify requests.  
5. **OpenAI Integration**: Calls GPT-3.5 for advanced, natural language replies.  
6. **Docker + Kubernetes** for easy container-based deployment.
