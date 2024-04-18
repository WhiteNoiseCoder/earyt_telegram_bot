package yt

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func (d Downloader) convertToAudio(filename string) (string, error) {
	audiofilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".mp3"
	log().Debugf("convert video %s to audio %s", filename, audiofilename)
	var buf_out bytes.Buffer
	var buf_err bytes.Buffer

	output := ffmpeg.Input(filename, ffmpeg.KwArgs{}).
		Output(audiofilename, ffmpeg.KwArgs{})

	if len(d.set.FfmpegPath) > 0 {
		output.SetFfmpegPath(d.set.FfmpegPath)
	}

	err := output.WithOutput(&buf_out, &buf_err).Run()

	if err != nil {
		log().Errorf("error on ffmpeg convert video %s to audio %s : %s", filename, audiofilename, buf_err.String())
		return "", fmt.Errorf("error on ffmpeg convert video %s to audio %s: %v", filename, audiofilename, err)
	}

	audiofile, err := os.Open(audiofilename)
	if err != nil {
		log().Errorf("error on open audio file for getting size %s", audiofilename)
	}
	audiofileStat, err := audiofile.Stat()
	var audiofileSize int64
	if err == nil {
		audiofileSize = audiofileStat.Size()
	} else {
		audiofileSize = 0
		log().Errorf("error on getting audio file %s size %v", audiofilename, err)
	}

	log().Tracef("finish convert video %s to audio %s (size %d) %s", filename, audiofilename, audiofileSize, buf_out.String())

	return audiofilename, nil
}
