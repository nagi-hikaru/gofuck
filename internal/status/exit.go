package status

import (
	"context"
	"errors"
)

const (
	ExitOK           = 0
	ExitNoCorrection = 1
	ExitError        = 2
	ExitCanceled     = 130
)

// ExitCode 将命令执行结果转换为进程退出码
func ExitCode(err error) int {
	switch {
	case err == nil:
		return ExitOK
	case errors.Is(err, ErrCanceled), errors.Is(err, context.Canceled):
		return ExitCanceled
	case errors.Is(err, ErrNoHistory), errors.Is(err, ErrNoCandidate):
		return ExitNoCorrection
	default:
		return ExitError
	}
}
