package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/model"
	"github.com/nagi-hikaru/gofuck/internal/rules"
	"github.com/nagi-hikaru/gofuck/internal/status"
)

type Fucker interface {
	RuleList() [][]string
	Fix(ctx context.Context, out io.Writer, script, outputFile string, yes bool) error
}

type FuckerImpl struct {
	rules []model.Rule
}

func NewFucker() Fucker {
	return &FuckerImpl{
		rules: rules.Rules,
	}
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

// Fix 查找、选择并输出修正后的命令
func (f *FuckerImpl) Fix(ctx context.Context, out io.Writer, script, outputFile string, yes bool) error {
	if strings.TrimSpace(script) == "" {
		return status.ErrNoHistory
	}

	var err error

	command := loadCommand(script, outputFile)
	candidates := find(ctx, command, f.rules)
	if err = ctx.Err(); err != nil {
		return err
	}
	if len(candidates) == 0 {
		return status.ErrNoCandidate
	}

	// -y 直接采用首个候选，否则交给用户选择
	selected := candidates[0]
	if !yes {
		selected, err = Select(ctx, os.Stdin, os.Stderr, candidates)
		if err != nil {
			return err
		}
	}
	// 输出前再次检查是否已取消
	if err = ctx.Err(); err != nil {
		return err
	}
	// 将修正命令写到 out
	_, err = fmt.Fprintln(out, selected)
	return err
}

// loadCommand 解析命令并读取 shell 已采集的诊断
func loadCommand(script, outputFile string) model.Command {
	words, _ := Parse(script)
	output, _ := readOutput(outputFile)

	return model.Command{
		Script: script,
		Words:  words,
		Output: output,
	}
}

// find 运行传入的规则，再对候选排序和去重
func find(ctx context.Context, command model.Command, ruleList []model.Rule) []string {
	var candidates []model.Candidate
	for _, rule := range ruleList {
		if ctx.Err() != nil {
			return nil
		}
		for index, script := range rule.Suggest(command) {
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
