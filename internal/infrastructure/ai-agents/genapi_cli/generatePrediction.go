package genapi_cli

import (
	"TarotBot/internal/infrastructure/logger"
	"context"
	"fmt"
)

func (c *GenApi) GeneratePrediction(ctx context.Context, userPrompt, initialPrompt, formCards string) (string, error) {
	l := logger.FromContext(ctx)

	prompt := c.formatPrompt(initialPrompt, userPrompt, formCards)

	requestBody := c.createGenAPIRequest(prompt)

	bodyBytes, err := c.executeGenAPIRequest(ctx, requestBody)
	if err != nil {
		l.Errorw("не удалось выполнить http запрос к genApi", "ошибка", err)
		return "", fmt.Errorf("не удалось выполнить http запрос к genApi: %w", err)
	}

	output, err := c.parseGenAPIResponse(ctx, bodyBytes)
	if err != nil {
		l.Errorw("не удалось декодировать ответ и извлечь контент от  genApi", "ошибка", err)
		return "", fmt.Errorf("не удалось декодировать ответ и извлечь контент от  genApi: %w", err)
	}

	return output, nil
}
