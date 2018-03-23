package scrapers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/araddon/dateparse"
	humanize "github.com/dustin/go-humanize"
	"github.com/tidwall/gjson"
)

// GetWeatherStatus reports when it's raining somewhere
func GetWeatherStatus(query string) (string, bool) {
	argument := strings.Replace(query, " ", "+", -1)
	location, found := GetLocation(argument)

	if !found {
		return "", false
	}

	resp, err := http.Get(fmt.Sprintf("https://graphdata.buienradar.nl/forecast/json/?lat=%v&lon=%v", location.Lat, location.Lng))

	if err != nil {
		return "", false
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", false
	}

	weatherJSON := string(body)

	msg := "No weather data found."

	forecasts := gjson.Get(weatherJSON, "forecasts").Array()
	raining := false

	if len(forecasts) > 0 {
		msg = "☀️☀️☀️ It's not going to rain in " + query
		if forecasts[0].Get("precipation").Num > 0 {
			msg = "🌧🌧🌧 It's now raining in " + query
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

	return msg, true
}
