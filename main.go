package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Bot
	token, err := common.GetEnvFromFile("BOT_TOKEN_FILE")

	if err != nil {
		panic(err)
	}

	botHandler := handler.NewBotHandler(db)

	opts := []bot.Option{
		bot.WithDefaultHandler(botHandler.Handle),
	}

	b, err := bot.New(token, opts...)

	if err != nil {
		panic(err)
	}

	// Start
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
