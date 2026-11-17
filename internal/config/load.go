package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Load 读取配置文件
func Load() error {
	// 拼接配置文件路径
	dir := os.Getenv("GOFUCK_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dir = filepath.Join(home, ".config")
	}
	path := filepath.Join(dir, "gofuck", "config.toml")

	// 加载配置文件
	loader := viper.New()
	loader.SetConfigFile(path)
	if err := loader.ReadInConfig(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}

	if err := loader.UnmarshalExact(&Config); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}
