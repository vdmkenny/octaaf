package scrapers

import (
	"errors"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// GetXKCD returns a random xkcd image or an error
func GetXKCD() ([]byte, error) {
	res, err := http.Get("https://c.xkcd.com/random/comic/")

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return nil, err
	}

	imageURL, found := doc.Find("#comic img").First().Attr("src")

	if !found {
		return nil, errors.New("no 'src' attribute found")
	}

	return loadImage(imageURL)
}
