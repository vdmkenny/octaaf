package main

import (
	"log"
	"os"

	"github.com/gobuffalo/envy"
	"gopkg.in/telegram-bot-api.v4"
)

// Octaaf is the global bot endpoint
var Octaaf *tgbotapi.BotAPI

func initBot() {
	var err error
	Octaaf, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		log.Panic(err)
	}

	env := envy.Get("GO_ENV", "development")
	Octaaf.Debug = env == "development"

	log.Printf("Authorized on account %s", Octaaf.Self.UserName)
}

func reply(message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = "markdown"
	Octaaf.Send(msg)
}
