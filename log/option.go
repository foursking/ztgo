package log

// Options logger options
type Options struct {
	// Level only logs higher than Level will be recorded
	Level string `toml:"level" json:"level"`

	// StdoutOff switch for standard output, default is false (prints log to stdout)
	StdoutOff bool `toml:"stdout_off" json:"stdout_off"`

	// File file writer
	File fileWriterOptions `toml:"file" json:"file"`

	// Kafka kafka writer
	Kafka kafkaWriterOptions `toml:"kafka" json:"kafka"`
}

// Option function to set logger options
type Option func(*Options)

const (
	_defaultLogLevel = "info"

	// file hook params
	_defaultMaxBackups = 50
	_defaultMaxAge     = 7
	_defaultRotateSize = 500
)

// DefaultOptions gets default log options
func DefaultOptions() Options {
	return Options{
		Level: _defaultLogLevel,
		File: fileWriterOptions{
			MaxSize:    _defaultRotateSize,
			MaxBackups: _defaultMaxBackups,
			MaxAge:     _defaultMaxAge,
			Compress:   false,
		},
	}
}

// SetOptions converts log config file to log options
func SetOptions(lo *Options) Option {
	return func(opts *Options) {
		if lo == nil {
			return
		}
		if lo.Level != "" {
			opts.Level = lo.Level
		}
		if lo.StdoutOff {
			opts.StdoutOff = true
		}
		if lo.File.Dir != "" {
			opts.File.Dir = lo.File.Dir
			if lo.File.MaxAge > 0 {
				opts.File.MaxAge = lo.File.MaxAge
			}
			if lo.File.MaxSize > 0 {
				opts.File.MaxSize = lo.File.MaxSize
			}
			if lo.File.MaxBackups > 0 {
				opts.File.MaxBackups = lo.File.MaxBackups
			}
			if lo.File.Compress {
				opts.File.Compress = true
			}
		}
		if lo.Kafka.Topic != "" && len(lo.Kafka.Brokers) > 0 {
			opts.Kafka.Topic = lo.Kafka.Topic
			opts.Kafka.Brokers = lo.Kafka.Brokers
		}
	}
}

// Level sets log level
// only trace/debug/info/warn/error/fatal could be used
func Level(l string) Option {
	return func(opts *Options) {
		if l == "" {
			l = _defaultLogLevel
		}
		opts.Level = l
	}
}

// StdoutOff switches stdout on or off
func StdoutOff(b bool) Option {
	return func(opts *Options) {
		opts.StdoutOff = b
	}
}

// FileDir file hook dir
func FileDir(dir string) Option {
	return func(opts *Options) {
		opts.File.Dir = dir
	}
}

// FileMaxSize max size per file
func FileMaxSize(size int) Option {
	return func(opts *Options) {
		opts.File.MaxSize = size
	}
}

// FileMaxBackups max backup files
func FileMaxBackups(num int) Option {
	return func(opts *Options) {
		opts.File.MaxBackups = num
	}
}

// FileMaxAge max file expire days
func FileMaxAge(days int) Option {
	return func(opts *Options) {
		opts.File.MaxAge = days
	}
}

// FileCompress judges whether to compress file
func FileCompress(compress bool) Option {
	return func(opts *Options) {
		opts.File.Compress = compress
	}
}

// KafkaTopic kafka topic
func KafkaTopic(topic string) Option {
	return func(opts *Options) {
		opts.Kafka.Topic = topic
	}
}

// KafkaBrokers kafka brokers, with ip:port each item
func KafkaBrokers(brokers []string) Option {
	return func(opts *Options) {
		opts.Kafka.Brokers = brokers
	}
}
