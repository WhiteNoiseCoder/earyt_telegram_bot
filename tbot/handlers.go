package tbot

import (
	"os"
	"path/filepath"

	"github.com/WhiteNoiseCoder/earyt/yt"
	"github.com/WhiteNoiseCoder/trouter"

	"github.com/flytam/filenamify"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// YT records manage interface
type YT interface {
	DownloadAudio(string) (yt.DownloadedInfo, error)
}

// Telegram interfaces for handle user query
type Handlers struct {
	YT
}

// Download youtube audio handler
func (h Handlers) TDownloadYTAudioHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	log().Infof("download youtube %s", update.Message.Text)
	fileInfo, err := h.DownloadAudio(update.Message.Text)
	if err != nil {
		return err
	}
	defer os.Remove(fileInfo.Path)
	log().Debugf("upload audio from %s", update.Message.Text)

	safeTitle, err := filenamify.Filenamify(fileInfo.Name, filenamify.Options{MaxLength: 140})
	if err == nil {
		safeTitle += filepath.Ext(fileInfo.Path)
	} else {
		log().Errorf("Error on create human readable filename %v", err)
	}

	audiofileRequest := tgbotapi.NewDocument(update.Message.Chat.ID, FileData{Path: fileInfo.Path, Name: safeTitle})
	_, err = bot.Send(audiofileRequest)
	if err != nil {
		log().Errorf("Error on send file %v", err)
	}
	log().Tracef("finish audio from %s", update.Message.Text)
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
