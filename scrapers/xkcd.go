package scrapers

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
)

// GetXKCD returns a random xkcd image or an error
func GetXKCD() ([]byte, error) {
	doc, err := goquery.NewDocument("https://c.xkcd.com/random/comic/")

	if err != nil {
		return nil, err
	}

	imageURL, found := doc.Find("#comic img").First().Attr("src")

	if !found {
		return nil, errors.New("no 'src' attribute found")
	}

	return loadImage(imageURL)
}
