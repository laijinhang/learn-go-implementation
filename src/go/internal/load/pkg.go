package load

import "go/build"

var (
	ModInit func()

	ModBinDir func()
	ModLookup func(parentPath string, parentIsStd bool, path string) (dir, realPath string, err error)
)

type Package struct {
	PackagePublic
	Internal PackageInternal

	GoFiles []string `json:",omitempty"`
}

type PackagePublic struct {
	Dir        string
	ImportPath string
}

type PackageInternal struct {
	Build   *build.Package
	Imports []*Package
}
