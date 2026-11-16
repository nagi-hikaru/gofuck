package app

import (
	"context"
	"io"
	"os"

	"atomicgo.dev/cursor"
	"github.com/nagi-hikaru/gofuck/internal/status"
	"github.com/pterm/pterm"
	"golang.org/x/term"
)

const maxVisibleCandidates = 3

func Select(ctx context.Context, in *os.File, out io.Writer, candidates []string) (string, error) {
	if len(candidates) == 0 {
		return "", status.ErrNoCandidate
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	output, ok := out.(*os.File)
	if !ok || !term.IsTerminal(int(in.Fd())) || !term.IsTerminal(int(output.Fd())) {
		return "", status.ErrTerminalRequired
	}

	fd := int(in.Fd())
	state, err := term.GetState(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, state)

	columns, height, _ := term.GetSize(int(output.Fd()))
	if columns <= 0 {
		columns = pterm.FallbackTerminalWidth
	}
	if height <= 0 {
		height = pterm.FallbackTerminalHeight
	}
	// PTerm 默认使用 stdout 确定宽度和绘制；stdout 留给选中的命令。
	previousArea, previousOutput := pterm.DefaultArea, pterm.GetDefaultOutput()
	pterm.DefaultArea.SetWriter(output)
	pterm.SetDefaultOutput(output)
	pterm.SetForcedTerminalSize(columns, height)
	cursor.SetTarget(output)
	defer func() {
		pterm.DefaultArea = previousArea
		pterm.SetDefaultOutput(previousOutput)
		pterm.SetForcedTerminalSize(0, 0)
		cursor.SetTarget(os.Stdout)
	}()

	options, originals := selectOptions(candidates, columns)
	canceled := false
	printer := pterm.DefaultInteractiveSelect.
		WithOptions(options).
		WithMaxHeight(min(maxVisibleCandidates, max(1, height-2))).
		WithFilter(false).
		WithOnInterruptFunc(func() { canceled = true })

	const promptLabel = "Fuck which"
	prompt := fitLine(promptLabel, max(1, columns-1))
	if columns > len(promptLabel+" [enter/ctrl+c]") {
		prompt = promptLabel + " [" + pterm.FgGreen.Sprint("enter") + "/" + pterm.FgRed.Sprint("ctrl+c") + "]"
	}
	selected, err := printer.Show(prompt)
	if err != nil {
		return "", err
	}
	if err = ctx.Err(); err != nil {
		return "", err
	}
	if canceled {
		return "", status.ErrCanceled
	}
	command, ok := originals[selected]
	if !ok {
		return "", status.ErrNoCandidate
	}
	return command, nil
}
