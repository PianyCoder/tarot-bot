package tg

import (
	"TarotBot/internal/infrastructure/logger"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	cooldownDuration    = 30 * time.Second
	maxUserPromptLength = 350
	//responses
	overMaxPromtLength = "⚠️ Хм, многовато ты написал. Попробуй сосредоточиться на главном!"
	afterStart         = "💜 Привет! Меня зовут Андрей, и я готов стать твоим личным тарологом!\n\n➡️ Напиши команду '/rasclad', чтобы мы могли продолжить..."
	sendMsgForRasclad  = "➡️ Итак, я тебя слушаю. Задай свой вопрос и, если хочешь получить более точную интерпретацию, добавь немного контекста... Расскажи, что тебя волнует?"
	unknownMsg         = "➡️ Я тебя не понимаю... Если хочешь получить расклад, используй команду /rasclad"
	timeWarning        = "⚠️ Не торопись! Ты слишком часто вызываешь команду /rasclad"
	//commands
	commandStart   = "start"
	commandRasclad = "rasclad"
)

func (b *Bot) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	l := logger.FromContext(ctx)
	l.Infow("сообщение пользователя", "пользователь", message.From.UserName, "текст сообщения", message.Text)

	if promptActive, ok := b.waitingForUserPrompt[message.Chat.ID]; ok && promptActive {
		userPrompt := message.Text

		if len(userPrompt) > maxUserPromptLength {
			l.Warnw("сообщение пользователя превышает лимит", "длина", len(userPrompt), "максимум", maxUserPromptLength)
			errorMsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%s\nРазмер твоего вопроса не должен превышать %d символов!", overMaxPromtLength, maxUserPromptLength))
			_, err := b.bot.Send(errorMsg)
			if err != nil {
				l.Errorw("не удалось отправить сообщение об ошибке длины", "ошибка", err, "сообщение", errorMsg)
			}
			return
		}

		delete(b.waitingForUserPrompt, message.Chat.ID)

		messages, err := b.Servicer.Rasclad(ctx, userPrompt)
		if err != nil {
			l.Errorw("не удалось сгенерировать расклад", "ошибка", err)
			errMsg := tgbotapi.NewMessage(message.Chat.ID, "Ой, что-то пошло не так. Попробуй еще раз!")
			_, _ = b.bot.Send(errMsg)
			return
		}

		for _, msg := range messages {
			switch v := msg.(type) {
			case *tgbotapi.MessageConfig:
				v.ChatID = message.Chat.ID
				_, err = b.bot.Send(*v)
				if err != nil {
					l.Errorw("не удалось отправить текстовое сообщение", "ошибка", err, "сообщение", msg)
				}
			case tgbotapi.MediaGroupConfig:
				v.ChatID = message.Chat.ID
				_, err = b.bot.Send(v)
				if err != nil {
					l.Errorw("не удалось отправить медиа-группу", "ошибка", err, "медиа-группа", v)
				}
			default:
				l.Errorw("неизвестный тип сообщения", "ошибка", msg)
			}
		}
		return
	}

	l.Info("пользователь отправил неизвестное сообщение")
	msg := tgbotapi.NewMessage(message.Chat.ID, unknownMsg)
	_, _ = b.bot.Send(msg)
}

func (b *Bot) handleCommand(ctx context.Context, message *tgbotapi.Message) error {
	l := logger.FromContext(ctx)
	l.Infow("сообщение пользователя", "пользователь", message.From.UserName, "текст сообщения", message.Text)

	switch message.Command() {
	case commandStart:
		msg := tgbotapi.NewMessage(message.Chat.ID, afterStart)
		_, err := b.bot.Send(msg)
		if err != nil {
			l.Errorw("не удалось отправить сообщение", "ошибка", err, "сообщение", msg)
			return fmt.Errorf("ошибка при отправке сообщения: %w", err)
		}
		return nil
	case commandRasclad:
		if !b.checkRascladRateLimit(message.Chat.ID) {
			msg := tgbotapi.NewMessage(message.Chat.ID, timeWarning)
			_, _ = b.bot.Send(msg)
			l.Warnw("пользователь вызывает команду /rasclad слишком быстро!", "chatID", message.Chat.ID)
			return nil
		}

		promptMsg := tgbotapi.NewMessage(message.Chat.ID, sendMsgForRasclad)
		_, err := b.bot.Send(promptMsg)
		if err != nil {
			l.Errorw("не удалось отправить сообщение", "ошибка", err, "сообщение", promptMsg)
			return fmt.Errorf("ошибка при отправке сообщения с запросом вопроса: %w", err)
		}

		b.waitingForUserPrompt[message.Chat.ID] = true

		return nil

	default:
		errMsg := tgbotapi.NewMessage(message.Chat.ID, unknownMsg)
		_, err := b.bot.Send(errMsg)
		if err != nil {
			l.Errorw("не удалось отправить сообщение о неизвестной команде", "ошибка", err, "сообщение", errMsg)
			return fmt.Errorf("ошибка при отправке сообщения о неизвестной команде: %w", err)
		}
		return nil
	}
}
