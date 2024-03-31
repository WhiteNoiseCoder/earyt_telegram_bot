package yt

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{"channel": "yt"})
}

// Class which download youtube records
type Downloader struct {
	set Settings
}

func CreateDownloader(set Settings) Downloader {
	return Downloader{set: set}
}

// func for download video from youtume and convert to audio
func (d Downloader) DownloadAudio(url string) (string, error) {
	videofilename, err := d.download(url)
	if err != nil {
		return "", fmt.Errorf("error on download from youtube: %v", err)
	}
	defer os.Remove(videofilename)
	audiofilename, err := d.convertToAudio(videofilename)
	if err != nil {
		return "", fmt.Errorf("error on convert to audio: %v", err)
	}
	return audiofilename, err
}
