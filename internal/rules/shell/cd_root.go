package shell

import (
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/config"
	"github.com/nagi-hikaru/gofuck/internal/engine"
)

func CdRootRule(cfg config.RuleConfig) engine.Rule {
	return engine.Rule{
		Name:     "cd_root",
		Enabled:  cfg.Enabled,
		Priority: cfg.Priority,
		Match:    cdRootMatcher,
		Fix:      cdRootFixer,
	}
}

func cdRootMatcher(command engine.Command) bool {
	return strings.TrimSpace(command.Script) == "cd~"
}

func cdRootFixer(_ engine.Command) []string {
	return []string{"cd ~"}
}
