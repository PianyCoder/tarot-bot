package service

import (
	"TarotBot/internal/infrastructure/ai-agents/genapi_cli"
	"TarotBot/internal/infrastructure/minio_cli"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Servicer interface {
	Rasclad(ctx context.Context, userPromt string) ([]tgbotapi.Chattable, error)
}

type Service struct {
	ObjectRepoAdapter minio_cli.ObjectRepoAdapter
	AiAdapter         genapi_cli.AiAdapter
}

func NewService(ObjectRepoAdapter minio_cli.ObjectRepoAdapter, AiAdapter genapi_cli.AiAdapter) Servicer {
	return &Service{
		ObjectRepoAdapter: ObjectRepoAdapter,
		AiAdapter:         AiAdapter,
	}
}
