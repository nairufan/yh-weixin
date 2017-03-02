package apperror

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type AppError struct {
	ErrorCode ErrorCode `json:"errorCode"`
	Message   string    `json:"message"`

	Target    string           `json:"target"`
	Details   []AppErrorDetail `json:"details"`

	cause     interface{}
	stack     []uintptr
	frames    []StackFrame
}

type AppErrorDetail struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

const MaxStackDepth = 100

func (err *AppError) CaptureStackTrace() *AppError {
	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(3, stack)

	err.stack = stack[:length]
	return err
}

func Wrap(cause interface{}) *AppError {
	if IsAppError(cause) {
		return cause.(*AppError)
	}

	appError := &AppError{
		ErrorCode: InternalServerError,
		Message:   "Wrap non AppError",
		cause:     cause,
	}
	appError.CaptureStackTrace()
	return appError
}

func WrapGeneric(cause interface{}) *AppError {
	if IsAppError(cause) {
		return cause.(*AppError)
	}

	appError := &AppError{
		ErrorCode: GenericWrap,
		Message:   "Wrap non AppError",
		cause:     cause,
	}
	appError.CaptureStackTrace()
	return appError
}

func IsAppError(err interface{}, args ...interface{}) bool {
	if err == nil {
		return false
	}

	switch e := err.(type) {
	case *AppError:
		if len(args) == 0 {
			return true
		}
		for _, c := range args {
			switch c.(type) {
			case ErrorCode:
				if e.ErrorCode == c.(ErrorCode) {
					return true
				}
			case string:
				if e.ErrorCode == ErrorCode(c.(string)) {
					return true
				}
			}
		}
		return false
	case error:
		return false
	default:
		return false
	}
}

func (err *AppError) StatusCode() int {
	if value, exists := errorCodesToStatusCode[err.ErrorCode]; exists {
		return value
	} else {
		return http.StatusInternalServerError
	}
}

func (err *AppError) Cause() interface{} {
	return err.cause
}

func (err *AppError) Error() string {
	result := "[" + string(err.ErrorCode) + "] " + err.Message

	if err.Target != "" {
		result += " (" + err.Target + ")"
	}

	return result
}

func (err *AppError) Stack() string {
	buf := bytes.Buffer{}

	for _, frame := range err.StackFrames() {
		buf.WriteString(frame.String())
	}

	return buf.String()
}

func (err *AppError) ErrorStack() string {
	return err.Error() + "\n" + err.Stack()
}

func (err *AppError) String() string {

	result := err.ErrorStack()

	if err.cause != nil {
		if IsAppError(err.cause) {
			result += strings.Replace(
				"\nCaused by: " + err.cause.(*AppError).String(),
				"\n", "\n  ", -1)
		} else {
			result += strings.Replace(
				fmt.Sprint("\nCaused by: ", err.cause),
				"\n", "\n  ", -1)
		}
	}

	return result
}

func (err *AppError) StackFrames() []StackFrame {
	if err.frames == nil {
		err.frames = make([]StackFrame, len(err.stack))

		for i, pc := range err.stack {
			err.frames[i] = *NewStackFrame(pc, i > 0 && err.frames[i - 1].InSigPanic)
		}
	}

	return err.frames
}
