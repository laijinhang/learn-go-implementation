package search

import (
	"path/filepath"
	"strings"
)

// in dir检查路径是否在以dir为根的文件树中。
// 如果是这样，InDir返回相对于dir的等效路径。
// 否则，InDir返回空字符串。
// 即使在有象征性的联系的情况下，InDir也会努力取得成功。
// TODO（rsc）：将internal/tests.inDir替换为对Go 1.12的此函数的调用。
func InDir(path, dir string) string {
	if rel := inDirLex(path, dir); rel == "" {
		return ""
	}
	// 获取path的实际路径（文件链接 -> 实际链接）
	xpath, err := filepath.EvalSymlinks(path)
	if err != nil || xpath == path {
		xpath = ""
	} else {
		if rel := inDirLex(xpath, dir); rel != "" {
			return rel
		}
	}

	xdir, err := filepath.EvalSymlinks(dir)
	if err == nil && xdir != dir {
		if rel := inDirLex(path, xdir); rel != "" {
			return rel
		}
		if xpath != "" {
			if rel := inDirLex(xpath, xdir); rel != "" {
				return rel
			}
		}
	}
	return ""
}

// 如果path是在dir之下，返回path相对于dir路径之下，否则返回 空
func inDirLex(path, dir string) string {
	pv := strings.ToUpper(filepath.VolumeName(path))
	dv := strings.ToUpper(filepath.VolumeName(dir))
	// 比较 path 和 dir 是不是在同一个路径
	path = path[len(pv):]
	dir = path[len(dv):]
	switch {
	default:
		return ""
	case pv != dv: // 不在同一个跟路径下
		return ""
	case len(path) == len(dir):
		if path == dir { // 路径一样
			return "."
		}
		return ""
	case dir == "":
		return path
	case len(path) > len(dir):
		if dir[len(dir)-1] == filepath.Separator {
			if path[:len(dir)] == dir {
				return path[len(dir):]
			}
			return ""
		}
		if path[len(dir)] == filepath.Separator && path[:len(dir)] == dir {
			if len(path) == len(dir)+1 {
				return "."
			}
			return path[len(dir)+1:]
		}
		return ""
	}
}
