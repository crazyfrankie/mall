package ioc

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"mall/pkg/zapx"
)

func InitLogger() *zap.Logger {
	encodeConfig := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encodeConfig), zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	customCore := zapx.NewCustomCore(core)
	logger := zap.New(customCore)

	return logger
}
