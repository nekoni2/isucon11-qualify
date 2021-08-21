package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func initLogger() {
	conf := zap.NewProductionConfig()
	// we might use third-party log analyzer to optimize
	conf.OutputPaths = []string{"/home/isucon/webapp/go/logs/error.log"}
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.Encoding = "console"
	if os.Getenv("ENV") == "dev" {
		conf.OutputPaths = []string{"stderr"}
	}
	// disable info logs when production
	conf.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if os.Getenv("ENV") == "prod" {
		conf.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	}
	l, _ := conf.Build()
	logger = l.Sugar()
}
