package conversation

import (
    "sync"
    "time"
)

// We'll store a short conversation log for each user
type Message struct {
    Role    string // "user" or "assistant"
    Content string
    Time    time.Time
}

var (
    convStore map[string][]Message
    mu        sync.Mutex
)

func InitConversationStore() {
    convStore = make(map[string][]Message)
}

// Add message to user's conversation
func AddMessage(userID, role, text string) {
    mu.Lock()
    defer mu.Unlock()
    convStore[userID] = append(convStore[userID], Message{
        Role:    role,
        Content: text,
        Time:    time.Now(),
    })

    // For demonstration, keep only last 5 messages
    if len(convStore[userID]) > 5 {
        convStore[userID] = convStore[userID][len(convStore[userID])-5:]
    }
}

// Get conversation for user
func GetConversation(userID string) []Message {
    mu.Lock()
    defer mu.Unlock()
    return convStore[userID]
}
