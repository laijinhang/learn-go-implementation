package version

import (
	"fmt"
	"go工具/internal/base"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var CmdVersion = &base.Command{
	UsageLine: "go version [-m] [-v] [file ...]",
	Short:     "print Go version",
	Long: `Version prints the build information for Go executables.

Go version reports the Go version userd to build each of the named
executable files.

If no files are named on the command line, go version prints its own
version information.

If a directory is named, go version walks that directory, recursively,
looking for recognized Go binaries and reporting their versions.
By default, go version does not report unrecognized files found
during a directory scan. The -v flag causes it to report unrecognized files.

Tje -m flag causes go version to print each executable's embedded
module version information, when available. In the output, the module
information consists of multiple lines following the version line, each

See also: go doc runtime/debug.BuildInfo.
`,
}

func runVersion(cmd *base.Command, args []string) {
	if len(args) == 0 {
		// go版本、操作系统、系统平台
		fmt.Printf("go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return
	}

	// 输出已编译文件程序版本
	for _, args := range args {
		info, err := os.Stat(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		if info.IsDir() {
			// 扫描指定目录下所有已编译的程序（递归遍历目录）

		} else {
			// 扫描指定文件

		}
	}
}

func scanDir(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

	})
}

// 检查是否为可执行文件
func isExe(file string, info os.FileInfo) bool {
	if runtime.GOOS == "windows" {
		// 检查后缀是否为 .exe
		return strings.HasSuffix(strings.ToLower(file), ".exe")
	}
	// 普通文件且为任何用户都是可以执行的
	return info.Mode().IsRegular() && info.Mode()&0111 != 0
}

func scanFile(file string, info os.FileInfo, mustPrint bool) {
	// 符号链接（不是快捷方式文件）
	if info.Mode()&os.ModeSymlink != 0 {
		// 仅处理文件符号链接
		i, err := os.Stat(file)
		// IsRegular() bool　判断是否是普通文件
		if err != nil || !i.Mode().IsRegular() {
			if mustPrint {
				fmt.Fprintf(os.Stderr, "%s: symlink\n", file)
			}
			return
		}
		info = i
	}
	if !isExe(file, info) {
		if mustPrint {
			fmt.Fprintf(os.Stderr, "%s: not executable file\n", file)
		}
		return
	}

	// 判断文件是不是 类Unix、Dos、Mac、aix平台下的可执行文件
	x, err := openExe(file)
	if err != nil {
		if mustPrint {
			fmt.Fprintf(os.Stderr, "%s: %v\n", file, err)
		}
		return
	}
	defer x.Close()

}

func findVers(x exe) (vers, mod string) {

}
