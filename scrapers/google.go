package scrapers

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

// Location contains the latitude & longitude
type Location struct {
	Lat float64
	Lng float64
}

// GetImages searches on google for images & returns an array of image urls
func GetImages(query string, safe bool) ([]string, error) {
	// Replace spaces with '+'
	query = strings.Replace(query, " ", "+", -1)

	url := "http://images.google.com/search?tbm=isch&q=" + query

	if safe {
		url += "&safe=on"
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36")
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return nil, err
	}

	var images []string

	doc.Find(".rg_di .rg_meta").Each(func(i int, token *goquery.Selection) {
		imageJSON := token.Text()
		imageURL := gjson.Get(imageJSON, "ou").String()

		if len(imageURL) > 0 {
			images = append(images, imageURL)
		}
	})

	return images, nil
}

// GetLocation returns a location based on the google maps API
func GetLocation(query string) (Location, bool) {
	res, err := http.Get("https://maps.google.com/maps/api/geocode/json?address=" + query + "&key=" + os.Getenv("GOOGLE_API_KEY"))

	if err != nil {
		return Location{0, 0}, false
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return Location{0, 0}, false
	}

	defer res.Body.Close()

	json := string(body)

	if !gjson.Get(json, "results.0.geometry.location").Exists() {
		return Location{0, 0}, false
	}

	location := Location{
		Lat: gjson.Get(json, "results.0.geometry.location.lat").Num,
		Lng: gjson.Get(json, "results.0.geometry.location.lng").Num}

	return location, true
}
