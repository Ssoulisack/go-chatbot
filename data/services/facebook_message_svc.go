package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-fiber/adapter"
	"go-fiber/bootstrap"
	"go-fiber/domain/models"
	"io"
	"net/http"
)

type MessageService interface {
	SendFacebookReply(senderID, reply string) error
}

type MessageServiceImpl struct {
	Client adapter.HttpRequest
}

func NewMessageService(client adapter.HttpRequest) MessageService {
	return &MessageServiceImpl{Client: client}
}

func (s *MessageServiceImpl) SendFacebookReply(senderID, reply string) error {
	var message = models.FacebookReplyRequest{
		Recipient: models.Recipient{
			ID: senderID,
		},
		Message: models.Message{
			Text: reply,
		},
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://graph.facebook.com/v18.0/me/messages?access_token=%s", bootstrap.GlobalEnv.Keys.FbPageAccessToken),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do("FB_REPLY", req)
	if err != nil {
		return fmt.Errorf("error sending Facebook message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body) // <-- ADD THIS
		return fmt.Errorf("facebook API returned non-200: %v, body: %s", resp.StatusCode, responseBody)
	}

	return nil
}
