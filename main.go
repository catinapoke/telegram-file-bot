package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/catinapoke/go-microservice/fileservice"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var database DatabaseOperator

func main() {
	fmt.Println("Hello world!")
	controller := fileservice.CreateController("fileservice", "3001")
	file_path, _ := controller.Get(0)
	if file_path == "omagad" { // Just don't want to remove fileservice dependency as I will use it later
		panic(1)
	}

	token, err := GetEnvFromFile("BOT_TOKEN_FILE")

	if err != nil {
		panic(err)
	}

	err = database.Start()
	if err != nil {
		panic(err)
	}
	defer database.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)

	if err != nil {
		panic(err)
	}
	// call methods.SetWebhook if needed

	mode := os.Getenv("LISTEN_MODE")
	if mode == "requests" {
		b.Start(ctx)
	} else if mode == "webhook" {
		go b.StartWebhook(ctx)
		http.ListenAndServe(":2000", b.WebhookHandler())
	} else {
		panic(fmt.Errorf("can't define listen mode"))
	}
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})

	if update.Message.Text == "/start" {
		start(ctx, b, update)
	}
}

func start(ctx context.Context, b *bot.Bot, update *models.Update) {

	userData := update.ChatMember.From
	user_row := Users{
		Id:           userData.ID,
		FirstName:    userData.FirstName,
		LastName:     userData.LastName,
		LanguageCode: userData.LanguageCode,
		Username:     userData.Username,
		StartUsage:   time.Now(),
	}

	err := database.CreateUser(&user_row)
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

}

func delete(userId int64, fileId string) {

}

func list(userid int64) {

}
