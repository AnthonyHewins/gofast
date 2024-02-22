package runner

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
func (s *Runner) Info(str string, args ...any) {
	fmt.Fprintf(
		s.logExporter,
		infoColor.Sprintf(str, args...)+"\n",
	)
}

func (s *Runner) Newline() {
	fmt.Fprintln(s.logExporter)
}

// denotes a group of steps in the server creation process
func (s *Runner) InfoBold(str string, args ...any) {
	fmt.Fprintf(
		s.logExporter,
		infoBoldColor.Sprintf(str, args...)+"\n",
	)
}

func (s *Runner) Fatal(str string, args ...any) {
	fmt.Fprintf(
		s.errExporter,
		errColor.Sprintf(str, args...)+"\n",
	)

	os.Exit(1)
}
