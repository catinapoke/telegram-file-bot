package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/catinapoke/go-microservice/fileservice"
	"github.com/catinapoke/telegram-file-bot/database"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type BotHandler struct {
	db database.DatabaseOperator
}

func NewBotHandler(db database.DatabaseOperator) *BotHandler {
	return &BotHandler{
		db: db,
	}
}

func (hndl *BotHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})

	if update.Message.Text == "/start" {
		hndl.start(ctx, b, update)
	}
}

func (hndl *BotHandler) start(ctx context.Context, b *bot.Bot, update *models.Update) {

	userData := update.Message.From
	user_row := database.User{
		Id:           userData.ID,
		FirstName:    userData.FirstName,
		LastName:     userData.LastName,
		LanguageCode: userData.LanguageCode,
		Username:     userData.Username,
		StartUsage:   time.Now(),
	}

	err := hndl.db.CreateUser(user_row)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Got error with creating user, please write to support!",
		})

		fmt.Printf("Got error: %v\n", err)
	}
}

func load(userId int64, telegramFileId string) {

}

func get(userId int64, fileId string) {
	controller := fileservice.CreateController("fileservice", "3001")
	file_path, _ := controller.Get(0)
	if file_path == "omagad" { // Just don't want to remove fileservice dependency as I will use it later
		panic(1)
	}
}

func delete(userId int64, fileId string) {

}

func list(userid int64) {

}
