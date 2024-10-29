package web

import "fmt"

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

func GetResponse(options ...func(*Response)) Response {
	resp := Response{
		Status: 200, // 默认状态
		Msg:    "",  // 默认消息
	}
	for _, opt := range options {
		opt(&resp)
	}
	return resp
}

// 设置具体参数的函数

func WithStatus(status int) func(*Response) {
	return func(r *Response) {
		r.Status = status
	}
}

func WithData(data interface{}) func(*Response) {
	return func(r *Response) {
		r.Data = data
	}
}

func WithMsg(msg string) func(*Response) {
	return func(r *Response) {
		r.Msg = msg
	}
}

// BusinessError 定义错误类型
type BusinessError struct {
	Message string
	Err     error
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func NewBusinessError(message string, err error) *BusinessError {
	return &BusinessError{
		Message: message,
		Err:     err,
	}
}
