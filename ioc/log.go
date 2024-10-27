package ioc

import (
	"mall/pkg/logger"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"mall/pkg/zapx"
)

func InitLogger() logger.Logger {
	encodeConfig := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encodeConfig), zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	customCore := zapx.NewCustomCore(core)
	l := zap.New(customCore)

	return logger.NewZapLogger(l)
}
