package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gobuffalo/envy"
	"gopkg.in/telegram-bot-api.v4"
)

var kaliCount *int

func main() {
	envy.Load("config/.env")

	connectDB()
	migrateDB()
	initBot()

	go initCrons()

	kaliID, _ := strconv.ParseInt(os.Getenv("TELEGRAM_ROOM_ID"), 10, 64)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := Octaaf.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if update.Message.Chat.ID == kaliID {
			*kaliCount = update.Message.MessageID
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "roll":
				go sendRoll(update.Message)
			case "m8ball":
				go m8Ball(update.Message)
			case "bodegem":
				go sendBodegem(update.Message)
			case "img", "img_sfw":
				go sendImage(update.Message)
			case "stallman":
				go sendStallman(update.Message)
			case "avatar":
				go sendAvatar(update.Message)
			case "search", "search_nsfw":
				go search(update.Message)
			case "where":
				go where(update.Message)
			case "count":
				go count(update.Message)
			case "weather":
				go weather(update.Message)
			case "what":
				go what(update.Message)
			case "bol":
				go bol(update.Message)
			}
		}

		if update.Message.MessageID%100000 == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("💯💯💯💯 YOU HAVE MESSAGE %v 💯💯💯💯", update.Message.MessageID))
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ParseMode = "markdown"
			Octaaf.Send(msg)
		}
	}
}
