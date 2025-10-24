package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
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

	stdoutWriter := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stdout
		w.TimeFormat = time.RFC3339Nano
	})
	stderrWriter := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stderr
		w.TimeFormat = time.RFC3339Nano
	})

	Logger = zerolog.New(levelSplitter{
		stdout: writerAdapter{stdoutWriter},
		stderr: writerAdapter{stderrWriter},
	}).With().Timestamp().Logger()
}

func InitSentry(dsn string) error {
	if dsn == "" {
		return nil
	}
	return sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		TracesSampleRate: 1.0,
	})
}

type writerAdapter struct {
	w zerolog.ConsoleWriter
}

func (a writerAdapter) Write(p []byte) (n int, err error) {
	return a.w.Write(p)
}

func (a writerAdapter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	return a.w.Write(p)
}

type levelSplitter struct {
	stdout zerolog.LevelWriter
	stderr zerolog.LevelWriter
}

func (s levelSplitter) Write(p []byte) (n int, err error) {
	return s.stdout.Write(p)
}

func (s levelSplitter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level >= zerolog.ErrorLevel {
		return s.stderr.WriteLevel(level, p)
	}
	return s.stdout.WriteLevel(level, p)
}

func commonFields(ctx context.Context, extra map[string]any) map[string]any {
	fields := map[string]any{}
	if ctx != nil {
		if reqID, ok := ctx.Value(TraceKey).(string); ok {
			fields[string(TraceKey)] = reqID
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
	sendToSentry(msg, fields)
	Logger.Error().Fields(fields).Msg(msg)
}

func Fatal(c context.Context, msg string, extra map[string]any) {
	fields := commonFields(c, extra)

	sendToSentry(msg, fields)
	sentry.Flush(2 * time.Second)

	Logger.Fatal().Fields(fields).Msg(msg)
	os.Exit(1)
}

func sendToSentry(msg string, fields map[string]any) {
	sentry.WithScope(func(scope *sentry.Scope) {
		for k, v := range fields {
			if k == string(TraceKey) {
				scope.SetTag("trace_id", fmt.Sprintf("%v", v))
			} else {
				scope.SetExtra(k, v)
			}
		}
		sentry.CaptureMessage(msg)
	})
}
