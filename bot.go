package main

import (
	"fmt"
	"log"
	"octaaf/models"
	"os"
	"strconv"
	"strings"

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

	KaliID, _ = strconv.ParseInt(os.Getenv("TELEGRAM_ROOM_ID"), 10, 64)
	ReporterID, _ = strconv.Atoi(envy.Get("REPORTER_ID", "-1"))
}

func reply(message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = "markdown"
	Octaaf.Send(msg)
}

func handle(message *tgbotapi.Message) {
	if message.Chat.ID == KaliID {
		KaliCount = message.MessageID

		if message.From.ID == ReporterID &&
			(strings.ToLower(message.Text) == "reported" || message.Sticker.FileID == "CAADBAAD5gEAAreTBA3s5qVy8bxHfAI") {
			DB.Save(&models.Report{})
		}
	}

	if message.IsCommand() {
		switch message.Command() {
		case "roll":
			sendRoll(message)
		case "m8ball":
			m8Ball(message)
		case "bodegem":
			sendBodegem(message)
		case "img", "img_sfw":
			sendImage(message)
		case "stallman":
			sendStallman(message)
		case "avatar":
			sendAvatar(message)
		case "search", "search_nsfw":
			search(message)
		case "where":
			where(message)
		case "count":
			count(message)
		case "weather":
			weather(message)
		case "what":
			what(message)
		case "bol":
			bol(message)
		}
	}

	if message.MessageID%100000 == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("💯💯💯💯 YOU HAVE MESSAGE %v 💯💯💯💯", message.MessageID))
		msg.ReplyToMessageID = message.MessageID
		msg.ParseMode = "markdown"
		Octaaf.Send(msg)
	}

}
