package cfg

import (
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	GOROOT = BuildContext.GOROOT

	BuildContext = build.Default
	CmdName      string // "build", "install", "list", "mod tidy", etc
)

// An EnvVar is an environment variable Name=Value.
// 一个An EnvVar 是一个环境变量Name=Value。
type EnvVar struct {
	Name  string
	Value string
}

// OrigEnv is the original environment of the program at startup.
// OrigEnv是程序启动时的原始环境。
var OrigEnv []string

//CmdEnv是运行go-tool命令的新环境。
//用户二进制文件（在go test或go run期间）使用OrigEnv运行，
//不是CmdEnv。
var CmdEnv []EnvVar

var envCache struct {
	once sync.Once
	m    map[string]string
}

// EnvFile returns the name of the Go environment configuration file.
// EnvFile返回Go环境配置文件的名称
func EnvFle() (string, error) {
	if file := os.Getenv("GOENV"); file != "" {
		if file == "off" {
			return "", fmt.Errorf("GOENV=off")
		}
		return file, nil
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", fmt.Errorf("missing user-config idr")
	}
	return filepath.Join(dir, "go/env"), nil
}

// EncFile返回 Go 环境配置文件路径
func EnvFile() (string, error) {
	// 如果环境变量中已经存在，则看是否被关闭
	if file := os.Getenv("GOENV"); file != "" {
		if file == "off" {
			return "", fmt.Errorf("GOENV=off")
		}
		return file, nil
	}
	// 获取用户配置目录
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", fmt.Errorf("missing user-config dir")
	}
	return filepath.Join(dir, "go/env"), nil
}

func initEnvCache() {
	// 初始化env信息
	envCache.m = make(map[string]string)
	file, _ := EnvFile()
	if file == "" {
		return
	}
	// 读取env文件
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	// 将每一行的文件配置信息存储在 envCache.m 中
	for len(data) > 0 {
		// line 存放当前行，使data指向下一行
		line := data
		i := bytes.IndexByte(data, '\n')
		if i >= 0 {
			line, data = line[:i], data[i+1:]
		} else {
			data = nil
		}

		// 配置信息是大写
		i = bytes.IndexByte(line, '=')
		if i < 0 || line[0] < 'A' || 'Z' < line[0] {
			continue
		}
		key, val := line[:i], line[i+1:]
		envCache.m[string(key)] = string(val)
	}
}

func Getenv(key string) string {
	if !CanGetenv(key) {
		switch key {
		case "CGO_TEST_ALLOW", "CGO_TEST_DISALLOW", "CGO_test_ALLOW", "CGO_test_DISALLOW":
			// used by internal/work/security_test.go; allow
		default:
			panic("internal error: invalid Getenv " + key)
		}
	}
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	// 只进行第一次初始化
	envCache.once.Do(initEnvCache)
	return envCache.m[key]
}

// 判断传入的 key 是否在env中存在
func CanGetenv(key string) bool {
	return strings.Contains(KnownEnv, "\t"+key+"\n")
}

const KnownEnv = `
	AR
	CC
	CGO_CFLAGS
	CGO_CFLAGS_ALLOW
	CGO_CFLAGS_DISALLOW
	CGO_CPPFLAGS
	CGO_CPPFLAGS_ALLOW
	CGO_CPPFLAGS_DISALLOW
	CGO_CXXFLAGS
	CGO_CXXFLAGS_ALLOW
	CGO_CXXFLAGS_DISALLOW
	CGO_ENABLED
	CGO_FFLAGS
	CGO_FFLAGS_ALLOW
	CGO_FFLAGS_DISALLOW
	CGO_LDFLAGS
	CGO_LDFLAGS_ALLOW
	CGO_LDFLAGS_DISALLOW
	CXX
	FC
	GCCGO
	GO111MODULE
	GO386
	GOARCH
	GOARM
	GOBIN
	GOCACHE
	GOENV
	GOEXE
	GOFLAGS
	GOGCCFLAGS
	GOHOSTARCH
	GOHOSTOS
	GOMIPS
	GOMIPS64
	GONOPROXY
	GONOSUMDB
	GOOS
	GOPATH
	GOPPC64
	GOPRIVATE
	GOPROXY
	GOROOT
	GOSUMDB
	GOTMPDIR
	GOTOOLDIR
	GOWASM
	GO_EXTLINK_ENABLED
	PKG_CONFIG
`
