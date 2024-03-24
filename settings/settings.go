package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/WhiteNoiseCoder/earyt/logger"
)

type Settings struct {
	Logger logger.Settings `json:"log"`
}

func ParseSettings(settingsPath string) (Settings, error) {

	absSettingsPath, _ := filepath.Abs(settingsPath)
	settingsFile, err := os.Open(absSettingsPath)
	if err != nil {
		return Settings{}, fmt.Errorf("Can't read setting file: %v", err)
	}
	defer settingsFile.Close()

	settingsJsonData := make([]byte, 2046)
	size, err := settingsFile.Read(settingsJsonData)
	if err != nil {
		return Settings{}, fmt.Errorf("Can't read setting file: %v", err)
	}

	settingsJsonData = settingsJsonData[0:size]

	var settings Settings
	err = json.Unmarshal(settingsJsonData, &settings)
	if err != nil {
		return Settings{}, fmt.Errorf("Can't read setting file: %v", err)
	}

	return settings, nil
}
