package rules

import (
	"sort"

	"github.com/nagi-hikaru/gofuck/internal/model"
)

var Rules []model.Rule

func init() {
	// 注册规则
	Rules = []model.Rule{
		GitRule(),
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
