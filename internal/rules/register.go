package rules

import (
	"sort"

	"github.com/nagi-hikaru/gofuck/internal/config"
	"github.com/nagi-hikaru/gofuck/internal/model"
)

var Rules []model.Rule

// Register 根据已加载的配置注册规则
func Register() {
	// 注册规则
	Rules = []model.Rule{
		GitRule(ruleConfigOrDefault("git")),
	}

	// 按 优先级 -> 字母 排序
	sort.SliceStable(Rules, func(i, j int) bool {
		a, b := Rules[i], Rules[j]
		if a.Priority != b.Priority {
			return a.Priority < b.Priority
		}
		return a.Name < b.Name
	})
}

func ruleConfigOrDefault(name string) model.RuleConfig {
	if cfg, ok := config.Config.RulesConfig[name]; ok {
		return cfg
	}
	return model.RuleConfig{Enabled: true, Priority: 0}
}
