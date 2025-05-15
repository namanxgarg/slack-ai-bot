package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// VerifySlackRequest verifies that the request is coming from Slack
func VerifySlackRequest(r *http.Request) error {
	timestamp := r.Header.Get("X-Slack-Request-Timestamp")
	signature := r.Header.Get("X-Slack-Signature")

	if timestamp == "" || signature == "" {
		return fmt.Errorf("missing Slack headers")
	}

	// Verify timestamp is within 5 minutes
	ts, err := time.Parse("", timestamp)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %v", err)
	}

	if time.Since(ts) > 5*time.Minute {
		return fmt.Errorf("request too old")
	}

	// Get request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %v", err)
	}
	// Restore body for later use
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Create signature base string
	baseString := fmt.Sprintf("v0:%s:%s", timestamp, string(body))

	// Calculate signature
	mac := hmac.New(sha256.New, []byte(os.Getenv("SLACK_SIGNING_SECRET")))
	mac.Write([]byte(baseString))
	expectedSignature := "v0=" + hex.EncodeToString(mac.Sum(nil))

	// Compare signatures
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}
