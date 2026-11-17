package shell

import (
	_ "embed"
	"fmt"
)

//go:embed history_zsh.sh
var zshHistory string

//go:embed history_bash.sh
var bashHistory string

//go:embed common.sh
var common string

// Script 返回可在指定 shell 中直接加载的集成脚本。
func Script(name string) (string, error) {
	switch name {
	case "zsh":
		return zshHistory + "\n" + common, nil
	case "bash":
		return bashHistory + "\n" + common, nil
	default:
		return "", fmt.Errorf("unsupported shell %q; choose zsh or bash", name)
	}
}
