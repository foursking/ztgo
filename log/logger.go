package log

import (
	mlog "github.com/micro/go-micro/v2/logger"
	"go.uber.org/zap"
)

// Logger qdgo logger
type Logger struct {
	l     *zap.SugaredLogger
	level zap.AtomicLevel
	opts  mlog.Options
}

func (l *Logger) Init(options ...mlog.Option) error {
	return nil
}

func (l *Logger) Options() mlog.Options {
	return l.opts
}

func (l *Logger) Error(err error) mlog.Logger {
	return l
}

func (l *Logger) Fields(fields map[string]interface{}) mlog.Logger {
	return l
}

func (l *Logger) Log(level mlog.Level, v ...interface{}) {
	switch level {
	case mlog.DebugLevel:
		l.l.Debug(v...)
	case mlog.InfoLevel:
		l.l.Info(v...)
	case mlog.WarnLevel:
		l.l.Warn(v...)
	case mlog.ErrorLevel:
		l.l.Error(v...)
	case mlog.FatalLevel:
		l.l.Fatal(v...)
	}
}

func (l *Logger) Logf(level mlog.Level, format string, v ...interface{}) {
	switch level {
	case mlog.DebugLevel:
		l.l.Debugf(format, v...)
	case mlog.InfoLevel:
		l.l.Infof(format, v...)
	case mlog.WarnLevel:
		l.l.Warnf(format, v...)
	case mlog.ErrorLevel:
		l.l.Errorf(format, v...)
	case mlog.FatalLevel:
		l.l.Fatalf(format, v...)
	}
}

func (l *Logger) Logw(level mlog.Level, msg string, kvs ...interface{}) {
	switch level {
	case mlog.DebugLevel:
		l.l.Debugw(msg, kvs...)
	case mlog.InfoLevel:
		l.l.Infof(msg, kvs...)
	case mlog.WarnLevel:
		l.l.Warnf(msg, kvs...)
	case mlog.ErrorLevel:
		l.l.Errorf(msg, kvs...)
	case mlog.FatalLevel:
		l.l.Fatalf(msg, kvs...)
	}
}

func (l *Logger) String() string {
	return "qdgo logger"
}
