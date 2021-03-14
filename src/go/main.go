// go 1.15.8
package main

import (
	"flag"
	"fmt"
	"learn-go-implementation/src/go/internal/base"
	"learn-go-implementation/src/go/internal/bug"
	"learn-go-implementation/src/go/internal/cfg"
	"learn-go-implementation/src/go/internal/envcmd"
	"learn-go-implementation/src/go/internal/help"
	"learn-go-implementation/src/go/internal/modload"
	"learn-go-implementation/src/go/internal/work"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	base.Go.Commands = []*base.Command{
		bug.CmdBuf,
		work.CmdBuild,
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	args := flag.Args()
	if len(args) < 1 {
		base.Usage()
	}

	if args[0] == "get" || args[0] == "help" {
		if modload.Init(); !modload.Enabled() {
		}
	}

	cfg.CmdName = args[0]
	// 处理 "go help <topic>" 命令
	if args[0] == "help" {
		help.Help(os.Stderr, args[1:])
		return
	}

	// Diagnose common mistake: GOPATH==GOROOT.
	// This setting is equivalent to not setting GOPATH at all,
	// which is not what most people want when they do it.
	//诊断常见错误：GOPATH==GOROOT。
	//这个设置相当于根本不设置GOPATH，
	//这不是大多数人做这件事时想要的。
	if gopath := cfg.BuildContext.GOPATH; filepath.Clean(gopath) == filepath.Clean(runtime.GOROOT()) {
		// 警告：GOPATH设为GOROOT是无效的
		fmt.Fprintf(os.Stderr, "warning: GOPATH set to GOROOT (%s) has no effect\n", gopath)
	} else {
		for _, p := range filepath.SplitList(gopath) {
			if p == "" {
				continue
			}
			if strings.HasPrefix(p, "~") {
				fmt.Fprintf(os.Stderr, "go: GOPATH entry cannot start with shell metacharacter '~': %q\n", p)
				os.Exit(2)
			}
			if !filepath.IsAbs(p) {
				if cfg.Getenv("GOPATH") == "" {
					cfg.BuildContext.GOPATH = ""
				} else {
					fmt.Fprintf(os.Stderr, "go: GOPATH entry is relative; must be absolute path: %q.\nFor more details see: 'go help gopath'\n", p)
					os.Exit(2)
				}
			}
		}
	}

	if fi, err := os.Stat(cfg.GOROOT); err != nil || !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "go: cannot find GOROOT directory: %v\n", cfg.GOROOT)
		os.Exit(2)
	}

	// 存储初始系统变量
	cfg.OrigEnv = os.Environ()
	cfg.CmdEnv = envcmd.MKEnv()
	for _, env := range cfg.CmdEnv {
		if os.Getenv(env.Name) != env.Value {
			os.Setenv(env.Name, env.Value)
		}
	}

BigCmdLoop:
	for bigCmd := base.Go; ; {
		for _, cmd := range bigCmd.Commands {
			if cmd.Name() != args[0] {
				continue
			}
			if len(cmd.Commands) > 0 {
				bigCmd = cmd
				args = args[1:]
				if len(args) == 0 {
					help.PrintUsage(os.Stdout, bigCmd) // 将信息从 os.Stdout 输出
					base.SetExitStatus(2)              // 设置退出状态为2
					base.Exit()                        // 执行回收函数
				}
			}
			if args[0] == "help" { // 输出help命令信息
				help.Help(os.Stdout, append(strings.Split(cfg.CmdName, " "), args[1:]...))
				return
			}
			cfg.CmdName += " " + args[0]
			continue BigCmdLoop
		}
	}
}

func init() {
	base.Usage = mainUsage
}

func mainUsage() {
	help.PrintUsage(os.Stderr, base.Go)
	os.Exit(2)
}
