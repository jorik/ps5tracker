package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/jorik/ps5tracker/internal"
	"gopkg.in/yaml.v3"
)

var (
	configFile       = flag.String("configFile", "./tests.yaml", "Config file")
	telegramChatID   = flag.Int64("telegramChatId", -1, "")
	telegramBotToken = flag.String("telegramBotToken", "", "")
	debug            = flag.Bool("debug", false, "Whether to only output to stdout")
	logLevel         = flag.String("loglevel", "info", "Log level: [debug, info, warning]")
)

func main() {
	flag.Parse()

	lvl, _ := logrus.ParseLevel(*logLevel)
	logrus.Infof("Setting log level to: %s", lvl)
	logrus.SetLevel(lvl)

	if *debug {
		logrus.Infof("Running debug mode")
	}
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

	var sinks []internal.ReportSink

	if *debug || *telegramBotToken == "" {
		sink := internal.NewPrintSink()
		sinks = append(sinks, sink)
	} else {
		telegramSink, err := internal.NewTelegramSink(*telegramBotToken, *telegramChatID)
		if err != nil {
			logrus.Fatalf("Could not create Telegram Sink: %v", err)
		}
		sinks = append(sinks, telegramSink)
	}

	for res := range results {
		for _, sink := range sinks {
			if err := sink.Send(res); err != nil {
				logrus.Errorf("Could not send message to sink: %v", err)
			}
		}
	}
}
