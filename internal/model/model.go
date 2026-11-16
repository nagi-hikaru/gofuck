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

// Candidate 保存修正命令、来源规则和排序权重
type Candidate struct {
	Script string
	Rule   string
	Weight int
}
