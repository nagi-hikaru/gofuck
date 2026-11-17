package app

import (
	"context"
	"sort"
	"strconv"

	"github.com/nagi-hikaru/gofuck/internal/model"
	"github.com/nagi-hikaru/gofuck/internal/rules"
)

type Fucker interface {
	LoadCommand(script, outputFile string) model.Command
	Find(ctx context.Context, command model.Command) []string
	AskAI(ctx context.Context, command model.Command) []string
	RuleList() [][]string
}

type FuckerImpl struct {
	rules []model.Rule
}

func NewFucker() Fucker {
	return &FuckerImpl{
		rules: rules.Rules,
	}
}

// LoadCommand 解析命令并读取 shell 已采集的诊断。
func (*FuckerImpl) LoadCommand(script, outputFile string) model.Command {
	words, _ := Parse(script)
	output, _ := readOutput(outputFile)

	return model.Command{
		Script: script,
		Words:  words,
		Output: output,
	}
}

// Find 运行已注册的规则，再对候选排序和去重。
func (f *FuckerImpl) Find(ctx context.Context, command model.Command) []string {
	var candidates []model.Candidate
	for _, rule := range f.rules {
		if ctx.Err() != nil {
			return nil
		}

		if !rule.Enabled {
			continue
		}

		if !rule.Match(command) {
			continue
		}

		for index, script := range rule.Fix(command) {
			candidates = append(candidates, model.Candidate{
				Script: script,
				Rule:   rule.Name,
				Weight: rule.Priority * (index + 1),
			})
		}
	}

	// 按 weight 排序
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Weight < candidates[j].Weight
	})

	// 去重
	seen := make(map[string]bool, len(candidates))
	scripts := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if !seen[candidate.Script] {
			seen[candidate.Script] = true
			scripts = append(scripts, candidate.Script)
		}
	}
	return scripts
}

func (f *FuckerImpl) AskAI(ctx context.Context, command model.Command) []string {
	return []string{}
}

// RuleList 返回各条规则的名称、状态和优先级。
func (f *FuckerImpl) RuleList() [][]string {
	rows := make([][]string, 0, len(f.rules))
	for _, rule := range f.rules {
		status := "disabled"
		if rule.Enabled {
			status = "enabled"
		}
		rows = append(rows, []string{rule.Name, status, strconv.Itoa(rule.Priority)})
	}
	return rows
}
