package zapx

import "go.uber.org/zap/zapcore"

type CustomCore struct {
	zapcore.Core
}

func NewCustomCore(core zapcore.Core) *CustomCore {
	return &CustomCore{
		Core: core,
	}
}

func (z CustomCore) Write(en zapcore.Entry, fields []zapcore.Field) error {
	modifiedFields := make([]zapcore.Field, len(fields))
	copy(modifiedFields, fields)

	for i, fd := range modifiedFields {
		if fd.Key == "phone" {
			phone := fd.String
			modifiedFields[i].String = phone[:3] + "****" + phone[7:]
		}
	}

	return z.Core.Write(en, modifiedFields)
}
