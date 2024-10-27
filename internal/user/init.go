package user

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"mall/internal/user/service/sms"
	"mall/internal/user/service/sms/memory"
	"mall/pkg/logger"
	"mall/pkg/zapx"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}

func InitLogger() logger.Logger {
	encodeConfig := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encodeConfig), zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	customCore := zapx.NewCustomCore(core)
	l := zap.New(customCore)

	return logger.NewZapLogger(l)
}
