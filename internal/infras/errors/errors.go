package errors

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/xerrors"
)

func New(code ErrorCode, err error, msg string) *Error {
	return &Error{
		code: code,
		err:  err,
		msg:  msg,
	}
}

func Newf(code ErrorCode, err error, format string, args ...any) *Error {
	msg := fmt.Sprintf(format, args...)
	return New(code, err, msg)
}

func Wrapf(err error, format string, args ...any) *Error {
	msg := fmt.Sprintf(format, args...)
	return New(Code(err), err, msg)
}

func Code(err error) ErrorCode {
	if err == nil {
		return OK
	}
	var e *Error
	if errors.As(err, &e) {
		return e.code
	}
	if errors.Is(err, context.Canceled) {
		return Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return DeadlineExceeded
	}

	return Unknown
}

type Error struct {
	code ErrorCode
	err  error
	msg  string
}

func (e *Error) Error() string {
	return fmt.Sprint(e)
}

func (e *Error) FormatError(p xerrors.Printer) (next error) {
	if e.msg == "" {
		p.Printf("code=%v", e.code)
	} else {
		p.Printf("%s (code=%v)", e.msg, e.code)
	}

	return e.err
}

func (e *Error) Format(s fmt.State, c rune) {
	xerrors.FormatError(e, s, c) // important!
	// will call e.FormatError()
}
