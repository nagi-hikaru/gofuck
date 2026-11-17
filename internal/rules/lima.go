package rules

import (
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/config"
	"github.com/nagi-hikaru/gofuck/internal/engine"
)

func limaRule(cfg config.RuleConfig) engine.Rule {
	return engine.Rule{
		Name:     "lima",
		Enabled:  cfg.Enabled,
		Priority: cfg.Priority,
		Match:    limaMatcher,
		Fix:      limaFixer,
	}
}

func limaMatcher(command engine.Command) bool {
	if len(command.Words) < 2 || command.Words[0].Value != "limactl" {
		return false
	}
	if !strings.Contains(command.Output, "unknown command") {
		return false
	}
	return true
}

func limaFixer(command engine.Command) []string {
	word := command.Words[1]
	var candidates []string
	for line := range strings.SplitSeq(command.Output, `\n`) {
		suggestion, ok := strings.CutPrefix(line, `\t`)
		if !ok {
			continue
		}

		trimmed := strings.TrimSpace(suggestion)
		candidates = append(candidates, command.Script[:word.Start]+trimmed+command.Script[word.End:])
	}
	return candidates
}
