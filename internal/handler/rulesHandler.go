package handler

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/nagi-hikaru/gofuck/internal/app"
	"github.com/spf13/cobra"
)

// RulesHandler 返回输出规则列表的处理函数。
func RulesHandler(fucker app.Fucker) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		rows := fucker.RuleList()
		out := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 8, 2, ' ', 0)
		if _, err := fmt.Fprintf(out, "total %d\n", len(rows)); err != nil {
			return err
		}
		for _, row := range rows {
			if _, err := fmt.Fprintln(out, strings.Join(row, "\t")); err != nil {
				return err
			}
		}
		return out.Flush()
	}
}
