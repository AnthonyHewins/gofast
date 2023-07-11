package cmdline

import (
	"fmt"
	"os"
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
