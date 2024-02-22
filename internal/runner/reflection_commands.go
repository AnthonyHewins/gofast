package runner

import (
	"runtime"
	"strings"
)

func (r *Runner) GoVersion() string {
	goVersion := "1.24"
	semanticVersionPieces := strings.Split(runtime.Version()[2:], ".")
	switch len(semanticVersionPieces) {
	case 2, 3:
		goVersion = strings.Join(semanticVersionPieces[0:2], ".")
	default:
		r.Info("unable to determine how to label your go version with this string: %s (got it from calling runtime.Version()). Expected semver: goXX.XX.XX", semanticVersionPieces)
		r.Info("going to default to %s", goVersion)
	}

	return goVersion
}
