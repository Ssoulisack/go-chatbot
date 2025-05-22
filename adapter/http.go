package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-fiber/bootstrap"
	"go-fiber/core/logs"
	"go-fiber/domain/models"
	"io"
	"io/ioutil"
	"log"

	"net/http"
)

type CustomHTTPClient struct {
	Client *http.Client
}

type HttpRequest interface {
	Do(log_key string, req *http.Request) (*http.Response, error)
}

// NewCustomHTTPClient creates a new instance of CustomHTTPClient.
func NewCustomHTTPClient(Client *http.Client) HttpRequest {
	return &CustomHTTPClient{
		Client: Client,
	}
}

func (c *CustomHTTPClient) Do(log_key string, req *http.Request) (*http.Response, error) {
	logs.Info(fmt.Sprintf("key:%v Method %v : %v", log_key, req.Method, req.URL))
	// Log the request payload
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logs.Info(fmt.Sprintf("Error reading request body: %v", err))
		} else {
			logs.Info(fmt.Sprintf("key:%v Request: %s", log_key, body))
			// Reset the request body for subsequent use
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}
	}
	req.Close = true
	// Send the request and get the response
	resp, err := c.Client.Do(req)

	// Log the response for non-GET and non-OPTIONS requests
	if err != nil {
		logs.Info(fmt.Sprintf("Error sending request: %v", err))
	} else if req.Method != "GET" && req.Method != "OPTIONS" {
		defer resp.Body.Close()
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Info(fmt.Sprintf("key:%s Error reading response body: %v", log_key, err))
		} else {
			logs.Info(fmt.Sprintf("key:%v Response: %s", log_key, responseBody))
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(responseBody))
		}
		logs.Info(fmt.Sprintf("REQUEST_PATH# '%v', RESPONSE_STATUS_CODE# '%v' , RESPONSE# '%v'", resp.Request.URL.RawPath, resp.StatusCode, string(responseBody)))
	}

	return resp, err
}

// utils/gemini.go or services/gemini_api.go

func CallGeminiAPI(input string) (string, error) {
	apiKey := bootstrap.GlobalEnv.Keys.GeminiApiKey
	path := bootstrap.GlobalEnv.Gemini.Path
	model := bootstrap.GlobalEnv.Gemini.Model

	url := fmt.Sprintf("%s/%s:generateContent?key=%s", path, model, apiKey)

	requestBody := models.GeminiRequest{
		Contents: []models.GeminiContent{
			{
				Parts: []models.GeminiPart{
					{
						Text: fmt.Sprintf(`You are a smart assistant. Product info:
						- SmartWatch X: Heart rate, Sleep tracker, GPS.
						- SmartSpeaker Z: Voice control, Wi-Fi, Alarm.

						Customer asks: %s`, input),
					},
				},
			},
		},
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("gemini request creation error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gemini API request error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("gemini API error %d: %s\n", resp.StatusCode, string(respBody))
		return "", fmt.Errorf("gemini API returned status %d", resp.StatusCode)
	}

	var geminiResp models.GeminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		log.Println("Failed to parse Gemini response:", err)
		return "", fmt.Errorf("gemini response parsing error: %w", err)
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("gemini returned no answer")
}
