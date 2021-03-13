package modload

import (
	"fmt"
	"go/internal/base"
	"go/internal/cfg"
	"go/internal/internal/search"
	"os"
	"path/filepath"
)

var (
	cwd            string
	initialized    bool    // 初始化
	mustUseModules = false // 必须使用modules

	modRoot string // mod路径

	CmdModInit bool // 运行 'go mod init' 命令
)

func Init() {
	// 已初始化
	if initialized {
		return
	}
	// 进行初始化
	initialized = true

	env := cfg.Getenv("GO11MODULE")
	switch env {
	default:
		base.Fatalf("go: unknown environment setting GO111MODULE=%s", env)
	case "auto", "":
		mustUseModules = false
	case "on":
		mustUseModules = true
	case "off":
		mustUseModules = false
		return
	}
	//禁止Git提示输入密码。
	//仅对2.3.0或更高版本有效，但避免
	//早期版本的提示太难了。
	//如果用户已显式设置GIT_TERMINAL_PROMPT=1，则保持
	//提示。
	//见golang.org/issue/9341和golang.org/issue/12706。
	if os.Getenv("GIT_TERMINAL_PROMPT") == "" {
		os.Setenv("GIT_TERMINAL_PROMPT", "0")
	}

	//通过Git禁用任何ssh连接池。
	//如果Git子进程将子进程分叉到后台以缓存新连接，
	//那个孩子一直开着标准输出/标准输出。Git子进程退出后，
	//os/exec希望能够读取stdout/stderr管道
	//直到EOF获得Git子过程之前退出的所有数据。
	//EOF直到孩子离开时才出现，因为孩子
	//正在握住管道的写端。
	//这很不幸，但至少出现了两次
	//（见golang.org/issue/13453和golang.org/issue/16104）
	//当它出现时会让用户感到困惑。
	//如果用户已显式设置GIT_SSH或GIT_SSH_命令，
	//假设他们知道自己在做什么，而不是踩在上面。
	//但默认为关闭ControlMaster。
	if os.Getenv("GIT_SSH") == "" && os.Getenv("GIT_SSH_COMMAND") == "" {
		os.Setenv("GIT_SSH_COMMAND", "ssh -o ControlMaster=no")
	}

	var err error
	// 获取当前程序运行的工作目录（绝对路径），linux下对应pwd
	cwd, err = os.Getwd()
	if err != nil {
		base.Fatalf("go: %v", err)
	}

	if CmdModInit {
		modRoot = cwd
	} else {
		modRoot = findModuleRoot(cwd)
		if modRoot == "" {
			if !mustUseModules {
				//GO111MODULE是“auto”，我们找不到模块根。
				//保持GOPATH模式。
				return
			}
		} else if search.InDir(modRoot, os.TempDir()) == "." {
			// 警告，mod路径在 /temp 之下
			modRoot = ""
			fmt.Fprintf(os.Stderr, "go: warning: ignoring go.mod in system temp root %v\n", os.TempDir())
		}
	}

	// 现在处于module模式，安装挂钩编译到工作
}

// mod 加载命令
func Enabled() bool {
	Init()
	return modRoot != "" || mustUseModules
}

func findModuleRoot(dir string) (root string) {
	// 去除路径中多余字符
	dir = filepath.Clean(dir)

	for {
		// 如果找到 go.mod 文件，则直接返回
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		// 获取当前路径的最后一个分隔符之前的部分（不包含分隔符）
		// 递归往前找
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		dir = d
	}
	return ""
}
