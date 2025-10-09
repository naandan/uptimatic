package utils

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func InitLogger(levelStr string) {
	var level zerolog.Level
	switch levelStr {
	case "DEBUG":
		level = zerolog.DebugLevel
	case "INFO":
		level = zerolog.InfoLevel
	case "WARN", "WARNING":
		level = zerolog.WarnLevel
	case "ERROR":
		level = zerolog.ErrorLevel
	case "FATAL":
		level = zerolog.FatalLevel
	default:
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	Logger = log.Output(os.Stdout)
}

func commonFields(c *gin.Context, extra map[string]any) map[string]any {
	fields := map[string]any{}

	if c != nil {
		if reqID, exists := c.Get("request_id"); exists {
			fields["request_id"] = reqID
		}
		fields["endpoint"] = c.FullPath()
	}

	if extra != nil {
		fields["extra"] = extra
	}

	return fields
}

func Debug(c *gin.Context, msg string, extra map[string]any) {
	fields := commonFields(c, extra)
	Logger.Debug().Fields(fields).Msg(msg)
}

func Info(c *gin.Context, msg string, extra map[string]any) {
	fields := commonFields(c, extra)
	Logger.Info().Fields(fields).Msg(msg)
}

func Warn(c *gin.Context, msg string, extra map[string]any) {
	fields := commonFields(c, extra)
	Logger.Warn().Fields(fields).Msg(msg)
}

func Error(c *gin.Context, msg string, extra map[string]any) {
	fields := commonFields(c, extra)
	Logger.Error().Fields(fields).Msg(msg)
}

func Fatal(c *gin.Context, msg string, extra map[string]any) {
	fields := commonFields(c, extra)
	Logger.Fatal().Fields(fields).Msg(msg)
}
