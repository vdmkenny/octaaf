package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gobuffalo/pop"
	"github.com/joho/godotenv"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	err := godotenv.Load("config/.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tx, e := pop.Connect("development")

	if e != nil {
		log.Panic(e)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		log.Panic(err)
	}

	kaliID, _ := strconv.ParseInt(os.Getenv("TELEGRAM_ROOM_ID"), 10, 64)
	kaliCount := 0

	go initCrons(bot, tx, &kaliCount)

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "roll":
				go sendRoll(bot, update.Message)
			case "m8ball":
				go m8Ball(bot, update.Message)
			case "bodegem":
				go sendBodegem(bot, update.Message)
			case "img", "img_sfw":
				go sendImage(bot, update.Message)
			case "stallman":
				go sendStallman(bot, update.Message)
			case "avatar":
				go sendAvatar(bot, update.Message)
			case "search", "search_nsfw":
				go search(bot, update.Message)
			case "where":
				go where(bot, update.Message)
			case "count":
				go count(bot, update.Message)
			case "weather":
				go weather(bot, update.Message)
			}
		}

		if update.Message.MessageID%100000 == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("💯💯💯💯 YOU HAVE MESSAGE %v 💯💯💯💯", update.Message.MessageID))
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ParseMode = "markdown"
			bot.Send(msg)
		}

		if update.Message.Chat.ID == kaliID {
			kaliCount = update.Message.MessageID
		}
	}
}
