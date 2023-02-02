package alogger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

const skipFrameCount = 3 // 3 is the caller of functions in this package

type ALogger struct {
	zerolog.Logger
}

func New() ALogger {
	//TODO need to add flag to enable debug mode
	//zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	return ALogger{
		log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().
			CallerWithSkipFrameCount(skipFrameCount).
			Logger(),
	}
}

func Disable() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func (aLogger ALogger) Infof(format string, v ...interface{}) {
	aLogger.Info().Msgf(format, v...)
}

func (aLogger ALogger) Error(e error) error {
	aLogger.Err(e).Send()
	return e
}

func (aLogger ALogger) Errorf(format string, v ...interface{}) {
	aLogger.Err(fmt.Errorf(format, v...)).Send()
}

func (aLogger ALogger) Debugf(format string, v ...interface{}) {
	aLogger.Debug().Msgf(format, v...)
}
