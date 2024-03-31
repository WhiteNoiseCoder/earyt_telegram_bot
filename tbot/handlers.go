package tbot

import (
	"fmt"
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

func (h Handlers) TStartHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	var startText string = `Hello!
	This is telegram bot for downloading audio from YouTube
		
	/start - write info
	/audio <YouTube URL> - download audio from YouTube
	<YouTube URL> - download audio from YouTube`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, startText)
	_, err := bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error on send start message %v", err)
	}
	return nil
}

// Download youtube audio handler
func (h Handlers) TDownloadYTAudioHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	log().Infof("download youtube %s", update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Prepare ...")
	msg.ReplyToMessageID = update.Message.MessageID
	sendedWaitMessage, errWaitMessage := bot.Send(msg)

	fileInfo, err := h.DownloadAudio(update.Message.Text)
	if err != nil {
		return err
	}
	defer os.Remove(fileInfo.Path)
	log().Debugf("upload audio from %s", update.Message.Text)

	if errWaitMessage == nil {
		bot.Send(tgbotapi.NewDeleteMessage(sendedWaitMessage.Chat.ID, sendedWaitMessage.MessageID))
	}
	safeTitle, err := filenamify.Filenamify(fileInfo.Name, filenamify.Options{MaxLength: 140})
	if err == nil {
		safeTitle += filepath.Ext(fileInfo.Path)
	} else {
		log().Errorf("Error on create human readable filename %v", err)
	}
	audiofileRequest := tgbotapi.NewDocument(update.Message.Chat.ID, FileData{Path: fileInfo.Path, Name: safeTitle})
	audiofileRequest.ReplyToMessageID = update.Message.MessageID
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

	startHandler := trouter.CreateHandlerKit(h.TStartHandler)
	router.AddHandler("^\\/start$", startHandler.Handler)
	downloadYTAudioHandler := trouter.CreateHandlerKit(h.TDownloadYTAudioHandler)
	router.AddHandler("^\\/audio$", downloadYTAudioHandler.Handler)
	router.AddDefaultHandler(downloadYTAudioHandler.Handler)
	router.Run()
}
