package shell

import (
	"os"
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/config"
	"github.com/nagi-hikaru/gofuck/internal/engine"
)

func RmRFRule(cfg config.RuleConfig) engine.Rule {
	return engine.Rule{
		Name:     "rm_rf",
		Enabled:  cfg.Enabled,
		Priority: cfg.Priority,
		Match:    rmrfMatcher,
		Fix:      rmrfFixer,
	}
}

func rmrfMatcher(command engine.Command) bool {
	if len(command.Words) < 2 || command.Words[0].Value != "rm" {
		return false
	}
	if !strings.Contains(command.Output, "is a directory") {
		return false
	}
	info, err := os.Lstat(command.Words[1].Value)
	return err == nil && info.IsDir()
}

func rmrfFixer(command engine.Command) []string {
	path := command.Words[1]
	rawPath := command.Script[path.Start:path.End]
	return []string{"rm -rf -- " + rawPath}
}
