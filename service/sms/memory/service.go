package memory

import (
	"context"
	"fmt"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (svc *Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}
