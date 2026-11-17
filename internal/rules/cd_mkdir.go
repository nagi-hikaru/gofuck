package rules

import (
	"os"
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/model"
)

func cdMkdirRule(cfg model.RuleConfig) model.Rule {
	return model.Rule{
		Name:     "cd_mkdir",
		Enabled:  cfg.Enabled,
		Priority: cfg.Priority,
		Match:    cdMkdirMatcher,
		Fix:      cdMkdirFixer,
	}
}

func cdMkdirMatcher(command model.Command) bool {
	if len(command.Words) < 2 || command.Words[0].Value != "cd" {
		return false
	}
	if !strings.Contains(command.Output, "no such file or directory") {
		return false
	}
	_, err := os.Lstat(command.Words[1].Value)
	return os.IsNotExist(err)
}

func cdMkdirFixer(command model.Command) []string {
	path := command.Words[1]
	rawPath := command.Script[path.Start:path.End]
	return []string{"mkdir -p -- " + rawPath + " && " + command.Script}
}
