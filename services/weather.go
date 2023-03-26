package services

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var c = cache.New(60*time.Minute, 60*time.Minute)

func GetWeather(cityCode string) string {
	cacheKey := fmt.Sprintf("weather-%s", strings.ToLower(cityCode))
	weather, found := c.Get(cacheKey)
	if found {
		fmt.Println("Get weather from cache")
		return weather.(string)
	}
	url := fmt.Sprintf("https://service.baomoi.com/weather-%s-4.json", cityCode)
	weather = baoMoiRequest(url)
	c.Set(cacheKey, weather, cache.DefaultExpiration)
	return weather.(string)
}

// baoMoiRequest get array byte from url
func baoMoiRequest(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Referer", "https://baomoi.com/")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("No response from request")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error when close response body")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}
	return string(body)
}
