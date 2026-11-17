package rules

import "github.com/nagi-hikaru/gofuck/internal/model"

func cdParentRule(cfg model.RuleConfig) model.Rule {
	return model.Rule{
		Name:     "cd_parent",
		Enabled:  cfg.Enabled,
		Priority: cfg.Priority,
		Match:    cdParentMatcher,
		Fix:      cdParentFixer,
	}
}

func cdParentMatcher(command model.Command) bool {
	return len(command.Words) == 1 && command.Words[0].Value == "cd.."
}

func cdParentFixer(_ model.Command) []string {
	return []string{"cd .."}
}
