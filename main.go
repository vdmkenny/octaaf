package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		log.Panic(err)
	}

	go cronJobs(bot)

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
	}
}

func cronJobs(bot *tgbotapi.BotAPI) {
	gocron.Every(1).Day().At("13:37").Do(sendGlobal, bot, "1337")
	gocron.Every(1).Day().At("16:20").Do(sendGlobal, bot, "420")

	<-gocron.Start()
}

func sendGlobal(bot *tgbotapi.BotAPI, message string) {
	// Wait 1.5 seconds because Telegram has bad NTP
	time.Sleep(1500)

	telegramRoomID, e := strconv.ParseInt("-1001090867629", 10, 64)

	if e != nil {
		log.Println("Invalid Telegram room id, not sending global messages.")
	}

	msg := tgbotapi.NewMessage(telegramRoomID, message)
	_, err := bot.Send(msg)

	if err != nil {
		log.Printf("Error while sending '%s': %s", message, err)
	}
}
