package rules

import (
	"slices"

	"github.com/nagi-hikaru/gofuck/internal/config"
	"github.com/nagi-hikaru/gofuck/internal/engine"
)

func clearRule(cfg config.RuleConfig) engine.Rule {
	return engine.Rule{
		Name:     "clear",
		Enabled:  cfg.Enabled,
		Priority: cfg.Priority,
		Match:    clearMatcher,
		Fix:      clearFixer,
	}
}

func clearMatcher(command engine.Command) bool {
	if len(command.Words) != 1 {
		return false
	}
	word := command.Words[0].Value
	if len(word) != len("clear") || word == "clear" {
		return false
	}
	letters := []byte(word)
	slices.Sort(letters)
	return string(letters) == "acelr"
}

func clearFixer(_ engine.Command) []string {
	return []string{"clear"}
}
