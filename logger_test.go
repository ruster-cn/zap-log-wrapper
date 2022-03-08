package zap_log_wrapper

import (
	"github.com/stretchr/testify/suite"
	"testing"

	"go.uber.org/zap"
)

type loggerTestSuite struct {
	suite.Suite
}

func (l *loggerTestSuite) SetupTest() {

}

func (l *loggerTestSuite) TestStdLog() {
	err := NewLogger(&LoggerConfiguration{
		Level:       "debug",
		Development: false,
		Encoding:    "text",
		OutputPath:  StdOut,
	})
	if err != nil {
		panic(err)
	}
	Debug("this is message", zap.String("key", "value"))
	Debugf("this is message %s", "test")
	Debugw("this is message", zap.String("key", "value"))
	Info("this is message", zap.String("key", "value"))
	Infof("this is message %s", "test")
	Infow("this is message", zap.String("key", "value"))
	Warn("this is message", zap.String("key", "value"))
	Warnf("this is message %s", "test")
	Warnw("this is message", zap.String("key", "value"))
}

func (l *loggerTestSuite) TestFileLog() {
	err := NewLogger(&LoggerConfiguration{
		Level:       "debug",
		Development: false,
		Encoding:    "text",
		OutputPath:  "std.log",
		Rotate: LogRotate{
			MaxSizeMB:  1,
			MaxBackups: 3,
			Compress:   true,
		},
	})
	if err != nil {
		panic(err)
	}

	Debug("this is message", zap.String("key", "value"))
	Debugf("this is message %s", "test")
	Debugw("this is message", zap.String("key", "value"))
	Info("this is message", zap.String("key", "value"))
	Infof("this is message %s", "test")
	Infow("this is message", zap.String("key", "value"))
	Warn("this is message", zap.String("key", "value"))
	Warnf("this is message %s", "test")
	Warnw("this is message", zap.String("key", "value"))
	Errorf("this is message", zap.String("key", "value"))

}

func TestResponseBuffer(t *testing.T) {
	suite.Run(t, &loggerTestSuite{})
}
