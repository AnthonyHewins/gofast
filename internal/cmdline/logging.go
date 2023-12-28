package cmdline

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	infoColor     = color.New(color.FgCyan)
	infoBoldColor = color.New(color.FgWhite, color.BgCyan, color.Bold)

	errColor      = color.New(color.FgWhite, color.Bold, color.BgRed)
	dirColor      = color.New(color.FgWhite, color.Bold, color.BgMagenta)
	templateColor = color.New(color.FgWhite, color.Bold, color.BgHiBlue)
	commandColor  = color.New(color.FgWhite, color.Bold, color.BgYellow)
)

// denotes a group of steps in the server creation process
func (s *App) info(str string, args ...any) {
	fmt.Fprintf(
		s.logExporter,
		infoColor.Sprintf(str, args...)+"\n",
	)
}

// denotes a group of steps in the server creation process
func (s *App) infoBold(str string, args ...any) {
	fmt.Fprintf(
		s.logExporter,
		infoBoldColor.Sprintf(str, args...)+"\n",
	)
}

func (s *App) fatal(str string, args ...any) {
	fmt.Fprintf(
		s.errExporter,
		errColor.Sprintf(str, args...)+"\n",
	)

	os.Exit(1)
}
