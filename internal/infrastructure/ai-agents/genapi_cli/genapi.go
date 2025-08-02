package genapi_cli

import (
	"context"
	"net/http"
	"time"
)

const (
	BaseGenAPIURL = "https://api.gen-api.ru/api/v1"
)

type AiAdapter interface {
	GeneratePrediction(ctx context.Context, userPrompt, initialPrompt, formCards string) (string, error)
}

type GenApi struct {
	Token      string
	NetworkID  string
	BaseURL    string
	HttpClient *http.Client
}

func NewClient(token, networkID string) AiAdapter {
	return &GenApi{
		Token:      token,
		NetworkID:  networkID,
		BaseURL:    BaseGenAPIURL,
		HttpClient: &http.Client{Timeout: 30 * time.Second},
	}
}
