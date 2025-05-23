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
						Text: fmt.Sprintf(`You are an Apple product assistant. Here's the product information:

						- iPhone 15 Pro:
						- A17 Pro chip
						- 6.1-inch Super Retina XDR display
						- ProMotion technology with 120Hz refresh rate
						- Triple camera system (48MP Main, Ultra Wide, Telephoto)
						- USB-C port and Dynamic Island
						- Titanium design, available in four finishes

						- iPhone 15:
						- A16 Bionic chip
						- 6.1-inch Super Retina XDR display
						- Dual-camera system (48MP Main, Ultra Wide)
						- USB-C charging
						- Ceramic Shield front, aluminum body

						- MacBook Pro 14-inch (M3 Pro):
						- Apple M3 Pro chip with 11-core CPU and 14-core GPU
						- 14.2-inch Liquid Retina XDR display
						- Up to 18 hours of battery life
						- MagSafe 3, Thunderbolt 4 ports
						- ProMotion, True Tone

						- MacBook Air 13-inch (M2):
						- Apple M2 chip
						- 13.6-inch Liquid Retina display
						- Fanless design
						- MagSafe charging
						- Up to 18 hours battery

						- Apple Watch Series 9:
						- S9 SiP chip with Double Tap gesture
						- Always-On Retina display
						- Blood Oxygen, ECG, Heart Rate monitoring
						- Workout tracking, sleep tracking
						- Water resistant to 50m

						- Apple Watch Ultra 2:
						- 49mm rugged titanium case
						- Up to 36 hours of battery life (72 in Low Power Mode)
						- Precision dual-frequency GPS
						- Depth gauge and EN13319 certified for diving

						- AirPods Pro (2nd Generation):
						- H2 chip for better noise cancellation and sound
						- Adaptive Transparency
						- Spatial Audio with dynamic head tracking
						- MagSafe charging case with Precision Finding

						- iPad Pro 12.9-inch (M2):
						- Apple M2 chip
						- Liquid Retina XDR display
						- Apple Pencil hover support
						- 5G capable
						- Face ID, USB-C with Thunderbolt

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
