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
	ErrModelAPIKeyRequired         = errors.New("model API key is required")
	ErrUnsupportedProvider         = errors.New("unsupported model provider")
	ErrNoModelsForProvider         = errors.New("no models configured for provider")
	ErrInvalidModelResponse        = errors.New("model returned invalid candidate data")
	ErrIncompleteModelResponse     = errors.New("model did not finish generating candidates")
	ErrModelResponseTooLarge       = errors.New("model response exceeded the size limit")
)
