package errors

import (
	"fmt"
	"strings"
)

type HttpErr struct {
	Err  error
	Code int
}

func (h *HttpErr) Error() string {
	return fmt.Sprintf("%s", h.Err)
}

func NewErr(code int, err error) error {
	return &HttpErr{
		Err:  err,
		Code: code,
	}
}

type ValidationErr struct {
	Errs map[string]string
	Code int
	Msg  string
}

func (v *ValidationErr) Error() string {
	var builder strings.Builder
	builder.WriteString(v.Msg)
	builder.WriteByte('\n')
	for field, err := range v.Errs {
		builder.WriteString(fmt.Sprintf("%s: %s\n", field, err))
	}
	return builder.String()
}

func NewValidationErr(msg string, errs map[string]string, code int) error {
	return &ValidationErr{
		Errs: errs,
		Code: code,
		Msg:  msg,
	}
}
