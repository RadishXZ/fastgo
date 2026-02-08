package errorsx

import (
	"errors"
	"fmt"
)

type ErrorX struct {
	// Code 表示错误的 HTTP 状态码, 用于与客户端进行交互时标识错误的类型
	Code int `json:"code,omitempty"`

	// Reason 表示错误发生的原因, 通常为业务错误码, 用于精准定位问题
	Reason string `json:"reason,omitempty"`

	// Message 表示简短的错误信息, 通常可直接暴露给用户查看
	Message string `json:"message,omitempty"`
}

// New 创建一个新的错误
func New(code int, reason string, format string, args ...any) *ErrorX {
	return &ErrorX{
		Code: code,
		Reason: reason,
		Message: fmt.Sprintf(format, args...),
	}
}

// Error 实现 error 接口中的 `Error` 方法
func (err *ErrorX) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s", err.Code, err.Reason, err.Message)
}

// WithMessage 设置错误的 Message 字段
func (err *ErrorX) WithMessage(format string, args ...any) *ErrorX {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

func FormError(err error) *ErrorX {
	if err == nil {
		return nil
	}

	if errx := new(ErrorX); errors.As(err, &errx) {
		return errx
	}

	return New(ErrInternal.Code, ErrInternal.Reason, err.Error())
}