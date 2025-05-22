package services

import (
	"go-fiber/adapter"
	"go-fiber/core/logs"
)

type GeminiServices interface {
	GenerateReply(input string) string
}

type GeminiServicesImpl struct{}

func NewGeminiServices() GeminiServices {
	return &GeminiServicesImpl{}
}

func (s *GeminiServicesImpl) GenerateReply(input string) string {
	reply, err := adapter.CallGeminiAPI(input)
	if err != nil {
		logs.Error("Error calling Gemini API: " + err.Error())
		return "Sorry, I have no answer at the moment."
	}
	return reply
}
