package zap_log_wrapper

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"sort"
	"time"
)

var (
	errNoEncoderNameSpecified = errors.New("no encoder name specified")

	_encoderNameToConstructor = map[string]func(zapcore.EncoderConfig) (zapcore.Encoder, error){
		"console": func(encoderConfig zapcore.EncoderConfig) (zapcore.Encoder, error) {
			return zapcore.NewConsoleEncoder(encoderConfig), nil
		},
		"json": func(encoderConfig zapcore.EncoderConfig) (zapcore.Encoder, error) {
			return zapcore.NewJSONEncoder(encoderConfig), nil
		},
	}
)

func newEncoder(name string, encoderConfig zapcore.EncoderConfig) (zapcore.Encoder, error) {
	if encoderConfig.TimeKey != "" && encoderConfig.EncodeTime == nil {
		return nil, fmt.Errorf("missing EncodeTime in EncoderConfig")
	}

	if name == "" {
		return nil, errNoEncoderNameSpecified
	}
	constructor, ok := _encoderNameToConstructor[name]
	if !ok {
		return nil, fmt.Errorf("no encoder registered for name %q", name)
	}
	return constructor(encoderConfig)
}

func BuildEncoder(encoding string, encoderConfig zapcore.EncoderConfig) (zapcore.Encoder, error) {
	return newEncoder(encoding, encoderConfig)
}

// Build constructs a logger from the Config and Options.
func Build(cfg zap.Config, rotate LogRotate, opts ...zap.Option) (*zap.Logger, error) {
	enc, err := BuildEncoder(cfg.Encoding, cfg.EncoderConfig)
	if err != nil {
		return nil, err
	}
	if cfg.Level == (zap.AtomicLevel{}) {
		return nil, fmt.Errorf("missing Level")
	}

	var output zapcore.WriteSyncer
	var outErr error

	if cfg.OutputPaths[0] == StdOut {
		output, _, outErr = zap.Open(StdOut)
		if outErr != nil {
			return nil, outErr
		}
	} else {
		output = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.OutputPaths[0],
			MaxSize:    rotate.MaxSizeMB,
			MaxBackups: rotate.MaxBackups,
			Compress:   rotate.Compress,
			LocalTime:  true,
		})
	}

	log := zap.New(zapcore.NewCore(enc, output, cfg.Level), BuildOptions(cfg)...)

	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}
	return log, nil
}

func BuildOptions(cfg zap.Config) []zap.Option {
	var opts []zap.Option

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zapcore.ErrorLevel
	if cfg.Development {
		stackLevel = zapcore.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if scfg := cfg.Sampling; scfg != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			var samplerOpts []zapcore.SamplerOption
			if scfg.Hook != nil {
				samplerOpts = append(samplerOpts, zapcore.SamplerHook(scfg.Hook))
			}
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				cfg.Sampling.Initial,
				cfg.Sampling.Thereafter,
				samplerOpts...,
			)
		}))
	}

	if len(cfg.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(cfg.InitialFields))
		keys := make([]string, 0, len(cfg.InitialFields))
		for k := range cfg.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, cfg.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	return opts
}
