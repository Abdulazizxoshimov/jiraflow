package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var issueKeyRe = regexp.MustCompile(`\b([A-Z][A-Z0-9]+-\d+)\b`)

func VerifyWebhookSignature(secret, body []byte, signature string) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func ExtractIssueKeys(text string) []string {
	matches := issueKeyRe.FindAllString(text, -1)
	seen := make(map[string]bool)
	var keys []string
	for _, m := range matches {
		if !seen[m] {
			seen[m] = true
			keys = append(keys, m)
		}
	}
	return keys
}

type Commit struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	URL       string `json:"url"`
	Timestamp string `json:"timestamp"`
	Author    struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
}

type PushEvent struct {
	Ref        string   `json:"ref"`
	Repository struct {
		FullName string `json:"full_name"`
		HTMLURL  string `json:"html_url"`
	} `json:"repository"`
	Commits []Commit `json:"commits"`
}

type PREvent struct {
	Action string `json:"action"`
	Number int    `json:"number"`
	PullRequest struct {
		Title   string `json:"title"`
		Body    string `json:"body"`
		State   string `json:"state"`
		HTMLURL string `json:"html_url"`
		MergedAt *string `json:"merged_at"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"pull_request"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func ParsePushEvent(body []byte) (*PushEvent, error) {
	var ev PushEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil, fmt.Errorf("github.ParsePushEvent: %w", err)
	}
	return &ev, nil
}

func ParsePREvent(body []byte) (*PREvent, error) {
	var ev PREvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil, fmt.Errorf("github.ParsePREvent: %w", err)
	}
	return &ev, nil
}

var closeKeywordsRe = regexp.MustCompile(`(?i)(?:fix(?:es|ed)?|close[sd]?|resolve[sd]?)\s+([A-Z][A-Z0-9]+-\d+)`)

func ExtractClosingKeys(text string) []string {
	matches := closeKeywordsRe.FindAllStringSubmatch(text, -1)
	seen := make(map[string]bool)
	var keys []string
	for _, m := range matches {
		key := strings.ToUpper(m[1])
		if !seen[key] {
			seen[key] = true
			keys = append(keys, key)
		}
	}
	return keys
}
