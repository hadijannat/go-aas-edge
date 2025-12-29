package ai

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// Agent handles AI-powered diagnostics using RAG (Retrieval Augmented Generation).
type Agent struct {
	client *openai.Client
}

// NewAgent creates a new AI agent. Returns nil client if API key is not set.
func NewAgent() *Agent {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return &Agent{client: nil}
	}
	return &Agent{client: openai.NewClient(key)}
}

// IsConfigured returns true if the AI agent has a valid OpenAI client.
func (a *Agent) IsConfigured() bool {
	return a.client != nil
}

// Diagnose uses RAG to answer questions about the asset based on live telemetry.
func (a *Agent) Diagnose(ctx context.Context, telemetry map[string]string, userQuestion string) (string, error) {
	if a.client == nil {
		// Fallback mock response when AI is not configured
		return fmt.Sprintf(
			"[AI MOCK]: My temperature is %s°C and RPM is %s. I am operating within normal parameters.",
			telemetry["Temperature"], telemetry["RPM"],
		), nil
	}

	// Dynamic Prompt Construction (The RAG Step)
	systemPrompt := fmt.Sprintf(`You are the embedded intelligence of a High-Speed Industrial Motor.
Current Telemetry:
- Temperature: %s°C (Max Safe: 85°C)
- RPM: %s (Max Safe: 3000)

Task: Answer the user's question based strictly on the current telemetry.
Be concise and speak as the machine. If values approach limits, warn the operator.
`, telemetry["Temperature"], telemetry["RPM"])

	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
				{Role: openai.ChatMessageRoleUser, Content: userQuestion},
			},
			MaxTokens:   150,
			Temperature: 0.7,
		},
	)

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
