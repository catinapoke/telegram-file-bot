package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/catinapoke/go-microservice/fileservice"
	"github.com/catinapoke/telegram-file-bot/internal/common"
	"github.com/catinapoke/telegram-file-bot/internal/database"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	StartMessage = "/start"
)

type BotHandler struct {
	db    database.DatabaseOperator
	files fileservice.FileServiceController
}

func NewBotHandler(db database.DatabaseOperator, service fileservice.FileServiceController) *BotHandler {

	return &BotHandler{
		db:    db,
		files: service,
	}
}

func (bh *BotHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	var err error

	defer func() {
		if err != nil { // TODO: Use errors.Is and write right errors output
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Got error while doing request, please write to support!",
			})

			fmt.Printf("Got error: %v\n", err)
		}
	}()

	if update.Message.Text == StartMessage {
		err = bh.start(ctx, b, update)
		return
	}

	if update.Message.Document != nil {
		err = bh.load(ctx, b, update)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}

func (bh *BotHandler) start(ctx context.Context, b *bot.Bot, update *models.Update) error {
	userData := update.Message.From
	user_row := database.User{
		Id:           userData.ID,
		FirstName:    userData.FirstName,
		LastName:     userData.LastName,
		LanguageCode: userData.LanguageCode,
		Username:     userData.Username,
		StartUsage:   time.Now(),
	}

	return bh.db.CreateUser(user_row)
}

func (bh *BotHandler) load(ctx context.Context, b *bot.Bot, update *models.Update) error {
	userId := update.Message.From.ID
	document := update.Message.Document
	file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: document.FileID})

	if err != nil {
		return err
	}

	id, err := bh.files.Set(file.FilePath)

	if err != nil {
		return err
	}

	filename := common.GenerateRandomString(24)
	fileRecord := database.FileName{
		Name:  filename,
		Id:    id,
		Share: 0,
		Owner: userId,
	}

	err = bh.db.CreateFile(fileRecord)
	return err
}

func (bh *BotHandler) get(userId int64, fileId string) {
	/*
		// file id of uploaded image
		inputFileData := "AgACAgIAAxkDAAIBOWJimnCJHQJiJ4P3aasQCPNyo6mlAALDuzEbcD0YSxzjB-vmkZ6BAQADAgADbQADJAQ"
		// or URL image path
		// inputFileData := "https://example.com/image.png"

		params := &bot.SendPhotoParams{
			ChatID:  chatID,
			Photo:   &models.InputFileString{Data: inputFileData},
		}

		bot.SendPhoto(ctx, params)*/
	/*
		fileContent, _ := os.ReadFile("/path/to/image.png")

		params := &bot.SendPhotoParams{
			ChatID:  chatID,
			Photo:   &models.InputFileUpload{Filename: "image.png", Data: bytes.NewReader(fileContent)},
		}

		bot.SendPhoto(ctx, params)*/
}

func (bh *BotHandler) delete(userId int64, fileId string) {

}

func (bh *BotHandler) list(userid int64) {

}
