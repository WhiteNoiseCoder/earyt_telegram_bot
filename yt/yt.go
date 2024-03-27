package yt

import "github.com/sirupsen/logrus"

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{"channel": "yt"})
}

// Class which download youtube records
type Downloader struct {
}

// func for download audio from youtube
func (m Downloader) DownloadAudio(url string) error {
	log().Infof("download audio url: %s", url)
	return nil
}
