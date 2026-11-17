package rules

import (
	"os"
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/model"
)

func rmrfRule(cfg model.RuleConfig) model.Rule {
	return model.Rule{
		Name:     "rm_rf",
		Enabled:  cfg.Enabled,
		Priority: cfg.Priority,
		Match:    rmrfMatcher,
		Fix:      rmrfFixer,
	}
}

func rmrfMatcher(command model.Command) bool {
	if len(command.Words) < 2 || command.Words[0].Value != "rm" {
		return false
	}
	if !strings.Contains(command.Output, "is a directory") {
		return false
	}
	info, err := os.Lstat(command.Words[1].Value)
	return err == nil && info.IsDir()
}

func rmrfFixer(command model.Command) []string {
	path := command.Words[1]
	rawPath := command.Script[path.Start:path.End]
	return []string{"rm -rf -- " + rawPath}
}
