package zap_log_wrapper

import (
	"errors"

	"go.uber.org/zap/zapcore"
)

type LoggerConfiguration struct {
	Level            Level    `yaml:"level"`
	Development      bool     `yaml:"development"`
	Encoding         Encoding `yaml:"encoding"`
	OutputPaths      []string `yaml:"output_paths"`
	ErrorOutputPaths []string `yaml:"error_output_paths"`
}

func (config *LoggerConfiguration) Default() *LoggerConfiguration {
	if config.Level == "" {
		config.Level = "debug"
	}
	if config.Encoding == "" {
		config.Encoding = TextEncoding
	}
	if len(config.OutputPaths) == 0 {
		config.OutputPaths = []string{"stdout"}
	}
	if len(config.ErrorOutputPaths) == 0 {
		config.ErrorOutputPaths = []string{"stdout"}
	}
	return config
}

func (config *LoggerConfiguration) Validate() error {
	if _, ok := logLevelMap[config.Level]; !ok {
		return errors.New("log level set error")
	}
	if config.Encoding != TextEncoding && config.Encoding != JsonEncoding {
		return errors.New("log encoding format set error")
	}
	if len(config.OutputPaths) == 0 {
		return errors.New("must set log output")
	}
	if len(config.ErrorOutputPaths) == 0 {
		return errors.New("must set log err output")
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
