package base

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type Command struct {
	//
	Run func(cmd *Command, args []string)

	UsageLine string

	// go 或 go help
	// 简短的内容
	Short string

	// go help <命令> 输出的内容输出
	// 命令完整的内容
	Long string

	// Flag是特定于此命令的一组标志。

	// 存储不同的Command打印内容
	// bug
	// build
	// clean
	// doc
	// dev
	// ...
	Commands []*Command
}

var Go = &Command{
	UsageLine: "go",
	Long:      "Go is a tool for managing Go source code.",
}

func (c *Command) LongName() string {
	name := c.UsageLine
	// 如果UsageLine行里存在 " ["
	if i := strings.Index(name, " ["); i >= 0 {
		name = name[:i]
	}
	if name == "go" {
		return ""
	}
	// 去掉前面带 "go " 的内容
	return strings.TrimPrefix(name, "go ")
}

func (c *Command) Name() string {
	name := c.LongName()
	// 该方法是返回 " " 在 name 中第一次出现的位置，如果没有找到，则返回-1
	if i := strings.LastIndex(name, " "); i >= 0 {
		name = name[i+1:]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stdout, "usage: %s\n", c.UsageLine)
	fmt.Fprintf(os.Stdout, "Run 'go help %s' for details.\n", c.LongName())
	SetExitStatus(2)
	Exit()
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

var exitStatus = 0
var exitMu sync.Mutex

func SetExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

var atExitFuncs []func()

func Exit() {
	for _, f := range atExitFuncs {
		f()
	}
	os.Exit(exitStatus)
}

func Fatalf(format string, args ...interface{}) {
	Errorf(format, args)
	Exit()
}

func Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
	SetExitStatus(1)
}

var Usage func()
