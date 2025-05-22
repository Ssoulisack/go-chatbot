package bootstrap

import (
	"net/http"
	"time"
)

func NewHttpClient() *http.Client {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return client
}
