package rules

import (
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/config"
	"github.com/nagi-hikaru/gofuck/internal/engine"
)

func cdParentRule(cfg config.RuleConfig) engine.Rule {
	return engine.Rule{
		Name:     "cd_parent",
		Enabled:  cfg.Enabled,
		Priority: cfg.Priority,
		Match:    cdParentMatcher,
		Fix:      cdParentFixer,
	}
}

func cdParentMatcher(command engine.Command) bool {
	return strings.TrimSpace(command.Script) == "cd.."
}

func cdParentFixer(_ engine.Command) []string {
	return []string{"cd .."}
}
