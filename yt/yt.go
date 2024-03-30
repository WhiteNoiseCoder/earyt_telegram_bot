package yt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	youtube "github.com/kkdai/youtube/v2"
	"github.com/sirupsen/logrus"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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

func (d Downloader) download(url string) (string, error) {
	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		return "", err
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return "", fmt.Errorf("error on get stream for download: %v", err)
	}
	defer stream.Close()

	file, err := os.CreateTemp(d.set.TempPath, "*.mp4")
	if err != nil {
		return "", fmt.Errorf("error on creating temp file for download: %v", err)
	}
	filename := file.Name()

	log().Infof("download audio url: %s to %s", url, filename)

	_, err = io.Copy(file, stream)
	if err != nil {
		file.Close()
		os.Remove(filename)
		return "", fmt.Errorf("error on saving downloaded stream: %v", err)
	}
	defer file.Close()
	return filename, nil
}

func (d Downloader) convertToAudio(filename string) (string, error) {
	audiofilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".mp3"
	var buf_out bytes.Buffer
	var buf_err bytes.Buffer
	err := ffmpeg.Input(filename, ffmpeg.KwArgs{}).
		Output(audiofilename, ffmpeg.KwArgs{}).SetFfmpegPath(d.set.FfmpegPath).
		WithOutput(&buf_out, &buf_err).Run()

	if err != nil {
		log().Errorf("error on ffmpeg convert video %s to audio %s : %s", filename, audiofilename, buf_err.String())
		return "", fmt.Errorf("error on ffmpeg convert video %s to audio %s: %v", filename, audiofilename, err)
	}
	if log().Logger.GetLevel() == logrus.TraceLevel {
		log().Tracef("convert video %s to audio %s: %s", filename, audiofilename, buf_out.String())
	} else {
		log().Tracef("convert video %s to audio %s", filename, audiofilename)
	}
	return audiofilename, nil
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
