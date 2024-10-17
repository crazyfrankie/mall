package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"mall/repository"
	"mall/service/sms"
)

var (
	ErrSendTooMany   = repository.ErrSendTooMany
	ErrVerifyTooMany = repository.ErrVerifyTooMany
)

const (
	codeTplId = ""
	secretKey = "BgrTwHrRffd6LMXZWXGJCaKZHGb5p5h8"
)

type CodeService struct {
	repo *repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo *repository.CodeRepository, sms sms.Service) *CodeService {
	return &CodeService{
		repo: repo,
		sms:  sms,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.GenerateCode()

	// 加密后再存储
	hash := svc.GenerateHMAC(code, secretKey)

	err := svc.repo.Store(ctx, biz, phone, hash)
	if err != nil {
		return err
	}

	// 发送
	err = svc.sms.Send(ctx, codeTplId, []string{code}, phone)

	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	enCode := svc.GenerateHMAC(inputCode, secretKey)

	return svc.repo.Verify(ctx, biz, phone, enCode)
}

func (svc *CodeService) GenerateCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var code strings.Builder
	for i := 0; i < 6; i++ {
		digit := rand.Intn(10)
		code.WriteString(strconv.Itoa(digit))
	}
	return code.String()
}

func (svc *CodeService) GenerateHMAC(code, key string) string {
	// 创建一个新的 HMAC 哈希对象，使用 SHA-256 哈希算法，并以 key 作为密钥。
	h := hmac.New(sha256.New, []byte(key))

	// 将输入的 code 数据写入到 HMAC 哈希对象中，进行哈希计算。
	h.Write([]byte(code))

	// 计算哈希值并返回其十六进制表示形式。
	return hex.EncodeToString(h.Sum(nil))
}
