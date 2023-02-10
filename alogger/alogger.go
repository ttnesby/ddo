package alogger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"time"
)

const skipFrameCount = 3 // 3 is the caller of functions in this package

type ALogger struct {
	zerolog.Logger
}

func New(debug bool) ALogger {

	// override default long names for Caller
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
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
