package types

import (
	"fmt"
	"strings"
)

func DBErr(err error) *Error {
	if err != nil {
		return &Error{
			Code:    500,
			Message: err.Error(),
		}
	}
	return nil
}

func NewErr(code int, format string, args ...interface{}) *Error {
	if len(args) < 1 {
		return &Error{Code: code, Message: format}
	}
	if strings.Count(format, "%") == len(args) {
		return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
	}

	msg := append([]interface{}{format}, args...)
	return &Error{Code: code, Message: fmt.Sprintln(msg...)}
}

type Error struct {
	Index   int
	Code    int
	Message string
}

func (e Error) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("%d, %s", e.Code, e.Message)
	}
	return e.Message
}
