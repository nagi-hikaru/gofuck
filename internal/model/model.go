package model

// Word 保存解析后的词和它在原命令中的字节范围 [Start, End)
type Word struct {
	Value string
	Start int
	End   int
}

// Command 保存待修正的命令和采集到的诊断输出
type Command struct {
	Script string // 源命令
	Words  []Word // 解析后词
	Output string // 输出
}

// Rule 定义规则的启用条件、优先级和候选生成函数
type Rule struct {
	Name     string
	Enabled  bool
	Priority int // 数值越小，越先检查
	Suggest  func(Command) []string
}

// Candidate 保存修正命令、来源规则和排序权重
type Candidate struct {
	Script string
	Rule   string
	Weight int
}
