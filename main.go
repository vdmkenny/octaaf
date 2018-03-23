package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gobuffalo/envy"
	"gopkg.in/telegram-bot-api.v4"
)

// KaliCount is an integer that holds the ID of the last send message in the Kali group
var KaliCount int

// KaliID is the ID of the kali group
var KaliID int64

// ReporterID is the id of the user who reports everyone
var ReporterID int

func main() {
	envy.Load("config/.env")

	connectDB()
	migrateDB()
	initBot()

	go initCrons()

	sendGlobal("I'm back up and running! 👌")

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		sendGlobal("I'm going down for onderhoud! ☝️")
		DB.Close()
		os.Exit(0)
	}()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := Octaaf.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		go handle(update.Message)
	}
}
