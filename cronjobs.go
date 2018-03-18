package main

import (
	"log"
	"time"

	"github.com/markbates/pop"
	"github.com/octaaf/models"

	"github.com/jasonlvhit/gocron"
	"gopkg.in/telegram-bot-api.v4"
)

func initCrons(bot *tgbotapi.BotAPI, tx *pop.Connection, kaliCountRef *int) {
	gocron.Every(1).Day().At("13:37").Do(sendGlobal, bot, "1337")
	gocron.Every(1).Day().At("16:20").Do(sendGlobal, bot, "420")
	gocron.Every(1).Day().At("00:00").Do(setKaliCount, tx, kaliCountRef)

	<-gocron.Start()
}

func setKaliCount(tx *pop.Connection, kaliCount *int64) {
	stat := models.KaliStat{MessageCount: *kaliCount, Date: time.Now()}

	_, err := tx.ValidateAndCreate(&stat)

	if err != nil {
		log.Printf("Unable to save the kali message count: %v", err)
	}
}
