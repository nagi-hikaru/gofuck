package rules

import (
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/model"
)

func gitRule(cfg model.RuleConfig) model.Rule {
	return model.Rule{
		Name:     "git_not_command",
		Enabled:  cfg.Enabled,
		Priority: cfg.Priority,
		Match:    gitMatcher,
		Fix:      gitFixer,
	}
}

func gitMatcher(command model.Command) bool {
	if len(command.Words) < 2 || command.Words[0].Value != "git" {
		return false
	}
	if !strings.Contains(command.Output, "not a git command") {
		return false
	}
	return true
}

func gitFixer(command model.Command) []string {
	word := command.Words[1]
	var candidates []string
	for line := range strings.SplitSeq(command.Output, "\n") {
		if !strings.HasPrefix(line, "\t") {
			continue
		}

		trimmed := strings.TrimSpace(line)
		candidates = append(candidates, command.Script[:word.Start]+trimmed+command.Script[word.End:])
	}
	return candidates
}
