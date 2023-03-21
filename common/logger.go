package common

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

// Return a new custom zap logger instance
func NewLogger(env string, logPath string) (*Logger, error) {
	if env == "test" {
		return &Logger{zap.NewNop()}, nil
	}

	var (
		l   *zap.Logger
		err error
	)

	if strings.EqualFold(env, "production") {
		l, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	} else {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		l, err = config.Build()
		if err != nil {
			return nil, err
		}
	}

	logFile, err := os.Create(logPath)
	if err != nil {
		return nil, err
	}

	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zap.DebugLevel)

	l = l.Named(SvcName).
		WithOptions(zap.AddStacktrace(zapcore.PanicLevel)).
		WithOptions(zap.WithFatalHook(zapcore.WriteThenGoexit)).
		WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, fileCore)
		}))

	return &Logger{
		l,
	}, nil
}
