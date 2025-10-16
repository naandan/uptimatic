package utils

import (
	"context"
	"os"
	"time"

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

	Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	// Logger = log.Output(os.Stdout)
}

func commonFields(ctx context.Context, extra map[string]any) map[string]any {
	fields := map[string]any{}
	if ctx != nil {
		if reqID, ok := ctx.Value("request_id").(string); ok {
			fields["request_id"] = reqID
		}
	}
	if extra != nil {
		fields["extra"] = extra
	}
	return fields
}

func Debug(c context.Context, msg string, extra map[string]any) {
	fields := commonFields(c, extra)
	Logger.Debug().Fields(fields).Msg(msg)
}

func Info(c context.Context, msg string, extra map[string]any) {
	fields := commonFields(c, extra)
	Logger.Info().Fields(fields).Msg(msg)
}

func Warn(c context.Context, msg string, extra map[string]any) {
	fields := commonFields(c, extra)
	Logger.Warn().Fields(fields).Msg(msg)
}

func Error(c context.Context, msg string, extra map[string]any) {
	fields := commonFields(c, extra)
	Logger.Error().Fields(fields).Msg(msg)
}

func Fatal(c context.Context, msg string, extra map[string]any) {
	fields := commonFields(c, extra)
	Logger.Fatal().Fields(fields).Msg(msg)
}
