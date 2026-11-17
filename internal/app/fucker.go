package app

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/ai"
	"github.com/nagi-hikaru/gofuck/internal/config"
	"github.com/nagi-hikaru/gofuck/internal/engine"
	"github.com/nagi-hikaru/gofuck/internal/rules"
	"github.com/nagi-hikaru/gofuck/internal/status"
)

type Fucker interface {
	LoadCommand(script, outputFile string) engine.Command
	Find(ctx context.Context, command engine.Command) []string
	AskAI(ctx context.Context, provider string, command engine.Command) ([]string, error)
	RuleList() [][]string
}

type FuckerImpl struct {
	rules []engine.Rule
}

func NewFucker() Fucker {
	return &FuckerImpl{
		rules: rules.Rules,
	}
}

// LoadCommand 解析命令并读取 shell 已采集的诊断。
func (*FuckerImpl) LoadCommand(script, outputFile string) engine.Command {
	words, _ := Parse(script)
	output, _ := readOutput(outputFile)

	return engine.Command{
		Script: script,
		Words:  words,
		Output: output,
	}
}

// Find 运行已注册的规则，再对候选排序和去重。
func (f *FuckerImpl) Find(ctx context.Context, command engine.Command) []string {
	// 保留本地规则原有的匹配行为，给模型的诊断则保留原文。
	command.Output = strings.ToLower(command.Output)
	var candidates []engine.Candidate
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
			candidates = append(candidates, engine.Candidate{
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

// AskAI 向指定模型请求修正候选。
func (*FuckerImpl) AskAI(ctx context.Context, provider string, command engine.Command) ([]string, error) {
	cfg, ok := config.Config.ProvidersConfig[provider]
	if !ok {
		return nil, status.ErrUnsupportedProvider
	}
	if len(cfg.Models) == 0 {
		return nil, status.ErrNoModelsForProvider
	}

	candidates, err := ai.Suggest(ctx, cfg, cfg.Models[0], ai.Input{
		Script: command.Script,
		Output: command.Output,
	})
	if err != nil {
		return nil, err
	}

	scripts := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		scripts = append(scripts, candidate.Script)
	}
	return scripts, nil
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
