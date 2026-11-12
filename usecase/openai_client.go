package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const openAIURL = "https://api.openai.com/v1/chat/completions"

// errOpenAIRefused means the model declined to answer.
var errOpenAIRefused = errors.New("openai: model refused")

// errOpenAICutOff means the answer hit the max output tokens setting.
var errOpenAICutOff = errors.New("openai: answer was cut off (max output too low)")

type openAIMessage struct {
	Role    string `json:"role"` // system, user or assistant
	Content string `json:"content"`
}

// openAIClient is a minimal Chat Completions client shared by Tanya Loka and
// the trip planner. Key, model and output limit are read from settings on
// every request, so admin changes apply without a restart.
type openAIClient struct {
	settings AISettingsSource
	http     *http.Client
}

func newOpenAIClient(settings AISettingsSource) *openAIClient {
	return &openAIClient{settings: settings, http: &http.Client{Timeout: 80 * time.Second}}
}

// enabled reports whether an API key is configured.
func (c *openAIClient) enabled() bool {
	return c != nil && c.settings != nil && c.settings.Current().APIKey != ""
}

// chat sends one request and returns the text of the first choice. When
// schema is set the answer is forced to match that JSON schema.
func (c *openAIClient) chat(ctx context.Context, tag string, messages []openAIMessage, temperature float64, schemaName string, schema map[string]any) (string, error) {
	s := c.settings.Current()
	payload := map[string]any{
		"model":                 s.Model,
		"temperature":           temperature,
		"messages":              messages,
		"max_completion_tokens": s.MaxOutputTokens,
	}
	if schema != nil {
		payload["response_format"] = map[string]any{
			"type": "json_schema",
			"json_schema": map[string]any{
				"name":   schemaName,
				"strict": true,
				"schema": schema,
			},
		}
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai status %d: %s", res.StatusCode, truncateRunes(string(data), 300))
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
				Refusal string `json:"refusal"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens        int `json:"prompt_tokens"`
			CompletionTokens    int `json:"completion_tokens"`
			PromptTokensDetails struct {
				CachedTokens int `json:"cached_tokens"`
			} `json:"prompt_tokens_details"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil || len(parsed.Choices) == 0 {
		return "", errors.New("openai: empty response")
	}
	log.Printf("[%s] openai tokens in=%d cached=%d out=%d", tag,
		parsed.Usage.PromptTokens, parsed.Usage.PromptTokensDetails.CachedTokens, parsed.Usage.CompletionTokens)

	choice := parsed.Choices[0]
	if choice.Message.Refusal != "" {
		return "", errOpenAIRefused
	}
	if choice.FinishReason == "length" {
		return "", errOpenAICutOff
	}
	return choice.Message.Content, nil
}
