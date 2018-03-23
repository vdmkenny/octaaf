package scrapers

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Search searches on duckduck go & returns the first url
func Search(query string, nsfw bool) (string, bool) {
	url := "https://duckduckgo.com/lite?k1=-1&q=" + query

	if nsfw {
		url += "&kp=-2"
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36")
	resp, err := client.Do(req)

	if err != nil {
		return "", false
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return "", false
	}

	return doc.Find(".result-link").First().Attr("href")
}
