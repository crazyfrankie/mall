package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type MiddlewareBuilder struct {
	allowReqBody   bool
	allowRespBody  bool
	allowUrlLength int
	allowReqLength int
	loggerFunc     func(ctx context.Context, al *AccessLog)
}

type AccessLog struct {
	// HTTP 方法
	Method string
	// URL 整个请求 URL
	URL      string
	ReqBody  string
	RespBody string
	Duration string
	Status   int
}

func NewMiddlewareBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc: fn,
	}
}

func (b *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	b.allowReqBody = true
	return b
}

func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}

func (b *MiddlewareBuilder) AllowUrlLength(length int) *MiddlewareBuilder {
	b.allowUrlLength = length
	return b
}

func (b *MiddlewareBuilder) AllowReqBodyLength(length int) *MiddlewareBuilder {
	b.allowReqLength = length
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		url := c.Request.URL.String()
		if len(url) > b.allowUrlLength {
			url = url[:b.allowUrlLength]
		}
		al := &AccessLog{
			Method: c.Request.Method,
			URL:    url,
		}

		if b.allowReqBody && c.Request.Body != nil {
			body, _ := c.GetRawData()
			// body 读完就没有了，需要放回去
			c.Request.Body = io.NopCloser(bytes.NewReader(body))

			if len(body) > b.allowReqLength {
				body = body[:b.allowReqLength]
			}
			// 很消耗 CPU 和内存的操作，因为会引起复制
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			c.Writer = &responseWriter{
				ResponseWriter: c.Writer,
				al:             al,
			}
		}

		defer func() {
			al.Duration = time.Since(start).String()
			b.loggerFunc(c, al)
		}()

		// 执行到业务逻辑
		c.Next()
	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) WriteString(data string) (int, error) {
	w.al.RespBody = data
	return w.ResponseWriter.WriteString(data)
}
