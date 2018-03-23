package scrapers

import (
	"errors"
	"math/rand"
	"path"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GetStallman returns either a random Stallman image, or an error
func GetStallman() ([]byte, error) {
	var url = "https://stallman.org/photos/rms-working/"

	doc, err := goquery.NewDocument(url)

	if err != nil {
		return nil, err
	}

	var pages []string

	doc.Find("img").Each(func(i int, token *goquery.Selection) {
		url, exists := token.Parent().Attr("href")
		if exists {
			pages = append(pages, url)
		}
	})

	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(len(pages))

	doc, err = goquery.NewDocument(url + pages[roll])

	if err != nil {
		return nil, err
	}

	imagePATH, found := doc.Find("img").First().Parent().Attr("href")

	if !found {
		return nil, errors.New("no 'href' attribute found")
	}

	return loadImage(url + path.Base(imagePATH))
}
