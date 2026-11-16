package rules

import (
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/model"
)

func LimaRule(config model.RuleConfig) model.Rule {
	return model.Rule{
		Name:     "lima",
		Enabled:  config.Enabled,
		Priority: config.Priority,
		Suggest:  LimaFixer,
	}
}

func LimaFixer(command model.Command) []string {
	if len(command.Words) < 2 || command.Words[0].Value != "limactl" {
		return nil
	}
	if !strings.Contains(command.Output, "unknown command") {
		return nil
	}

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
