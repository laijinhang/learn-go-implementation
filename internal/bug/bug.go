package bug

import (
	"go工具/internal/base"
)

var CmdBuf = &base.Command{
	Run:       runBug,
	UsageLine: "go bug",
	Short:     "start a bug report", // 启动错误报告
	Long: `
Bug opens the default browser and starts a new bug report.
The report includes useful system information.
	`, // Bug打开默认浏览器并启动新的Bug报告。
	// 报告包括有用的系统信息。
	Commands: nil,
}

func init() {
}

func runBug(cmd *base.Command, args []string) {
	if len(args) > 0 {
		base.Fatalf("go bug: bug takes no arguments") // base.Fatalf("go bug:bug不带参数")
	}
}
