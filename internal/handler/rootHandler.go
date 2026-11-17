package handler

import (
	"fmt"
	"os"
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/app"
	"github.com/nagi-hikaru/gofuck/internal/status"
	"github.com/spf13/cobra"
)

// RootHandler 返回根命令的修正处理函数。
func RootHandler(fucker app.Fucker) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		script := os.Getenv("GOFUCK_SCRIPT")

		yes, err := cmd.Flags().GetBool("yes")
		if err != nil {
			return err
		}
		provider, err := cmd.Flags().GetString("model")
		if err != nil {
			return err
		}
		if strings.TrimSpace(script) == "" {
			return status.ErrNoHistory
		}

		ctx := cmd.Context()
		command := fucker.LoadCommand(script, os.Getenv("GOFUCK_OUTPUT_FILE"))
		var candidates []string
		if cmd.Flags().Changed("model") {
			candidates, err = fucker.AskAI(ctx, provider, command)
			if err != nil {
				return err
			}
		} else {
			candidates = fucker.Find(ctx, command)
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if len(candidates) == 0 {
			return status.ErrNoCandidate
		}

		// -y 直接采用首个候选，否则交给用户选择。
		selected := candidates[0]
		if !yes {
			selected, err = Select(ctx, os.Stdin, cmd.ErrOrStderr(), candidates)
			if err != nil {
				return err
			}
		}
		// 输出前再检查一次是否已取消。
		if err := ctx.Err(); err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), selected)
		return err
	}
}
