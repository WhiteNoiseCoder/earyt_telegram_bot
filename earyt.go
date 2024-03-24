package main

import (
	"os"

	"github.com/WhiteNoiseCoder/earyt/logger"
	"github.com/WhiteNoiseCoder/earyt/settings"
	"github.com/sirupsen/logrus"
)

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{"channel": "main"})
}

func main() {
	settingspath := "settings.json"
	if len(os.Args) > 1 {
		command := os.Args[1]
		path := os.Args[2]
		if command == "--settings" {
			settingspath = path
		}
	}

	set, err := settings.ParseSettings(settingspath)
	if err != nil {
		log().Errorf("Error in start param: %v\n", err)
		return
	}

	loggerHolder, err := logger.SetUp(&set.Logger)
	if err != nil {
		log().Errorf("Error on setup logger: %v\n", err)
		return
	}
	defer loggerHolder.Close()

	log().Infof("start telegram server")
}
