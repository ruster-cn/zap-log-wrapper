package zap_log_wrapper

import (
	"errors"

	"go.uber.org/zap/zapcore"
)

const StdOut = "stdout"

type LoggerConfiguration struct {
	Level       Level     `yaml:"level"`
	Development bool      `yaml:"development"`
	Encoding    Encoding  `yaml:"encoding"`
	OutputPath  string    `yaml:"output_path"`
	Rotate      LogRotate `yaml:"rotate"`
}

type LogRotate struct {
	MaxSizeMB  int  `yaml:"max_size_mb"`
	MaxBackups int  `yaml:"max_backups"`
	Compress   bool `yaml:"compress"`
}

func (config *LoggerConfiguration) Default() {
	if config.Level == "" {
		config.Level = "debug"
	}
	if config.Encoding == "" {
		config.Encoding = TextEncoding
	}
	if config.OutputPath == "" {
		config.OutputPath = StdOut
	}

	if config.Rotate.MaxSizeMB == 0 {
		config.Rotate.MaxSizeMB = 1024
	}

	if config.Rotate.MaxBackups == 0 {
		config.Rotate.MaxBackups = 3
		config.Rotate.Compress = true
	}

}

func (config *LoggerConfiguration) Validate() error {
	config.Default()
	if _, ok := logLevelMap[config.Level]; !ok {
		return errors.New("log level set error")
	}
	if config.Encoding != TextEncoding && config.Encoding != JsonEncoding {
		return errors.New("log encoding format set error")
	}

	if config.Rotate.MaxSizeMB <= 0 {
		return errors.New("log rotate maxsize must greater 0")
	}

	if config.Rotate.MaxBackups <= 0 {
		return errors.New("log rotate max backups must greater 0")
	}

	return nil
}

var logLevelMap = map[Level]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	FatalLevel: zapcore.FatalLevel,
}

type Level string

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

type Encoding string

const (
	TextEncoding = "text"
	JsonEncoding = "json"
)

var logEncodingMap = map[Encoding]string{
	TextEncoding: "console",
	JsonEncoding: "json",
}
