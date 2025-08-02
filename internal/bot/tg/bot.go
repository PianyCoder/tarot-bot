package tg

import (
	"TarotBot/internal/infrastructure/logger"
	"TarotBot/internal/service"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type Bot struct {
	bot                     *tgbotapi.BotAPI
	Servicer                service.Servicer
	waitingForUserPrompt    map[int64]bool
	lastRascladRequestTimes sync.Map
}

func NewBot(bot *tgbotapi.BotAPI, service service.Servicer) *Bot {
	return &Bot{
		bot:                     bot,
		Servicer:                service,
		waitingForUserPrompt:    make(map[int64]bool),
		lastRascladRequestTimes: sync.Map{},
	}
}

func (b *Bot) Start(ctx context.Context) error {
	logger.FromContext(ctx).Infof("авторизирован аккаунт %s", b.bot.Self.UserName)

	updates := b.initUpdatesChannel()

	b.handleUpdates(ctx, updates)

	return nil
}

func (b *Bot) handleUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			_ = b.handleCommand(ctx, update.Message)
			continue
		}

		b.handleMessage(ctx, update.Message)
	}
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
