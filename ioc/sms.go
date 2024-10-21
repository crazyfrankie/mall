package ioc

import (
	"mall/internal/user/service/sms"
	"mall/internal/user/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
