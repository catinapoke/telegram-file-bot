package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/catinapoke/go-microservice/fileservice"
	"github.com/catinapoke/telegram-file-bot/common"
	"github.com/catinapoke/telegram-file-bot/database"
	"github.com/catinapoke/telegram-file-bot/handler"
	"github.com/go-telegram/bot"
)

func main() {
	fmt.Println("Starting bot!")

	// Database
	var db database.DatabaseOperator
	err := db.Start()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// FileService
	fileserviceUrl := os.Getenv("FILESERVICE_URL")
	if fileserviceUrl == "" {
		fileserviceUrl = "http://127.0.0.1:3001"
	}

	service := fileservice.CreateControllerByUrl(fileserviceUrl)

	// Bot
	token, err := common.GetEnvFromFile("BOT_TOKEN_FILE")

	if err != nil {
		panic(err)
	}

	botHandler := handler.NewBotHandler(db, service)

	opts := []bot.Option{
		bot.WithDefaultHandler(botHandler.Handle),
	}

	b, err := bot.New(token, opts...)

	if err != nil {
		panic(err)
	}

	// Start
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

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
