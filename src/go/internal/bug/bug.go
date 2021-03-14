package bug

import (
	"bytes"
	"fmt"
	"io"
	"learn-go-implementation/src/go/internal/base"
	"learn-go-implementation/src/go/internal/cfg"
	"path/filepath"
	"runtime"
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
	CmdBuf.Flag.BoolVar(&cfg.BuildV, "v", false, "")
}

func runBug(cmd *base.Command, args []string) {
	if len(args) > 0 {
		base.Fatalf("go bug: bug takes no arguments") // base.Fatalf("go bug:bug不带参数")
	}
	var buf bytes.Buffer
	buf.WriteString(bugHeader)
	printGoVersion(&buf)

	buf.WriteString("### Does this issue reproduce with the latest release?\n\n\n")
	printEnvDetails(&buf)
	buf.WriteString(bugFooter)

	//body := buf.String()
	//url := "https://github.com/golang/go/issues/new?body=" + urlpkg.QueryEscape(body)
	//if !web.OpenBrowser(url) {
	//	fmt.Print("Please file a new issue at golang.org/issue/new using this template:\n\n")
	//	fmt.Print(body)
	//}
}

const bugHeader = `<!-- Please answer these questions before submitting your issue. Thanks! -->

`
const bugFooter = `### What did you do?

<!--
If possible, provide a recipe for reproducing the error.
A complete runnable program is good.
A link on play.golang.org is best.
-->

`

func printGoVersion(w io.Writer) {
	fmt.Fprintf(w, "### What version of Go are you using (`go version`)?\n\n")
	fmt.Fprintf(w, "<pre>\n")
	fmt.Fprintf(w, "$ go version\n")
	printCmdOut(w, "", "go", "version")
	fmt.Fprintf(w, "</pre>\n")
	fmt.Fprintf(w, "\n")
}

func printEnvDetails(w io.Writer) {
	fmt.Fprintf(w, "### What operating system and processor architecture are you using (`go env`)?\n\n")
	fmt.Fprintf(w, "<details><summary><code>go env</code> Output</summary><br><pre>\n")
	fmt.Fprintf(w, "$ go env\n")
	printCmdOut(w, "", "go", "env")
	printGoDetails(w)
	printOSDetails(w)
	printCDetails(w)
	fmt.Fprintf(w, "</pre></details>\n\n")
}

func printGoDetails(w io.Writer) {
	printCmdOut(w, "GOROOT/bin/go version: ", filepath.Join(runtime.GOROOT(), "bin/go"), "version")
	printCmdOut(w, "GOROOT/bin/go tool compile -V: ", filepath.Join(runtime.GOROOT(), "bin/go"), "tool", "compile", "-V")
}

func printCDetails(w io.Writer) {}

func printOSDetails(w io.Writer) {}

func printCmdOut(w io.Writer, prefix, path string, args ...string) {

}
