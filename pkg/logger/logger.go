package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

// NewLogger - return an instance of logger
// "debug,info,warn,error,dpanic,panic,fatal"
func NewLogger(level string) *Logger {
	var lg *Logger
	atomicLevel := zap.NewAtomicLevel()
	err := atomicLevel.UnmarshalText([]byte(level))
	if err != nil {
		panic(err)
	}

	cfg := zap.Config{
		Encoding:         "json",
		Level:            atomicLevel,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			TimeKey:     "ts",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	lg = &Logger{logger}

	return lg
}
