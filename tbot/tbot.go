package tbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Server struct {
	bot *tgbotapi.BotAPI
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{"channel": "tbot"})
}

// return UserName which registred in telegram bot
func (ser *Server) UserName() string {
	return ser.bot.Self.UserName
}

func StartServer(set Settings) (*Server, error) {
	server := new(Server)
	var err error
	server.bot, err = tgbotapi.NewBotAPI(set.Token)
	if err != nil {
		return nil, fmt.Errorf("failed create telegram bot: %v", err)
	}
	return server, nil
}
