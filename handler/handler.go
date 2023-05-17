package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/catinapoke/go-microservice/fileservice"
	"github.com/catinapoke/telegram-file-bot/common"
	"github.com/catinapoke/telegram-file-bot/database"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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

func (hndl *BotHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {

	var err error
	if update.Message.Text == "/start" {
		err = hndl.start(ctx, b, update)
	} else if update.Message.Document != nil {
		err = hndl.load(ctx, b, update)
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   update.Message.Text,
		})
	}

	if err != nil { // TODO: Use errors.Is and write right errors output
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Got error while doing request, please write to support!",
		})

		fmt.Printf("Got error: %v\n", err)
	}
}

func (hndl *BotHandler) start(ctx context.Context, b *bot.Bot, update *models.Update) error {
	userData := update.Message.From
	user_row := database.User{
		Id:           userData.ID,
		FirstName:    userData.FirstName,
		LastName:     userData.LastName,
		LanguageCode: userData.LanguageCode,
		Username:     userData.Username,
		StartUsage:   time.Now(),
	}

	return hndl.db.CreateUser(user_row)
}

func (hndl *BotHandler) load(ctx context.Context, b *bot.Bot, update *models.Update) error {
	userId := update.Message.From.ID
	document := update.Message.Document
	file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: document.FileID})

	if err != nil {
		return err
	}

	id, err := hndl.files.Set(file.FilePath)

	if err != nil {
		return err
	}

	filename := common.GenerateRandomString(24)
	file = database.FileName{}

	hndl.db.CreateFile()
}

func (hndl *BotHandler) get(userId int64, fileId string) {
	controller := fileservice.CreateController("fileservice", "3001")
	file_path, _ := controller.Get(0)
	if file_path == "omagad" { // Just don't want to remove fileservice dependency as I will use it later
		panic(1)
	}

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

func (hndl *BotHandler) delete(userId int64, fileId string) {

}

func (hndl *BotHandler) list(userid int64) {

}
