package zap_log_wrapper

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.uber.org/zap"
)

type loggerTestSuite struct {
	suite.Suite
}

func (l *loggerTestSuite) SetupTest() {
	err := NewLogger(&LoggerConfiguration{
		Level:            "debug",
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stdout", "./main.log"},
		ErrorOutputPaths: []string{"stdout", "./main.log"},
	})
	if err != nil {
		panic(err)
	}
}

func (l *loggerTestSuite) TestDebug() {
	Debug("this is message", zap.String("key", "value"))
	Debugf("this is message %s", "test")
	Debugw("this is message", zap.String("key", "value"))
	Info("this is message", zap.String("key", "value"))
	Infof("this is message %s", "test")
	Infow("this is message", zap.String("key", "value"))
	Warn("this is message", zap.String("key", "value"))
	Warnf("this is message %s", "test")
	Warnw("this is message", zap.String("key", "value"))
	w := zapcore.AddSync(&lumberjack.Logger{})
	zapcore.NewCore()
}

func TestResponseBuffer(t *testing.T) {
	suite.Run(t, &loggerTestSuite{})
}
