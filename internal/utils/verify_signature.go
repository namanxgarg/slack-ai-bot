package utils

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
    "strconv"
    "strings"
    "time"
)

func VerifySlackSignature(r *http.Request, signingSecret string) bool {
    if signingSecret == "" {
        return false
    }

    // Slack headers
    timestamp := r.Header.Get("X-Slack-Request-Timestamp")
    slackSig := r.Header.Get("X-Slack-Signature")

    // Check replay attack
    ts, err := strconv.ParseInt(timestamp, 10, 64)
    if err != nil {
        return false
    }
    if time.Now().Unix()-ts > 60*5 {
        // older than 5 minutes => potential replay
        return false
    }

    // Slack signature base string
    rawBody, err := io.ReadAll(r.Body)
    if err != nil {
        return false
    }
    // re-inject into request for downstream reading
    r.Body.Close()
    r.Body = io.NopCloser(strings.NewReader(string(rawBody)))

    baseString := fmt.Sprintf("v0:%s:%s", timestamp, rawBody)
    mac := hmac.New(sha256.New, []byte(signingSecret))
    mac.Write([]byte(baseString))
    computedSig := "v0=" + hex.EncodeToString(mac.Sum(nil))

    return hmac.Equal([]byte(computedSig), []byte(slackSig))
}
