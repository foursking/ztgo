package log

import (
	"fmt"
	"os"
	"sync"

	"git.code.oa.com/qdgo/core/config/env"
	"git.code.oa.com/qdgo/core/metadata"

	mlog "github.com/micro/go-micro/v2/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// key names in log entries
const (
	_timeKey       = "time"
	_levelKey      = "level"
	_nameKey       = "logger"
	_messageKey    = "msg"
	_stacktraceKey = "stacktrace"
)

var (
	logger *Logger   // global business logger
	once   sync.Once // ensures that global logger be instantiated only once

	// log encoded options
	encCfg = zapcore.EncoderConfig{
		TimeKey:        _timeKey,
		LevelKey:       _levelKey,
		NameKey:        _nameKey,
		CallerKey:      metadata.Caller,
		MessageKey:     _messageKey,
		StacktraceKey:  _stacktraceKey,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
)

func init() {
	once.Do(func() {
		if logger == nil {
			Init()
		}
	})
}

// Init inits log
func Init(opts ...Option) {
	options := DefaultOptions()
	for _, o := range opts {
		o(&options)
	}
	lvl := zap.NewAtomicLevelAt(zapLevel(options.Level))
	logger = &Logger{level: lvl}
	cores := zapCores(&options, &lvl)
	tee := zapcore.NewTee(cores...)
	basicFields := []zap.Field{
		{Key: "env", Type: zapcore.StringType, String: env.DeployEnv},
		{Key: "app", Type: zapcore.StringType, String: env.AppName},
		{Key: "hostname", Type: zapcore.StringType, String: env.Hostname},
		{Key: "server", Type: zapcore.StringType, String: env.ServerIP},
	}
	l := zap.New(tee,
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.AddStacktrace(zap.DPanicLevel),
		zap.Fields(basicFields...),
	)
	// logger for business codes
	logger.l = l.Sugar()
	// replaces go-micro logger
	mlog.DefaultLogger = logger
}

func zapCores(opts *Options, lvl *zap.AtomicLevel) (cores []zapcore.Core) {
	enc := zapcore.NewJSONEncoder(encCfg)
	if !opts.StdoutOff {
		sw := zapcore.AddSync(zapcore.Lock(os.Stdout))
		cores = append(cores, zapcore.NewCore(enc, sw, lvl))
	}
	if opts.File.Dir != "" {
		fw := fileWriter(&opts.File)
		cores = append(cores, zapcore.NewCore(enc, zapcore.AddSync(fw), lvl))
	}
	if opts.Kafka.Topic != "" && len(opts.Kafka.Brokers) > 0 {
		kw, err := newKafkaWriter(&opts.Kafka)
		if err != nil {
			fmt.Printf("log: new kafka writer(%+v) error(%v)", opts.Kafka, err)
			return
		}
		cores = append(cores, zapcore.NewCore(enc, zapcore.AddSync(kw), lvl))
	}
	return
}

// Sync flushes buffered logs
// call it before application stopped
func Sync() (err error) {
	// https://github.com/uber-go/zap/issues/370
	// because sync to stdout always has error, ignore here temporary
	_ = logger.l.Sync()
	return
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	logger.Log(mlog.DebugLevel, args...)
}

// Info logs a info message
func Info(args ...interface{}) {
	logger.Log(mlog.InfoLevel, args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	logger.Log(mlog.WarnLevel, args...)
}

// Error logs a error message
func Error(args ...interface{}) {
	logger.Log(mlog.ErrorLevel, args...)
}

// Fatal logs an error message then exit
func Fatal(args ...interface{}) {
	logger.Log(mlog.FatalLevel, args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	logger.Logf(mlog.DebugLevel, format, args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	logger.Logf(mlog.InfoLevel, format, args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	logger.Logf(mlog.WarnLevel, format, args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	logger.Logf(mlog.ErrorLevel, format, args...)
}

// Fatalf logs a formatted fatal message then exit
func Fatalf(format string, args ...interface{}) {
	logger.Logf(mlog.FatalLevel, format, args...)
}

// Debugw logs a debug message with some additional key-value pairs
func Debugw(msg string, kvs ...interface{}) {
	logger.Logw(mlog.DebugLevel, msg, kvs...)
}

// Infow logs a info message with some additional key-value pairs
func Infow(msg string, kvs ...interface{}) {
	logger.Logw(mlog.InfoLevel, msg, kvs...)
}

// Warnw logs a warning message with some additional key-value pairs
func Warnw(msg string, kvs ...interface{}) {
	logger.Logw(mlog.WarnLevel, msg, kvs...)
}

// Errorw logs a error message with some additional key-value pairs
func Errorw(msg string, kvs ...interface{}) {
	logger.Logw(mlog.ErrorLevel, msg, kvs...)
}

// Fatalw logs a fatal message with some additional key-value pairs, then exit
func Fatalw(msg string, kvs ...interface{}) {
	logger.Logw(mlog.FatalLevel, msg, kvs...)
}
