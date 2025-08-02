package service

import (
	"TarotBot/internal/domain/tarot"
	"TarotBot/internal/infrastructure/logger"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/minio/minio-go/v7"
	"io"
	"math/rand"
	"strings"
	"time"
)

const (
	defaultInitialPromt = "Ты опытный таро-консультант, который делает профессиональные интерпретации карт. Ты должен проанализировать карты, которые выпали пользователю, и высказать свое мнение, опираясь на его запрос. Будь прямым, хладнокровным, говори честно и прямо, не бойся задеть или обидеть. Не лей воды и избегай пафоса, ты современный таролог. Каждый абзац необходимо начинать со смайлика \"🔮\". Не приветствуй и не прощайся - сразу переходи к делу. Количество символов в твоем сообщении должно находиться в диапазоне от 800 до 1000."
	defaultBucketName   = "tarot-images"
	maxCards            = 3
	afterRascladCommand = "🔮 Итак, начнем. Будем гадать в рамках системы Таро Манара:\n\n"
	replyPromt          = "📝 Ты спросил:\n"
	listCards           = "🕯 Тебе выпали следующие карты:\n\n"
	interpretation      = "➡️ Что это может для тебя значить?\n\n"
)

func (s *Service) formatRasclad(ctx context.Context, userPromt string, cards []tarot.Card) ([]tgbotapi.Chattable, error) {
	l := logger.FromContext(ctx)
	var messages []tgbotapi.Chattable

	var sb strings.Builder
	formCards := s.formatCards(sb, cards)
	l.Infow("карты сформированы", formCards)

	genResult, err := s.AiAdapter.GeneratePrediction(ctx, userPromt, defaultInitialPromt, formCards)
	if err != nil {
		l.Errorw("не удалось сгенерировать предсказание", "ошибка", err)
		return nil, fmt.Errorf("formatRasclad: Ошибка при генерации предсказания: %w", err)
	}

	l.Infow("получено предсказание от apiGen", "предсказание", genResult)

	sb.WriteString(afterRascladCommand)
	sb.WriteString(replyPromt)
	sb.WriteString(userPromt)
	sb.WriteString(listCards)
	sb.WriteString(formCards)
	firstMessageText := sb.String()

	mediaGroup, hasImages := s.createMediaGroup(ctx, cards, firstMessageText)

	if hasImages {
		mediaGroupConfig := tgbotapi.NewMediaGroup(0, mediaGroup)
		messages = append(messages, mediaGroupConfig)

	} else {
		l.Info("не удалось получить карточки, отправляем только текст")
		textMessage := tgbotapi.NewMessage(0, firstMessageText)
		messages = append(messages, &textMessage)
	}

	interpretationText := fmt.Sprintf("%s%s", interpretation, genResult)
	interpretationMessage := tgbotapi.NewMessage(0, interpretationText)
	messages = append(messages, &interpretationMessage)

	return messages, nil
}

func (s *Service) makeSpread(deck tarot.Deck) []tarot.Card {
	keys := make([]int, len(deck))
	resultCards := make([]tarot.Card, maxCards)

	for i := 0; i < len(deck); i++ {
		keys[i] = i
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for i := 0; i < maxCards; i++ {
		resultCards[i] = deck[keys[i]]
	}

	return resultCards
}

func (s *Service) createMediaGroup(ctx context.Context, cards []tarot.Card, caption string) ([]interface{}, bool) {
	l := logger.FromContext(ctx)
	mediaGroup := make([]interface{}, 0, len(cards))
	firstImageCaptionSet := false

	for i, card := range cards {
		imageBytes, err := s.getS3ImageBytes(ctx, card.ID, defaultBucketName)
		if err != nil {
			l.Errorw("не получилось взять изображение для карты", "карта", card.ID, "ошибка", err)
			continue
		}

		photo := tgbotapi.FileBytes{Name: fmt.Sprintf("%d.jpg", card.ID), Bytes: imageBytes}
		photoMessage := tgbotapi.NewInputMediaPhoto(photo)

		if i == 0 && !firstImageCaptionSet && caption != "" {
			l.Info("установлен заголовок caption для первого изображения")
			photoMessage.Caption = caption
			firstImageCaptionSet = true
		}

		mediaGroup = append(mediaGroup, photoMessage)
	}

	return mediaGroup, len(mediaGroup) > 0
}

func (s *Service) formatCards(sb strings.Builder, cards []tarot.Card) string {
	for _, c := range cards {
		sb.WriteString(fmt.Sprintf("❗️%s:\n%s\n\n", c.Title, c.Description))
	}

	return sb.String()
}

func (s *Service) getS3ImageBytes(ctx context.Context, cardID int64, bucket string) ([]byte, error) {
	l := logger.FromContext(ctx)
	key := fmt.Sprintf("%d.jpg", cardID)

	object, err := s.ObjectRepoAdapter.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		l.Errorf("ошибка при получениие объекта из MiniO cardID: %d, bucket: %s, key: %s", cardID, bucket, key)
		return nil, fmt.Errorf("ошибка при получении объекта %s из Minio: %w", key, err)
	}
	defer object.Close()

	imageData, err := io.ReadAll(object)
	if err != nil {
		l.Errorw("не удалось прочитать данные изображения", "ошибка", err, "изображение", key)
		return nil, fmt.Errorf("ошибка при чтении данных изображения %s: %w", key, err)
	}

	l.Infof("успешно получено изображение cardID: %d", cardID)
	return imageData, nil
}
