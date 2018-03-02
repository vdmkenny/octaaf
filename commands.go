package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
	"gopkg.in/telegram-bot-api.v4"
)

func sendRoll(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	rand.Seed(time.Now().UnixNano())
	roll := strconv.Itoa(rand.Intn(9999999999-1000000000) + 1000000000)
	points := [9]string{"👌 Dubs", "🙈 Trips", "😱 Quads", "🤣😂 Penta", "👌👌🤔🤔😂😂 Hexa", "🙊🙉🙈🐵 Septa", "🅱️Octa", "💯💯💯 El Niño"}
	var dubscount int8 = -1

	for i := len(roll) - 1; i > 0; i-- {
		if roll[i] == roll[i-1] {
			dubscount++
		} else {
			break
		}
	}

	if dubscount > -1 {
		roll = points[dubscount] + " " + roll
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, roll)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func m8Ball(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	if len(message.CommandArguments()) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Oi! You have to ask question hé 🖕")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	answers := [20]string{"👌 It is certain",
		"👌 It is decidedly so",
		"👌 Without a doubt",
		"👌 Yes definitely",
		"👌 You may rely on it",
		"👌 As I see it, yes",
		"👌 Most likely",
		"👌 Outlook good",
		"👌 Yes",
		"👌 Signs point to yes",
		"☝ Reply hazy try again",
		"☝ Ask again later",
		"☝ Better not tell you now",
		"☝ Cannot predict now",
		"☝ Concentrate and ask again",
		"🖕 Don't count on it",
		"🖕 My reply is no",
		"🖕 My sources say no",
		"🖕 Outlook not so good",
		"🖕 Very doubtful"}
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(19)
	msg := tgbotapi.NewMessage(message.Chat.ID, answers[roll])
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func sendBodegem(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewLocation(message.Chat.ID, 50.8614773, 4.211304)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func sendImage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	if len(message.CommandArguments()) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "What do you expect me to do? 🤔🤔🤔🤔")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	query := "http://images.google.com/search?tbm=isch&q=" + strings.Replace(message.CommandArguments(), " ", "+", -1)

	if message.Command() == "img_sfw" {
		query += "&safe=on"
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", query, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36")
	resp, _ := client.Do(req)

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	var images []string

	doc.Find(".rg_di .rg_meta").Each(func(i int, token *goquery.Selection) {
		imageJSON := token.Text()
		imageURL, err := jsonparser.GetString([]byte(imageJSON), "ou")

		if err == nil {
			images = append(images, imageURL)
		}
	})

	for _, url := range images {

		res, _ := http.Get(url)

		content, _ := ioutil.ReadAll(res.Body)
		bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: content}
		msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

		msg.Caption = message.CommandArguments()
		msg.ReplyToMessageID = message.MessageID
		_, e := bot.Send(msg)

		if e == nil {
			return
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "I did not find images for the query: `"+message.CommandArguments()+"`")
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = "markdown"
	bot.Send(msg)
}
