package runner

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func (s *Runner) Run(command *exec.Cmd) {
	fmt.Fprintf(
		s.logExporter,
		commandColor.Sprint("  CMD   ")+
			fmt.Sprintf(
				" %s %s\n",
				filepath.Base(command.Path),
				strings.Join(command.Args, " "),
			),
	)

	buf, err := command.Output()
	if err != nil {
		s.Fatal("failed running command: %v. Your app should still be built, but gofast couldn't do all the work for you", err)
	}

	fmt.Fprint(s.logExporter, string(buf))
}
