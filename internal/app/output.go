package app

import (
	"bytes"
	"io"
	"os"
)

const maxOutput = 128 * 1024

// readOutput 读取本次修正的诊断
func readOutput(path string) (string, bool) {
	if path == "" {
		return "", false
	}

	file, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxOutput+1))
	if err != nil {
		return "", false
	}

	if len(data) > maxOutput {
		lastNewline := bytes.LastIndexByte(data[:maxOutput], '\n')
		if lastNewline == -1 {
			// 限制范围内没有完整的一行
			return "", false
		}
		data = data[:lastNewline+1]
	}

	return string(data), true
}
