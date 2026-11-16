package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nagi-hikaru/gofuck/internal/app"
	"github.com/nagi-hikaru/gofuck/internal/cli"
	"github.com/nagi-hikaru/gofuck/internal/status"
)

func main() {
	os.Exit(run())
}

func run() int {
	// 创建上下文
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 创建 root 命令
	root := cli.New(app.NewFucker())

	// 执行 root 命令
	err := root.ExecuteContext(ctx)
	if err != nil {
		fmt.Fprintln(root.ErrOrStderr(), err)
	}

	// 返回 exit code
	return status.ExitCode(err)
}
