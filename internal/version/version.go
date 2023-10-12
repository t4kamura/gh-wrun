package version

import (
	"github.com/hashicorp/go-version"
	"github.com/t4kamura/gh-wrun/internal/subproc"
)

func CheckGhVersion(required string) (bool, error) {
	v, err := subproc.GetGhVersion()
	if err != nil {
		return false, err
	}

	ver, err := version.NewVersion(v)
	requiredVer, err := version.NewVersion(required)

	return ver.GreaterThanOrEqual(requiredVer), err
}
