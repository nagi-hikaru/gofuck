package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/uniseg"
)

func selectOptions(candidates []string, columns int) ([]string, map[string]string) {
	// PTerm 给选中标记和间距占两列，再留一列避免自动换行。
	width := max(1, columns-3)
	options := make([]string, 0, len(candidates))
	originals := make(map[string]string, len(candidates))
	for i, command := range candidates {
		text := commandLabel(command)
		label := text
		if _, exists := originals[label]; exists || label == "" || uniseg.StringWidth(label) > width {
			// 显示文字可能重复，用序号区分并保留原命令。
			for n := i + 1; ; n++ {
				suffix := fmt.Sprintf(" [%d]", n)
				label = fitLine(text, max(1, width-len(suffix))) + suffix
				if _, exists := originals[label]; !exists {
					break
				}
			}
		}
		options = append(options, label)
		originals[label] = command
	}
	return options, originals
}

// commandLabel 保留命令原有的引号，只转义可能影响终端显示的字符。
func commandLabel(command string) string {
	var label strings.Builder
	for _, r := range command {
		if strconv.IsPrint(r) {
			label.WriteRune(r)
		} else {
			escaped := strconv.QuoteRune(r)
			label.WriteString(escaped[1 : len(escaped)-1])
		}
	}
	return label.String()
}

// fitLine 按终端列宽截短显示文本，不改变真正返回的命令。
func fitLine(text string, columns int) string {
	if uniseg.StringWidth(text) <= columns {
		return text
	}
	const ellipsis = "…"
	var end int
	remaining := columns - uniseg.StringWidth(ellipsis)
	state := -1
	for rest := text; rest != ""; {
		cluster, next, width, nextState := uniseg.FirstGraphemeClusterInString(rest, state)
		if width > remaining {
			break
		}
		end += len(cluster)
		remaining -= width
		rest, state = next, nextState
	}
	return text[:end] + ellipsis
}
