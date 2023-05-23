package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/catinapoke/go-microservice/fileservice"
	"github.com/catinapoke/telegram-file-bot/internal/config"
	"github.com/catinapoke/telegram-file-bot/internal/database"
	"github.com/catinapoke/telegram-file-bot/internal/handler"
	"github.com/go-telegram/bot"
)

const (
	WebhookUrl = ":2000"
)

func main() {
	fmt.Println("Starting bot!")

	if err := config.LoadConfig(); err != nil {
		log.Fatalln(err)
	}

	// Database
	var db database.DatabaseOperator
	err := db.Start()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// FileService
	service := fileservice.CreateControllerByUrl(config.Config.FileServiceUrl)

	// Bot
	botHandler := handler.NewBotHandler(db, service)

	opts := []bot.Option{
		bot.WithDefaultHandler(botHandler.Handle),
	}

	b, err := bot.New(config.Config.BotToken, opts...)

	if err != nil {
		log.Fatalln(err)
	}

	// Start
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	switch config.Config.ListenMode {
	case "requests":
		b.Start(ctx)
	case "webhook":
		go b.StartWebhook(ctx)
		http.ListenAndServe(WebhookUrl, b.WebhookHandler())
	default:
		log.Fatalln(fmt.Errorf("can't define listen mode"))
	}
}
