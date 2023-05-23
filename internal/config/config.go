package config

import (
	"os"

	"github.com/catinapoke/telegram-file-bot/internal/common"
)

type AppConfig struct {
	BotToken       string
	FileServiceUrl string
	ListenMode     string

	Database struct {
		Url      string
		Name     string
		Password string
		Username string
	}
}

var Config = AppConfig{}

func LoadConfig() error {
	config := &AppConfig{}
	var err error

	config.ListenMode = os.Getenv("LISTEN_MODE")

	config.BotToken, err = common.GetEnvFromFile("BOT_TOKEN_FILE")

	if err != nil {
		return common.NewConfigLoadError("BOT_TOKEN_FILE", err)
	}

	fileserviceUrl := os.Getenv("FILESERVICE_URL")
	if fileserviceUrl == "" {
		fileserviceUrl = "http://127.0.0.1:3001"
	}
	config.FileServiceUrl = fileserviceUrl

	addr, err := common.GetEnv("DATABASE_URL")
	if err != nil {
		return common.NewConfigLoadError("DATABASE_URL", err)
	}
	config.Database.Url = addr

	user, err := common.GetEnvFromFile("POSTGRES_USER_FILE")
	if err != nil {
		return common.NewConfigLoadError("POSTGRES_USER_FILE", err)
	}
	config.Database.Username = user

	password, err := common.GetEnvFromFile("POSTGRES_PASSWORD_FILE")
	if err != nil {
		return common.NewConfigLoadError("POSTGRES_PASSWORD_FILE", err)
	}
	config.Database.Password = password

	db, err := common.GetEnvFromFile("POSTGRES_DB_FILE")
	if err != nil {
		return common.NewConfigLoadError("POSTGRES_DB_FILE", err)
	}
	config.Database.Name = db

	return nil
}
