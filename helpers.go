package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
	"gopkg.in/telegram-bot-api.v4"
)

type location struct {
	lat float64
	lng float64
}

func getLocation(query string) (location, bool) {
	resp, _ := http.Get("https://maps.google.com/maps/api/geocode/json?address=" + query + "&key=" + os.Getenv("GOOGLE_API_KEY"))
	body, _ := ioutil.ReadAll(resp.Body)
	json := string(body)

	if !gjson.Get(json, "results.0.geometry.location").Exists() {
		return location{0, 0}, false
	}

	location := location{
		lat: gjson.Get(json, "results.0.geometry.location.lat").Num,
		lng: gjson.Get(json, "results.0.geometry.location.lng").Num}

	return location, true
}

func getMessageConfig(message *tgbotapi.Message, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = "markdown"
	return msg
}

func sendGlobal(message string) {
	// Wait 1.5 seconds because Telegram has bad NTP
	time.Sleep(1500)

	telegramRoomID, e := strconv.ParseInt(os.Getenv("TELEGRAM_ROOM_ID"), 10, 64)

	if e != nil {
		log.Println("Invalid Telegram room id, not sending global messages.")
	}

	msg := tgbotapi.NewMessage(telegramRoomID, message)
	_, err := Octaaf.Send(msg)

	if err != nil {
		log.Printf("Error while sending '%s': %s", message, err)
	}
}
