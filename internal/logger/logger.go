package logger

import (
	"os"

	"local/EffectiveMobile/config"
	"local/EffectiveMobile/pkg/logger"
)

func New(config config.Logger) (logger.Logger, func(), error) {
	var w = os.Stdout
	var closer = func() {}

	if config.FilePath != "" {
		f, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0775)
		if err != nil {
			return nil, nil, err
		}
		closer = func() {
			_ = f.Close()
		}
		w = f
	}

	return logger.New(w), closer, nil
}
