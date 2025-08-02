package genapi_cli

import (
	"TarotBot/internal/infrastructure/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *GenApi) formatPrompt(initialPrompt, userPrompt, formCards string) string {
	return fmt.Sprintf("%s\nВопрос пользователя:\n%s\nПользователю выпали карты\n%s",
		initialPrompt, userPrompt, formCards)
}

func (c *GenApi) createGenAPIRequest(prompt string) GenAPIRequest {
	return GenAPIRequest{
		Messages: []GenAPIMessage{
			{
				Role: "user",
				Content: []GenAPIContent{
					{
						Type: "text",
						Text: prompt,
					},
				},
			},
		},
		IsSync:           true,
		Stream:           false,
		N:                1,
		FrequencyPenalty: 0,
		MaxTokens:        16384,
		PresencePenalty:  0,
		Temperature:      1,
		TopP:             1,
		ResponseFormat:   `{"type":"text"}`,
	}
}

func (c *GenApi) executeGenAPIRequest(ctx context.Context, requestBody GenAPIRequest) ([]byte, error) {
	l := logger.FromContext(ctx)

	generateURL := fmt.Sprintf("%s/networks/%s", c.BaseURL, c.NetworkID)
	jsonPayload, err := json.Marshal(requestBody)
	if err != nil {
		l.Errorw("не удалось замаршалить реквест для genApi", "ошибка", err)
		return nil, fmt.Errorf("ошибка сериализации запроса: %w", err)
	}
	log.Printf("Запрос к GenAPI: %s", string(jsonPayload))

	req, err := http.NewRequest("POST", generateURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		l.Errorw("не удалось создать http запрос к genApi", "ошибка", err, "url", generateURL)
		return nil, fmt.Errorf("ошибка создания HTTP запроса: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		l.Errorw("не удалось выполнить http запрос к genApi", "ошибка", err, "url", generateURL)
		return nil, fmt.Errorf("ошибка запроса : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		l.Errorf("genApi вернул ошибку (статус %d): %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("ошибка API GenAPI (статус %d): %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Errorw("не удалось прочитать тело ответа от genApi", "ошибка", err)
		return nil, fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	return bodyBytes, nil
}

func (c *GenApi) parseGenAPIResponse(ctx context.Context, bodyBytes []byte) (string, error) {
	l := logger.FromContext(ctx)

	bodyString := string(bodyBytes)
	l.Infow("получен полный ответ от genApi", "текст", bodyString)

	var genAPIResponse GenAPIResponse
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&genAPIResponse); err != nil {
		l.Errorw("не удалось декодировать ответ от genApi", "ошибка", err)
		return "", fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	if len(genAPIResponse.Response) > 0 {
		output := genAPIResponse.Response[0].Message.Content
		return output, nil
	} else {
		l.Errorw("genApi вернул пустой ответ")
		return "", fmt.Errorf("genApi вернул пустой массив output")
	}
}
