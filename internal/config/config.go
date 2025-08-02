package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	Debug            bool   `env:"DEBUG" envDefault:"false"`
	GenConfig        GenConfig
	MinioConfig      MinioConfig
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("не удалось загрузить .env: %w", err)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("не удалось запарсить переменные окружения .env: %w", err)
	}

	return &cfg, nil
}
