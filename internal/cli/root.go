package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/nagi-hikaru/gofuck/internal/app"
	"github.com/nagi-hikaru/gofuck/internal/status"
	"github.com/spf13/cobra"
)

// New 创建根命令
func New(fucker app.Fucker) *cobra.Command {
	root := &cobra.Command{
		Use:           "fuck",
		Short:         "Fix your previous shell command",
		Long:          "Suggest corrections for your previous shell command. Use -y to run the first suggestion without confirmation.",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			script := os.Getenv("GOFUCK_SCRIPT")
			if os.Getenv("GOFUCK_PHASE") == "check" {
				if _, ok := app.Parse(script); !ok {
					return status.ErrNoCandidate
				}
				return nil
			}

			// 读取 -y 或 --yes 的值
			yes, err := cmd.Flags().GetBool("yes")
			if err != nil {
				return err
			}
			return fucker.Fix(cmd.Context(), cmd.OutOrStdout(), script, os.Getenv("GOFUCK_OUTPUT_FILE"), yes)
		},
	}

	rules := &cobra.Command{
		Use:   "rules",
		Short: "List correction rules",
		Long:  "Show each correction rule, whether it is enabled, and its effective priority.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
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
		},
	}

	root.Flags().BoolP("yes", "y", false, "Run the first suggestion without confirmation")

	root.AddCommand(rules)

	return root
}
