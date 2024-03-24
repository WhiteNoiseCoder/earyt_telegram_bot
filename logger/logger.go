package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is logger holder struct
type Holder struct {
	file *os.File
}

// SetUp is Logger constructor
func SetUp(settings *Settings) (*Holder, error) {
	holder := new(Holder)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	if len(settings.Path) > 0 {
		var err error
		holder.file, err = os.OpenFile(settings.Path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed open file to log: %v", err)
		}
		logrus.SetOutput(holder.file)
	}
	return holder, nil
}

func (h *Holder) Close() {
	if h.file != nil {
		h.file.Close()
	}
}
