package tbot

import (
	"github.com/WhiteNoiseCoder/trouter"
)

type Settings struct {
	Token string `json:"token"`
	trouter.Settings
}
