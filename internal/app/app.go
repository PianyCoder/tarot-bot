package app

import (
	"TarotBot/internal/bot/tg"
	"TarotBot/internal/config"
	"TarotBot/internal/infrastructure/ai-agents/genapi_cli"
	"TarotBot/internal/infrastructure/logger"
	"TarotBot/internal/infrastructure/minio_cli"
	"TarotBot/internal/service"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(ctx context.Context) error {
	zapLog, err := logger.Init()
	if err != nil {
		return fmt.Errorf("не удалось инициализировать логгер: %w", err)
	}

	ctx = logger.WithContext(ctx, zapLog)
	l := logger.FromContext(ctx)

	cfg, err := config.Load()
	if err != nil {
		l.Errorw("не удалось инициализировать конфиг:", "ошибка", err)
		return fmt.Errorf("не удалось инициализировать конфиг: %w", err)
	}

	minioCli, err := minio_cli.NewClient(cfg.MinioConfig.BaseURL, cfg.MinioConfig.User, cfg.MinioConfig.Password, cfg.MinioConfig.Port, cfg.MinioConfig.UseSSL)
	if err != nil {
		l.Errorw("не удалось инициализировать минио:", "ошибка", err)
		return fmt.Errorf("не удалось инициализировать минио: %w", err)
	}

	genCli := genapi_cli.NewClient(cfg.GenConfig.GenAPIToken, cfg.GenConfig.GenAPINetworkID)

	service := service.NewService(minioCli, genCli)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		l.Errorw("не удалось инициализировать бота:", "ошибка", err)
		return fmt.Errorf("не удалось инициализировать бота: %w", err)
	}

	bot.Debug = cfg.Debug

	tgBot := tg.NewBot(bot, service)

	if err = tgBot.Start(ctx); err != nil {
		l.Errorw("не удалось запустить бота:", "ошибка", err)
		return fmt.Errorf("не удалось запустить бота: %w", err)
	}

	return nil
}
