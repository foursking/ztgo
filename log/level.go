package log

import (
	mlog "github.com/micro/go-micro/v2/logger"
	"go.uber.org/zap/zapcore"
)

func goMicroLevel(l string) mlog.Level {
	switch l {
	case "trace":
		return mlog.TraceLevel
	case "debug":
		return mlog.DebugLevel
	case "info":
		return mlog.InfoLevel
	case "warn":
		return mlog.WarnLevel
	case "error":
		return mlog.ErrorLevel
	case "fatal":
		return mlog.FatalLevel
	default:
		return mlog.InfoLevel
	}
}

func zapLevel(l string) zapcore.Level {
	switch l {
	case "trace", "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// SetLevel sets logger level
func SetLevel(l string) {
	lo := mlog.WithLevel(goMicroLevel(l))
	mlog.DefaultLogger = mlog.NewLogger(lo)
	logger.level.SetLevel(zapLevel(l))
	Infof("log level changed to %s", l)
}
