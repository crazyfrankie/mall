package ioc

import (
	"mall/service/sms"
	"mall/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
