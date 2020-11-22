package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/corpix/uarand"
	"github.com/sirupsen/logrus"
)

const (
	// scrape interval is hardcoded to 5 seconds now
	scrapeInterval = time.Second * 5
)

var (
	httpClient = http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Second * 10,
	}
)

func StartScrape(target Website, done chan struct{}, result chan string) {
	logrus.Infof("Starting scrape for %s", target.Url)

	for {
		select {
		case <-time.Tick(scrapeInterval):
			if doScrape(target) {
				result <- fmt.Sprintf("🎉🎮 PS5 beschikbaar! %s", target.Url)
			}
		case <-done:
			return
		}
	}
}

func doScrape(target Website) bool {
	req, err := http.NewRequest("GET", target.Url, nil)
	if err != nil {
		logrus.Errorf("Could not prepare HTTP request for url %s: %v", target.Url, err)
		return false
	}

	// Set a random User-Agent to try to avoid being blocked by websites
	req.Header.Set("User-Agent", uarand.GetRandom())

	res, err := httpClient.Do(req)
	if err != nil {
		logrus.Errorf("Could not do HTTP request for URL %s: %v", target.Url, err)
		return false
	}

	if res.StatusCode != http.StatusOK {
		logrus.Infof("Received a non-200 statusCode (%s): %s", res.Status, target.Url)
		return false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Errorf("Could not read response body: %v", err)
		return false
	}

	str := string(body)

	for _, keyword := range target.Keywords {
		if !strings.Contains(str, keyword) {
			logrus.Debugf("KeywordAppears failed: %s did not contain keyword '%s'", target.Url, keyword)
			return false
		}

		logrus.Debugf("KeywordAppears: %s contains keyword '%s'", target.Url, keyword)
	}

	for _, keyword := range target.KeywordsNotAppearing {
		if strings.Contains(str, keyword) {
			logrus.Debugf("KeywordNotAppearing failed: %s contains blocked keyword '%s'", target.Url, keyword)
			return false
		}

		logrus.Debugf("KeywordNotAppearing: %s does not container keyword '%s'", target.Url, keyword)
	}

	return true
}
