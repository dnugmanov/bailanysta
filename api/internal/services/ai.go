package services

import (
	"context"
	"fmt"
	"strings"

	"bailanysta/api/internal/pkg/ai"
)

type AIService struct {
	client *ai.Client
}

type GenerateTextRequest struct {
	Prompt      string  `json:"prompt" validate:"required,min=1,max=2000"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	Context     string  `json:"context,omitempty"` // Additional context for generation
}

type GenerateTextResponse struct {
	Text  string `json:"text"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
}

type GeneratePostRequest struct {
	Topic     string `json:"topic" validate:"required,min=3,max=200"`
	Course    string `json:"course,omitempty"`
	Module    string `json:"module,omitempty"`
	Style     string `json:"style,omitempty"` // academic, casual, professional
	MaxTokens int    `json:"max_tokens,omitempty"`
}

type GenerateCommentRequest struct {
	PostContent string `json:"post_content" validate:"required"`
	Style       string `json:"style,omitempty"`
	MaxTokens   int    `json:"max_tokens,omitempty"`
}

func NewAIService(client *ai.Client) *AIService {
	return &AIService{client: client}
}

func (s *AIService) GenerateText(ctx context.Context, req GenerateTextRequest) (*GenerateTextResponse, error) {
	// Enhance prompt with context if provided
	prompt := req.Prompt
	if req.Context != "" {
		prompt = fmt.Sprintf("Context: %s\n\nRequest: %s", req.Context, req.Prompt)
	}

	// Set default values
	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 500
	}
	if maxTokens > 4000 {
		maxTokens = 4000
	}

	temperature := req.Temperature
	if temperature <= 0 {
		temperature = 0.7
	}
	if temperature > 2.0 {
		temperature = 2.0
	}

	text, err := s.client.GenerateText(ctx, prompt, maxTokens, temperature)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	response := &GenerateTextResponse{
		Text:  strings.TrimSpace(text),
		Model: "openai/gpt-oss-120b",
	}

	return response, nil
}

func (s *AIService) GeneratePost(ctx context.Context, req GeneratePostRequest) (*GenerateTextResponse, error) {
	// Build enhanced prompt for post generation
	var promptBuilder strings.Builder

	promptBuilder.WriteString("Create an educational post about: ")
	promptBuilder.WriteString(req.Topic)
	promptBuilder.WriteString("\n\n")

	if req.Course != "" {
		promptBuilder.WriteString(fmt.Sprintf("Course: %s\n", req.Course))
	}
	if req.Module != "" {
		promptBuilder.WriteString(fmt.Sprintf("Module: %s\n", req.Module))
	}

	// Style instructions
	style := req.Style
	if style == "" {
		style = "academic"
	}

	switch style {
	case "academic":
		promptBuilder.WriteString("Style: Write in an academic tone suitable for students. Include key concepts, explanations, and be informative.\n")
	case "casual":
		promptBuilder.WriteString("Style: Write in a casual, conversational tone that's easy to understand.\n")
	case "professional":
		promptBuilder.WriteString("Style: Write in a professional tone with practical insights.\n")
	default:
		promptBuilder.WriteString("Style: Write in an informative and engaging manner.\n")
	}

	promptBuilder.WriteString("\nMake the post comprehensive but concise. Include relevant hashtags at the end.")
	promptBuilder.WriteString("\n\nWrite the post content:")

	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 800
	}

	text, err := s.client.GenerateText(ctx, promptBuilder.String(), maxTokens, 0.7)
	if err != nil {
		return nil, fmt.Errorf("failed to generate post: %w", err)
	}

	response := &GenerateTextResponse{
		Text:  strings.TrimSpace(text),
		Model: "openai/gpt-oss-120b",
	}

	return response, nil
}

func (s *AIService) GenerateComment(ctx context.Context, req GenerateCommentRequest) (*GenerateTextResponse, error) {
	// Build prompt for comment generation
	var promptBuilder strings.Builder

	promptBuilder.WriteString("Write a thoughtful comment for this post:\n\n")
	promptBuilder.WriteString(fmt.Sprintf("Post content: %s\n\n", req.PostContent))

	// Style instructions
	style := req.Style
	if style == "" {
		style = "constructive"
	}

	switch style {
	case "constructive":
		promptBuilder.WriteString("Style: Write a constructive comment that adds value, asks questions, or shares related insights.\n")
	case "appreciative":
		promptBuilder.WriteString("Style: Write an appreciative comment that shows engagement and positive feedback.\n")
	case "questioning":
		promptBuilder.WriteString("Style: Write a comment that asks relevant questions to deepen the discussion.\n")
	default:
		promptBuilder.WriteString("Style: Write an engaging and relevant comment.\n")
	}

	promptBuilder.WriteString("\nKeep the comment concise and natural. Write the comment:")

	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 200
	}

	text, err := s.client.GenerateText(ctx, promptBuilder.String(), maxTokens, 0.8)
	if err != nil {
		return nil, fmt.Errorf("failed to generate comment: %w", err)
	}

	response := &GenerateTextResponse{
		Text:  strings.TrimSpace(text),
		Model: "openai/gpt-oss-120b",
	}

	return response, nil
}

func (s *AIService) ValidateConnection(ctx context.Context) error {
	return s.client.ValidateConnection(ctx)
}

// Helper methods for specific use cases

func (s *AIService) GenerateStudyNotes(ctx context.Context, topic, course string) (*GenerateTextResponse, error) {
	prompt := fmt.Sprintf("Create comprehensive study notes about '%s'", topic)
	if course != "" {
		prompt += fmt.Sprintf(" for the course '%s'", course)
	}
	prompt += ". Include key concepts, definitions, important points to remember, examples, and detailed explanations. Format as markdown with headers and lists."

	return s.GenerateText(ctx, GenerateTextRequest{
		Prompt:      prompt,
		MaxTokens:   2500,
		Temperature: 0.3, // Lower temperature for more factual content
	})
}

func (s *AIService) GenerateQuiz(ctx context.Context, topic, course string) (*GenerateTextResponse, error) {
	prompt := fmt.Sprintf("Create a 5-question quiz about '%s'", topic)
	if course != "" {
		prompt += fmt.Sprintf(" from the course '%s'", course)
	}
	prompt += ". Include multiple choice questions with 4 options each and indicate the correct answers."

	return s.GenerateText(ctx, GenerateTextRequest{
		Prompt:      prompt,
		MaxTokens:   600,
		Temperature: 0.5,
	})
}

func (s *AIService) ExplainConcept(ctx context.Context, concept, context string) (*GenerateTextResponse, error) {
	prompt := fmt.Sprintf("Provide a detailed explanation of the concept '%s'", concept)
	if context != "" {
		prompt += fmt.Sprintf(" in the context of %s", context)
	}
	prompt += ". Include: 1) Clear definition, 2) Key characteristics, 3) Practical examples, 4) How it works, 5) Why it's important. Use simple language but be comprehensive. Format as markdown with headers."

	return s.GenerateText(ctx, GenerateTextRequest{
		Prompt:      prompt,
		MaxTokens:   2000,
		Temperature: 0.6,
	})
}
