package main

import (
	"github.com/sirupsen/logrus"
)

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{"channel": "main"})
}

func main() {
	log().Infof("start")
}
