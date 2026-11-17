package handler

import (
	"fmt"

	"github.com/nagi-hikaru/gofuck/internal/shell"
	"github.com/spf13/cobra"
)

// InitHandler 输出指定 Shell 的集成脚本。
func InitHandler(cmd *cobra.Command, args []string) error {
	script, err := shell.Script(args[0])
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(cmd.OutOrStdout(), script)
	return err
}
