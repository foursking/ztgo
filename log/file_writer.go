package log

import (
	"fmt"

	"github.com/foursking/ztgo/config/env"

	"gopkg.in/natefinch/lumberjack.v2"
)

type fileWriterOptions struct {
	Dir        string `toml:"dir" json:"dir"`                 // 日志文件目录
	MaxSize    int    `toml:"max_size" json:"max_size"`       // 日志切割大小，超过就切割新文件，单位 MB
	MaxBackups int    `toml:"max_backups" json:"max_backups"` // 保留最多日志文件个数，超过就清理
	MaxAge     int    `toml:"max_age" json:"max_age"`         // 日志文件保留天数
	Compress   bool   `toml:"compress" json:"compress"`       // 是否压缩
}

type fileLogger struct {
	*lumberjack.Logger
}

// Sync closes file logger
func (f *fileLogger) Sync() error {
	return f.Logger.Close()
}

func fileWriter(opt *fileWriterOptions) *fileLogger {
	logName := env.AppName
	if logName == "" {
		logName = "ztgo"
	}
	return &fileLogger{Logger: &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.log", opt.Dir, logName),
		MaxSize:    opt.MaxSize,
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge,
		Compress:   opt.Compress,
	}}
}
