package main

import (
	"octaaf/models"

	"github.com/jasonlvhit/gocron"
)

func initCrons() {
	gocron.Every(1).Day().At("13:37").Do(sendGlobal, "1337")
	gocron.Every(1).Day().At("16:20").Do(sendGlobal, "420")
	gocron.Every(1).Day().At("00:00").Do(setKaliCount)

	<-gocron.Start()
}

func setKaliCount() {
	lastCount := models.MessageCount{}

	err := DB.Last(&lastCount)

	count := models.MessageCount{
		Count: kaliCount,
		Diff:  0,
	}

	if err == nil {
		count.Diff = (kaliCount - lastCount.Count)
	}

	DB.Save(&count)
}
