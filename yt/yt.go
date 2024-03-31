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

// Data of youtube downloaded content
type DownloadedInfo struct {
	Path string
	Name string
}

// func for download video from youtume and convert to audio
func (d Downloader) DownloadAudio(url string) (DownloadedInfo, error) {
	videoFileInfo, err := d.download(url)
	if err != nil {
		return DownloadedInfo{}, fmt.Errorf("error on download from youtube: %v", err)
	}
	defer os.Remove(videoFileInfo.Path)
	audiofilename, err := d.convertToAudio(videoFileInfo.Path)
	if err != nil {
		return DownloadedInfo{}, fmt.Errorf("error on convert to audio: %v", err)
	}
	return DownloadedInfo{Path: audiofilename, Name: videoFileInfo.Name}, err
}
