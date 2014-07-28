package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qiniu/log"
)

const prefix = " ==> "

// --------------------------------------------------------------------

func New(msg string) error {
	return errors.New(msg)
}

// --------------------------------------------------------------------

type errorDetailer interface {
	ErrorDetail() string
}

func Detail(err error) string {
	if e, ok := err.(errorDetailer); ok {
		return e.ErrorDetail()
	}
	return prefix + err.Error()
}

// --------------------------------------------------------------------

type ErrorInfo struct {
	Err error
	Why error
	Cmd []interface{}
}

func Info(err error, cmd ...interface{}) *ErrorInfo {
	if e, ok := err.(*ErrorInfo); ok {
		err = e.Err
	}
	return &ErrorInfo{Cmd: cmd, Err: err}
}

func (r *ErrorInfo) Error() string {
	return r.Err.Error()
}

func (r *ErrorInfo) ErrorDetail() string {
	e := prefix + r.Err.Error() + " ~ " + fmt.Sprintln(r.Cmd...)
	if r.Why != nil {
		e += Detail(r.Why)
	}
	return e
}

func (r *ErrorInfo) Detail(err error) *ErrorInfo {
	r.Why = err
	return r
}

func (r *ErrorInfo) Method() (cmd string, ok bool) {
	if len(r.Cmd) > 0 {
		if cmd, ok = r.Cmd[0].(string); ok {
			if pos := strings.Index(cmd, " "); pos > 1 {
				cmd = cmd[:pos]
			}
		}
	}
	return
}

func (r *ErrorInfo) LogMessage() string {
	detail := r.ErrorDetail()
	if cmd, ok := r.Method(); ok {
		detail = cmd + " failed:\n" + detail
	}
	return detail
}

// deprecated. please use (*ErrorInfo).LogWarn
//
func (r *ErrorInfo) Warn() *ErrorInfo {
	log.Std.Output("", log.Lwarn, 2, r.LogMessage())
	return r
}

func (r *ErrorInfo) LogWarn(reqId string) *ErrorInfo {
	log.Std.Output(reqId, log.Lwarn, 2, r.LogMessage())
	return r
}

func (r *ErrorInfo) LogError(reqId string) *ErrorInfo {
	log.Std.Output(reqId, log.Lerror, 2, r.LogMessage())
	return r
}

func (r *ErrorInfo) Log(level int, reqId string) *ErrorInfo {
	log.Std.Output(reqId, level, 2, r.LogMessage())
	return r
}

// --------------------------------------------------------------------

func Err(err error) error {
	for {
		if e, ok := err.(*ErrorInfo); ok {
			err = e.Err
		} else {
			break
		}
	}
	return err
}

// --------------------------------------------------------------------
