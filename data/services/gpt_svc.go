package services

import (
	"context"
	"fmt"
	"go-fiber/bootstrap"
	"log"

	"github.com/sashabaranov/go-openai"
)

type GPTServices interface {
	GenerateReply(input string) string
}

type GPTServicesImpl struct{}

func (s *GPTServicesImpl) GenerateReply(input string) string {
	client := openai.NewClient(bootstrap.GlobalEnv.Keys.OpenaiApiKey)

	prompt := fmt.Sprintf(`You are a smart assistant. Product info:
	- SmartWatch X: Heart rate, Sleep tracker, GPS.
	- SmartSpeaker Z: Voice control, Wi-Fi, Alarm.

	Customer asks: %s`, input)

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a helpful customer service assistant.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})

	if err != nil {
		log.Println("OpenAI error:", err)
		return "Sorry, I'm having trouble responding right now."
	}

	return resp.Choices[0].Message.Content
}

func NewGPTServices() GPTServices {
	return &GPTServicesImpl{}
}
