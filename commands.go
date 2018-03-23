package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"octaaf/scrapers"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	humanize "github.com/dustin/go-humanize"

	"github.com/PuerkitoBio/goquery"
	"github.com/o1egl/govatar"
	"github.com/tidwall/gjson"
	"gopkg.in/telegram-bot-api.v4"
)

func sendRoll(message *tgbotapi.Message) {
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
	reply(message, roll)
}

func count(message *tgbotapi.Message) {
	reply(message, fmt.Sprintf("%v", message.MessageID))
}

func whoami(message *tgbotapi.Message) {
	reply(message, fmt.Sprintf("%v", message.From.ID))
}

func m8Ball(message *tgbotapi.Message) {

	if len(message.CommandArguments()) == 0 {
		reply(message, "Oi! You have to ask question hé 🖕")
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
	Octaaf.Send(msg)
}

func sendBodegem(message *tgbotapi.Message) {
	msg := tgbotapi.NewLocation(message.Chat.ID, 50.8614773, 4.211304)
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func where(message *tgbotapi.Message) {
	argument := strings.Replace(message.CommandArguments(), " ", "+", -1)

	location, found := scrapers.GetLocation(argument)

	if !found {
		reply(message, "This place does not exist 🙈🙈🙈🤔🤔�")
		return
	}

	msg := tgbotapi.NewLocation(message.Chat.ID, location.Lat, location.Lng)
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func what(message *tgbotapi.Message) {
	query := message.CommandArguments()
	resp, _ := http.Get(fmt.Sprintf("https://api.duckduckgo.com/?q=%v&format=json&no_html=1&skip_disambig=1", query))
	body, _ := ioutil.ReadAll(resp.Body)

	result := gjson.Get(string(body), "AbstractText").String()

	if len(result) == 0 {
		reply(message, fmt.Sprintf("What is this *%v* you speak of? 🤔", query))
		return
	}

	reply(message, fmt.Sprintf("*%v:* %v", query, result))
}

func weather(message *tgbotapi.Message) {
	argument := strings.Replace(message.CommandArguments(), " ", "+", -1)

	location, found := scrapers.GetLocation(argument)

	if !found {
		reply(message, "No data found 🙈🙈🙈🤔🤔🤔")
		return
	}

	resp, _ := http.Get(fmt.Sprintf("https://graphdata.buienradar.nl/forecast/json/?lat=%v&lon=%v", location.Lat, location.Lng))
	body, _ := ioutil.ReadAll(resp.Body)
	weatherJSON := string(body)

	msg := "No weather data found."

	forecasts := gjson.Get(weatherJSON, "forecasts").Array()
	raining := false

	if len(forecasts) > 0 {
		msg = "☀️☀️☀️ It's not going to rain in " + message.CommandArguments()
		if forecasts[0].Get("precipation").Num > 0 {
			msg = "🌧🌧🌧 It's now raining in " + message.CommandArguments()
			raining = true
		}
	}

	for _, forecast := range forecasts {
		if raining && forecast.Get("precipation").Num == 0 {
			msg += ", but it's expected to stop "
			rain, err := dateparse.ParseAny(forecast.Get("datetime").String())
			if err != nil {
				msg += " in " + forecast.Get("datetime").String()
			} else {
				msg += humanize.Time(rain)
			}
			break
		} else if forecast.Get("precipation").Num > 0 {
			rain, err := dateparse.ParseAny(forecast.Get("datetime").String())
			if err != nil {
				msg = "🌦🌦🌦 Expected rain from " + forecast.Get("datetime").String()
			} else {
				msg = "🌦🌦🌦 Expected rain " + humanize.Time(rain)
			}
			break
		}
	}

	reply(message, "*Weather:* "+msg)
}

func sendAvatar(message *tgbotapi.Message) {
	img, err := govatar.GenerateFromUsername(govatar.MALE, message.From.UserName)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, img)

	bytes := tgbotapi.FileBytes{Name: "avatar.png", Bytes: buf.Bytes()}
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func bol(message *tgbotapi.Message) {
	bolURL := "https://www.bol.com/nl/nieuwsbrieven.html?country=BE"
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	resp, _ := client.Get(bolURL)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	token := "bogusTokenValue"

	doc.Find(".newsletter_subscriptions input").Each(func(i int, node *goquery.Selection) {
		name, found := node.Attr("name")
		if found && name == "token" {
			token, _ = node.Attr("value")
		}
	})

	data := url.Values{
		"emailAddress":          {message.CommandArguments()},
		"subscribedNewsLetters": {"DAGAANBIEDINGEN", "SOFT_OPTIN", "HARD_OPTIN", "B2B"},
		"token":                 {token},
		"updateNewsletters":     {"Voorkeuren+opslaan"}}

	req, _ := http.NewRequest("POST", bolURL, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36")
	client.Do(req)

	reply(message, fmt.Sprintf("Succesfully subscribed *%v* to the bol.com mailing lists!", message.CommandArguments()))
}

func search(message *tgbotapi.Message) {
	if len(message.CommandArguments()) == 0 {
		reply(message, "What do you expect me to do? 🤔🤔🤔🤔")
		return
	}

	url, found := scrapers.Search(message.CommandArguments(), message.Command() == "search_nsfw")

	if found {
		reply(message, url)
		return
	}

	reply(message, "I found nothing 😱😱😱")
}

func sendStallman(message *tgbotapi.Message) {

	image, err := scrapers.GetStallman()

	if err != nil {
		reply(message, "Stallman went bork? 🤔🤔🤔🤔")
		return
	}

	bytes := tgbotapi.FileBytes{Name: "stally.jpg", Bytes: image}
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

	msg.Caption = message.CommandArguments()
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func sendImage(message *tgbotapi.Message) {
	if len(message.CommandArguments()) == 0 {
		reply(message, fmt.Sprintf("What am I to do, @%v? 🤔🤔🤔🤔", message.From.UserName))
		return
	}

	images, err := scrapers.GetImages(message.CommandArguments(), message.Command() == "img_sfw")

	if err != nil {
		reply(message, "Something went wrong!")
	}

	timeout := time.Duration(2 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	for _, url := range images {

		res, err := client.Get(url)

		if err != nil {
			continue
		}

		content, err := ioutil.ReadAll(res.Body)

		if err != nil {
			continue
		}

		bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: content}
		msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

		msg.Caption = message.CommandArguments()
		msg.ReplyToMessageID = message.MessageID
		_, e := Octaaf.Send(msg)

		if e == nil {
			return
		}
	}

	reply(message, "I did not find images for the query: `"+message.CommandArguments()+"`")
}

func xkcd(message *tgbotapi.Message) {
	image, err := scrapers.GetXKCD()

	if err != nil {
		reply(message, "Failed to parse XKCD image")
	}

	bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: image}
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

	msg.Caption = message.CommandArguments()
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}
