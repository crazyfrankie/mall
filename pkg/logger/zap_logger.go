package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(l *zap.Logger) Logger {
	return &ZapLogger{
		logger: l,
	}
}

func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.logger.Debug(msg, z.toZapLogger(args)...)
}

func (z *ZapLogger) Info(msg string, args ...Field) {
	z.logger.Info(msg, z.toZapLogger(args)...)
}

func (z *ZapLogger) Error(msg string, args ...Field) {
	z.logger.Error(msg, z.toZapLogger(args)...)
}

func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.logger.Warn(msg, z.toZapLogger(args)...)
}

func (z *ZapLogger) toZapLogger(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, ag := range args {
		res = append(res, zap.Any(ag.Key, ag.Val))
	}

	return res
}
