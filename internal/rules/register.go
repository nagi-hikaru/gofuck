package rules

import (
	"sort"

	"github.com/nagi-hikaru/gofuck/internal/config"
	"github.com/nagi-hikaru/gofuck/internal/engine"
)

var Rules []engine.Rule

// Register 根据已加载的配置注册规则
func Register() {
	// 注册规则
	Rules = []engine.Rule{
		gitRule(ruleConfigOrDefault("git_not_command")),
		limaRule(ruleConfigOrDefault("lima")),
		cdMkdirRule(ruleConfigOrDefault("cd_mkdir")),
		cdParentRule(ruleConfigOrDefault("cd_parent")),
		cdRootRule(ruleConfigOrDefault("cd_root")),
		rmrfRule(ruleConfigOrDefault("rm_rf")),
		clearRule(ruleConfigOrDefault("clear")),
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

func ruleConfigOrDefault(name string) config.RuleConfig {
	if cfg, ok := config.Config.RulesConfig[name]; ok {
		return cfg
	}
	return config.RuleConfig{Enabled: true, Priority: 0}
}
