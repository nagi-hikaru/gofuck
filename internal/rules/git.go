package rules

import (
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/model"
)

func GitRule() model.Rule {
	return model.Rule{
		Name:     "git",
		Enabled:  true,
		Priority: 1,
		Suggest:  gitFixer,
	}
}

func gitFixer(command model.Command) []string {
	if len(command.Words) < 2 || command.Words[0].Value != "git" {
		return nil
	}
	word := command.Words[1]
	if !strings.Contains(
		command.Output,
		"git: '"+word.Value+"' is not a git command",
	) {
		return nil
	}

	var candidates []string
	for _, line := range strings.Split(command.Output, "\n") {
		if !strings.HasPrefix(line, "\t") {
			continue
		}

		trimmed := strings.TrimSpace(line)
		candidates = append(candidates, command.Script[:word.Start]+trimmed+command.Script[word.End:])
	}
	return candidates
}
