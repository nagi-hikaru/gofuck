package engine

// Rule 定义规则的启用条件、优先级和候选生成函数
type Rule struct {
	Name     string
	Enabled  bool
	Priority int // 数值越小，越先检查
	Match    func(Command) bool
	Fix      func(Command) []string
}
