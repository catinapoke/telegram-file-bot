package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/catinapoke/go-microservice/fileservice"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func GetEnvFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("can't retrieve env value %s: %d", path, err)
	}

	return string(data), nil
}

func main() {
	fmt.Println("Hello world!")
	controller := fileservice.CreateController("fileservice", "3001")
	_, err := controller.Get(0)

	if err != nil {
		panic(fmt.Errorf("can't get zero file: %d", err))
	}

	token, err := GetEnvFromFile("BOT_TOKEN_FILE")

	if err != nil {
		panic(err)
	}

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

	go b.StartWebhook(ctx)

	http.ListenAndServe(":2000", b.WebhookHandler())
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
