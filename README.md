# ps5tracker

This application is written in [Go](http://golang.org/).

PS5 Tracker is a sample application that scrapes websites every 5 seconds looking for certain phrases. If a phrase, or combination of phrases, is found, it will send a telegram message.

The phrases/websites are loaded from a YAML file. When running locally, it will use the `cmd/tracker/tests/yaml` file. When running in Kubernetes, it's using a mounted config file.

# Running PS5Tracker


To run this application, execute the following command in `cmd/tracker`:

```
go run main.go --telegramBotToken $BOT_TOKEN --telegramChatId $CHAT_ID
```
