package tbot

import (
	"github.com/WhiteNoiseCoder/trouter"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DownloadYTAudioHandler = func(yturl string) error

// YT records manage interface
type YT interface {
	DownloadAudio(string) error
}

// Telegram interfaces for handle user query
type Handlers struct {
	YT
}

// Download youtube audio handler
func (h Handlers) TDownloadYTAudioHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	err := h.DownloadAudio(update.Message.Text)
	return err
}

// Start telegram query handling
func (ser *Server) Start(h *Handlers, set *trouter.Settings) {
	router := trouter.NewTRouter(ser.bot, set)

	downloadYTAudioHandler := trouter.CreateHandlerKit(h.TDownloadYTAudioHandler)
	router.AddHandler("^\\/audio$", downloadYTAudioHandler.Handler)
	router.AddDefaultHandler(downloadYTAudioHandler.Handler)
	router.Run()
}
