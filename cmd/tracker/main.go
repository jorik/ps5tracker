package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/jorik/ps5tracker/internal"
	"gopkg.in/yaml.v3"
)

var (
	configFile       = flag.String("configFile", "./tests.yaml", "Config file")
	telegramChatID   = flag.Int64("telegramChatId", -1, "")
	telegramBotToken = flag.String("telegramBotToken", "", "Leave empty to output to STDOUT")
	logLevel         = flag.String("loglevel", "info", "Log level: [debug, info, warning]")
)

func main() {
	flag.Parse()

	lvl, _ := logrus.ParseLevel(*logLevel)
	logrus.Infof("Setting log level to: %s", lvl)
	logrus.SetLevel(lvl)

	c, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Printf("Unable to read config file contents: %v\n", err)
	}

	var config internal.Config
	if err = yaml.Unmarshal(c, &config); err != nil {
		log.Printf("Could not unmarshal config json")
	}

	done := make(chan struct{})
	results := make(chan string)

	for _, c := range config.Websites {
		go internal.StartScrape(c, done, results)
	}

	var sink internal.ReportSink

	if *telegramBotToken == "" {
		sink = internal.NewPrintSink()
	} else {
		sink, err = internal.NewTelegramSink(*telegramBotToken, *telegramChatID)
		if err != nil {
			logrus.Fatalf("Could not create Telegram Sink: %v", err)
		}
	}

	startupMessage := fmt.Sprintf("🎮🕵🏻 Starting PS5 Tracker with %d tracked websites:", len(config.Websites))
	for _, c := range config.Websites {
		startupMessage += fmt.Sprintf("\n - %s", c.Url)
	}
	_ = sink.Send(startupMessage)

	for res := range results {
		if err := sink.Send(res); err != nil {
			logrus.Errorf("Could not send message to sink: %v", err)
		}
	}
}
