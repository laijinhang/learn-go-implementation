package cache

import (
	"fmt"
	"go/internal/cfg"
	"os"
	"path/filepath"
	"sync"
)

func Default() *Cache {
	defaultOnce.Do(initDefaultCache)
	return defaultCache
}

var (
	defaultOnce  sync.Once
	defaultCache *Cache
)

func initDefaultCache() {
	dir := DefaultDir()
}

var (
	defaultDirOnce sync.Once
	defaultDir     string
	defaultDirErr  error
)

func DefaultDir() string {
	defaultDirOnce.Do(func() {
		defaultDir = cfg.Getenv("GOCACHE")
		if filepath.IsAbs(defaultDir) || defaultDir == "off" {
			return
		}
		if defaultDir != "" {
			defaultDir = "off"
			defaultDirErr = fmt.Errorf("GOCACHE is not an absolute path")
			return
		}
		// 计算默认地址
		dir, err := os.UserCacheDir()
		if err != nil {
			defaultDir = "off"
			defaultDirErr = fmt.Errorf("GOCACHE is not defined and %v", err)
			return
		}
		defaultDir = filepath.Join(dir, "go-build")
	})

	return defaultDir
}
