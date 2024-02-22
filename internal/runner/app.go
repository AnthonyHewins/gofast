package runner

import (
	"io"
	"io/fs"
)

type Runner struct {
	logExporter     io.Writer
	errExporter     io.Writer
	permissionLevel fs.FileMode
}

func NewRunner(logExporter, errExporter io.Writer, permLevel fs.FileMode) *Runner {
	return &Runner{
		logExporter:     logExporter,
		errExporter:     errExporter,
		permissionLevel: permLevel,
	}
}
