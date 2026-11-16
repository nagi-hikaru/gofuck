package status

import (
	"errors"
)

var (
	ErrNoHistory                   = errors.New("no command to correct")
	ErrNoCandidate                 = errors.New("no correction found")
	ErrCanceled                    = errors.New("canceled")
	ErrInteractiveSelectionPending = errors.New("interactive selection is pending; use -y in this exercise")
	ErrDiagnosticOutputLimit       = errors.New("diagnostic output exceeded the limit")
	ErrTerminalRequired            = errors.New("interactive confirmation requires a terminal; use -y to accept the first suggestion")
)
