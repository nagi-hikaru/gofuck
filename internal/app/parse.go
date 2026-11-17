package app

import (
	"strings"

	"github.com/nagi-hikaru/gofuck/internal/engine"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

// Parse 解析单条简单命令，返回各词的值和原文中的字节范围
func Parse(script string) ([]engine.Word, bool) {
	file, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil || len(file.Stmts) != 1 {
		return nil, false
	}
	stmt := file.Stmts[0]
	call, ok := stmt.Cmd.(*syntax.CallExpr)
	if !ok || stmt.Negated || stmt.Background || len(stmt.Redirs) > 0 ||
		len(call.Assigns) > 0 || len(call.Args) == 0 {
		return nil, false
	}

	words := make([]engine.Word, 0, len(call.Args))
	for _, word := range call.Args {
		// 只接受不依赖 shell 展开的字面参数
		for _, part := range word.Parts {
			switch node := part.(type) {
			case *syntax.Lit:
				if strings.ContainsAny(node.Value, "*?[{}~!") {
					return nil, false
				}
			case *syntax.SglQuoted:
				if node.Dollar {
					return nil, false
				}
			case *syntax.DblQuoted:
				if node.Dollar {
					return nil, false
				}
				for _, inner := range node.Parts {
					if _, ok := inner.(*syntax.Lit); !ok {
						return nil, false
					}
				}
			default:
				return nil, false
			}
		}
		value, err := expand.Literal(nil, word)
		if err != nil {
			return nil, false
		}
		words = append(words, engine.Word{
			Value: value, Start: int(word.Pos().Offset()), End: int(word.End().Offset()),
		})
	}
	return words, true
}
