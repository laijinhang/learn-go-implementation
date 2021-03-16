package auth

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type netrcLine struct {
	machine  string // 机器
	login    string
	password string
}

var (
	netrcOnce sync.Once
	netrc     []netrcLine
	netrcErr  error
)

func parseNetrc(date string) []netrcLine {
	var nrc []netrcLine
	var l netrcLine
	inMacro := false
	for _, line := range strings.Split(date, "\n") {
		if inMacro {
			if line == "" {
				inMacro = false
			}
			continue
		}

		f := strings.Fields(line)
		i := 0
		for ; i < len(f)-1; i += 2 {
			switch f[i] {
			case "machine":
				l = netrcLine{machine: f[i+1]}
			case "default":
				break
			case "login":
				l.login = f[i+1]
			case "password":
				l.password = f[i+1]
			case "macdef":
				inMacro = true
			}
			if l.machine != "" && l.login != "" && l.password != "" {
				nrc = append(nrc, l)
				l = netrcLine{}
			}
		}
		// 只能有一个默认令牌，并且必须在所有机器令牌之后
		if i < len(f) && f[i] == "default" {
			break
		}
	}

	return nrc
}

func netrcPath() (string, error) {
	if env := os.Getenv("NETRC"); env != "" {
		return env, nil
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	base := ".netrc"
	if runtime.GOOS == "windows" {
		base = "_netrc"
	}
	return filepath.Join(dir, base), nil
}

func readNetrc() {
	path, err := netrcPath()
	if err != nil {
		netrcErr = err
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			netrcErr = err
		}
		return
	}

	netrc = parseNetrc(string(data))
}
