package service

import (
	"TarotBot/internal/domain/tarot"
	"TarotBot/internal/infrastructure/logger"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) Rasclad(ctx context.Context, userPromt string) ([]tgbotapi.Chattable, error) {
	l := logger.FromContext(ctx)
	l.Infof("получен userPromt: %s", userPromt)

	var messages []tgbotapi.Chattable

	deck := tarot.NewDeck()
	l.Info("собрана колода")

	cards := s.makeSpread(deck)
	l.Info("карты перетасованы и выбраны", "карты", cards)

	messages, err := s.formatRasclad(ctx, userPromt, cards)
	if err != nil {
		l.Errorw("не удалось сформировать расклад", "ошибка", err, "запрос пользователя", userPromt)
		return nil, fmt.Errorf("ошибка при форматировании расклада: %w", err)
	}
	l.Info("Расклад сформирован")
	return messages, nil
}
