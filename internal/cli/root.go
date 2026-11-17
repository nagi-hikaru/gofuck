package cli

import (
	"github.com/nagi-hikaru/gofuck/internal/app"
	"github.com/nagi-hikaru/gofuck/internal/handler"
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
		RunE:          handler.RootHandler(fucker),
	}

	rules := &cobra.Command{
		Use:   "rules",
		Short: "List correction rules",
		Long:  "Show each correction rule, whether it is enabled, and its effective priority.",
		Args:  cobra.NoArgs,
		RunE:  handler.RulesHandler(fucker),
	}

	init := &cobra.Command{
		Use:       "init <shell>",
		Short:     "Print the shell integration script",
		Long:      "Print the shell integration script for Zsh or Bash. Load it in your shell startup file to enable command correction.",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"zsh", "bash"},
		RunE:      handler.InitHandler,
	}

	root.Flags().BoolP("yes", "y", false, "Run the first suggestion without confirmation")
	root.Flags().StringP("model", "m", "", "Send the previous command and stderr to a model (deepseek)")

	root.AddCommand(rules, init)

	return root
}
