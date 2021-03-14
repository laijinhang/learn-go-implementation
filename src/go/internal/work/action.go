package work

type Builder struct {
	WorkDir string // 临时工作目录（以filepath.Separator结尾）
	//actionCache map[cache]
}

type cacheKey struct {
	mode string
}

type Action struct {
	Mode    string
	Package *load.Pa
}
