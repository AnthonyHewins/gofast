package cmdline

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

func (s *App) run(cmd string, args ...string) {
	fmt.Fprintf(
		s.logExporter,
		commandColor.Sprint("  CMD   ")+fmt.Sprintf(" %s %s\n", cmd, strings.Join(args, " ")),
	)

	command := exec.Command(cmd, args...)
	command.Dir = "./" + s.AppName
	buf, err := command.Output()
	if err != nil {
		s.fatal("failed running command: %v. Your app should still be built, but gofast couldn't do all the work for you", err)
	}

	fmt.Fprintf(s.logExporter, string(buf))
}

func (s *App) dir(dir string) {
	fmt.Fprintf(
		s.logExporter,
		dirColor.Sprint("  DIR   ")+" "+color.MagentaString(dir)+"\n",
	)

	err := os.MkdirAll(dir, s.permissionLevel)
	if err != nil {
		s.fatal(err.Error())
	}
}

func (s *App) template(tmplPath, targetPath string) {
	fmt.Fprintf(
		s.logExporter,
		templateColor.Sprint("TEMPLATE")+" "+color.HiBlueString(targetPath)+"\n",
	)

	buf, err := f.ReadFile(tmplPath)
	if err != nil {
		s.fatal(err.Error())
	}

	tmpl, err := template.New(tmplPath).Parse(string(buf))
	if err != nil {
		s.fatal("template located at embedded filesystem for %s that targets path %s is invalid: %v", tmplPath, targetPath, err)
	}

	intermediateBuffer := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(intermediateBuffer, s)
	if err != nil {
		s.fatal("failed executing template: %v", err)
	}

	output := bytes.TrimSpace(intermediateBuffer.Bytes())
	if len(output) == 0 {
		fmt.Fprintf(
			s.logExporter,
			templateColor.Sprint("  SKIP  ")+" You passed some input that resulted in this file being empty; skipping\n",
		)
	}

	err = os.WriteFile(targetPath, append(output, '\n'), s.permissionLevel)
	if err != nil {
		s.fatal("failed creating the output file %s: %v", targetPath, err)
	}
}
