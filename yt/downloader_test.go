package yt

import (
	"testing"

	youtube "github.com/kkdai/youtube/v2"
)

// TestParseMimeTypeAudioMp4a calls yt.parseMimeType
// for audio format
func TestParseMimeTypeAudioMp4a(t *testing.T) {
	typeCodec := parseMimeType("audio/mp4; codecs=\"mp4a.40.2\"")
	if typeCodec != mp4a {
		t.Fatalf("type must be mp4a")
	}
}

// TestParseMimeTypeVideo calls yt.parseMimeType
// for audio with video format
func TestParseMimeTypeVideo(t *testing.T) {
	typeCodec := parseMimeType("video/mp4; codecs=\"avc1.64001F, mp4a.40.2\"")
	if typeCodec != mp4a {
		t.Fatalf("type must be mp4a")
	}
}

// TestParseMimeTypeAudioOpus calls yt.parseMimeType
// for audio format opus
func TestParseMimeTypeAudioOpus(t *testing.T) {
	typeCodec := parseMimeType("audio/webm; codecs=\"opus\"")
	if typeCodec != opus {
		t.Fatalf("type must be opus")
	}
}

// TestParseMimeTypeWrong calls yt.parseMimeType
// for wrong audio format
func TestParseMimeTypeWrong(t *testing.T) {
	typeCodec := parseMimeType("wrongtype")
	if typeCodec != unknown {
		t.Fatalf("type must be opus")
	}
}

// TestBestFormatForAudio call bestFormatForAudio
// sort by AudioQuality
func TestBestFormatForAudioByAudioQuality(t *testing.T) {
	format := youtube.FormatList{
		youtube.Format{AudioQuality: "AUDIO_QUALITY_LOW", MimeType: "codecs=\"opus\"", FPS: 0, AudioChannels: 2},
		youtube.Format{AudioQuality: "AUDIO_QUALITY_MEDIUM", MimeType: "video/mp4; codecs=\"avc1.64001F, mp4a.40.2\"", FPS: 24, AudioChannels: 2},
	}
	bestFormat, err := bestFormatForAudio(format)
	if err != nil {
		t.Fatalf("error on sort format %v", err)
	}
	if bestFormat.AudioQuality != "AUDIO_QUALITY_MEDIUM" {
		t.Fatalf("best format AudioQuality must be AUDIO_QUALITY_MEDIUM but %s", bestFormat.AudioQuality)
	}
}

// TestBestFormatForAudioByFPS call bestFormatForAudio
// sort by FPS
func TestBestFormatForAudioByFPS(t *testing.T) {
	format := youtube.FormatList{
		youtube.Format{AudioQuality: "AUDIO_QUALITY_MEDIUM", MimeType: "audio/mp4; codecs=\"mp4a.40.2\"", FPS: 0, AudioChannels: 2},
		youtube.Format{AudioQuality: "AUDIO_QUALITY_MEDIUM", MimeType: "video/mp4; codecs=\"avc1.64001F, mp4a.40.2\"", FPS: 24, AudioChannels: 2},
	}
	bestFormat, err := bestFormatForAudio(format)
	if err != nil {
		t.Fatalf("error on sort format %v", err)
	}
	if bestFormat.FPS != 0 {
		t.Fatalf("best format FPS must be 0 but %d", bestFormat.FPS)
	}
}

// TestBestFormatForAudioByMimeType call bestFormatForAudio
// sort by MimeType
func TestBestFormatForAudioByMimeType(t *testing.T) {
	format := youtube.FormatList{
		youtube.Format{AudioQuality: "AUDIO_QUALITY_MEDIUM", MimeType: "MimeType: audio/webm; codecs=\"opus\"", FPS: 0, AudioChannels: 2},
		youtube.Format{AudioQuality: "AUDIO_QUALITY_MEDIUM", MimeType: "video/mp4; codecs=\"avc1.64001F, mp4a.40.2\"", FPS: 0, AudioChannels: 2},
	}
	bestFormat, err := bestFormatForAudio(format)
	if err != nil {
		t.Fatalf("error on sort format %v", err)
	}
	if bestFormat.MimeType != "MimeType: audio/webm; codecs=\"opus\"" {
		t.Fatalf("best format FPS must be \"MimeType: audio/webm; codecs=\"opus\"\" but %s", bestFormat.MimeType)
	}
}

// TestBestFormatForAudioWithoutAudio call bestFormatForAudio
// without audio
func TestBestFormatForAudioWithoutAudio(t *testing.T) {
	format := youtube.FormatList{
		youtube.Format{AudioQuality: "AUDIO_QUALITY_MEDIUM", MimeType: "MimeType: audio/webm; codecs=\"opus\"", FPS: 0, AudioChannels: 0},
		youtube.Format{AudioQuality: "AUDIO_QUALITY_MEDIUM", MimeType: "video/mp4; codecs=\"avc1.64001F, mp4a.40.2\"", FPS: 0, AudioChannels: 0},
	}
	_, err := bestFormatForAudio(format)
	if err == nil {
		t.Fatalf("error on sort format %v", err)
	}
}
