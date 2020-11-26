package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/cfg"
)

const (
	LoggerLevelDebug = "DEBUG"
	LoggerLevelInfo  = "INFO"
	LoggerLevelWarn  = "WARN"
	LoggerLevelError = "ERROR"
	LoggerLevelFatal = "FATAL"
	LoggerLevelPanic = "PANIC"
)

func NewLogger(cfg *cfg.AppCfg) zerolog.Logger {
	switch cfg.Logger.Level {
	case LoggerLevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LoggerLevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LoggerLevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LoggerLevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LoggerLevelFatal:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	}

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    cfg.Logger.NoColoredOutput,
		TimeFormat: time.RFC3339,
	}

	output.FormatLevel = func(i interface{}) string {
		var v string

		if ii, ok := i.(string); ok {
			ii = strings.ToUpper(ii)
			switch ii {
			case LoggerLevelDebug, LoggerLevelError, LoggerLevelFatal, LoggerLevelInfo, LoggerLevelWarn, LoggerLevelPanic:
				v = fmt.Sprintf("%-5s", ii)
			default:
				v = ii
			}
		}

		return fmt.Sprintf("| %s |", v)
	}

	return zerolog.New(output).With().Timestamp().Logger()
}
