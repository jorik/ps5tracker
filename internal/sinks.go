package internal

import (
	"fmt"

	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ReportSink interface {
	Send(message string) error
}

type telegramSink struct {
	botAPI *tgbotapi.BotAPI
	chatID int64
}

func (t telegramSink) Send(message string) error {
	msg := tgbotapi.NewMessage(t.chatID, message)
	_, err := t.botAPI.Send(msg)
	return err
}

func NewTelegramSink(token string, chatID int64) (ReportSink, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("could not create telegram bot: %w", err)
	}

	return telegramSink{botAPI: bot, chatID: chatID}, nil
}

type printSink struct{}

func (p printSink) Send(message string) error {
	logrus.Print(message)
	return nil
}

func NewPrintSink() ReportSink {
	return printSink{}
}
