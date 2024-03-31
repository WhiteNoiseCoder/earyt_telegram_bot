package yt

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	youtube "github.com/kkdai/youtube/v2"
)

type audioCodec uint32

const (
	unknown audioCodec = iota
	opus
	mp4a
)

func parseMimeType(mimeTipe string) audioCodec {
	regexMimeExpression := "(audio|video)/(mp4|webm); codecs=\"(.*)\""
	regMime := regexp.MustCompile(regexMimeExpression)
	match := regMime.FindStringSubmatch(mimeTipe)

	if len(match) == 4 {
		codecsString := match[3]
		if codecsString == "opus" {
			return opus
		}
		regexCodecsExpression := "([^.^,^ ]*)\\.([^.^,^ ]*)*"
		regCodecs := regexp.MustCompile(regexCodecsExpression)
		codecs := regCodecs.FindAllStringSubmatch(codecsString, -1)
		for _, codec := range codecs {

			if codec[1] == "mp4a" {
				return mp4a
			} else if codec[1] == "opus" {
				return opus
			}
		}
	}
	log().Errorf("can't parse mime type %s", mimeTipe)
	return unknown
}

func bestFormatForAudio(formats youtube.FormatList) (*youtube.Format, error) {

	if len(formats) == 0 {
		return nil, fmt.Errorf("there isn't formats in video")
	}

	audioformats := formats.WithAudioChannels()
	if len(audioformats) == 0 {
		return nil, fmt.Errorf("there isn't audio formats in video")
	}

	sort.Slice(audioformats, func(i, j int) bool {
		audioQuality := map[int]int{}
		for _, index := range []int{i, j} {
			if strings.Contains(audioformats[index].AudioQuality, "AUDIO_QUALITY_LOW") {
				audioQuality[index] = 1
			} else if strings.Contains(audioformats[index].AudioQuality, "AUDIO_QUALITY_MEDIUM") {
				audioQuality[index] = 2
			} else if strings.Contains(audioformats[index].AudioQuality, "AUDIO_QUALITY_HIGHT") {
				audioQuality[index] = 3
			}
		}
		if audioQuality[i] != audioQuality[j] {
			return audioQuality[i] < audioQuality[j]
		}

		if formats[i].FPS != formats[j].FPS {
			return formats[i].FPS > formats[j].FPS // less fps is better
		}

		mimeTypeI := parseMimeType(audioformats[i].MimeType)
		mimeTypeJ := parseMimeType(audioformats[j].MimeType)
		return mimeTypeI < mimeTypeJ
	})

	return &audioformats[len(audioformats)-1], nil
}

func (d Downloader) download(url string) (DownloadedInfo, error) {
	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		return DownloadedInfo{}, err
	}

	bestAudioFormat, err := bestFormatForAudio(video.Formats)
	if err != nil {
		return DownloadedInfo{}, fmt.Errorf("error on choice audio format for download: %v", err)
	}

	stream, size, err := client.GetStream(video, bestAudioFormat)
	if err != nil {
		return DownloadedInfo{}, fmt.Errorf("error on get stream for download: %v", err)
	}
	defer stream.Close()

	file, err := os.CreateTemp(d.set.TempPath, "*.mp4")
	if err != nil {
		return DownloadedInfo{}, fmt.Errorf("error on creating temp file for download: %v", err)
	}
	filename := file.Name()

	log().Infof("download youtube file, url: %s to %s, size is: %d, format is: AudioQuality: %s, MimeType: %s, FPS: %d", url, filename, size, bestAudioFormat.AudioQuality, bestAudioFormat.MimeType, bestAudioFormat.FPS)

	_, err = io.Copy(file, stream)
	if err != nil {
		file.Close()
		os.Remove(filename)
		return DownloadedInfo{}, fmt.Errorf("error on saving downloaded stream: %v", err)
	}
	defer file.Close()
	log().Tracef("finish download youtube file, url: %s", url)

	return DownloadedInfo{Path: filename, Name: video.Title}, nil
}
