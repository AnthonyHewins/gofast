package runner

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

func (s *Runner) Run(command *exec.Cmd) {
	fmt.Fprintf(
		s.logExporter,
		commandColor.Sprint("  CMD   ")+fmt.Sprintf(" %s %s\n", command.Path, strings.Join(command.Args, " ")),
	)

	buf, err := command.Output()
	if err != nil {
		s.Fatal("failed running command: %v. Your app should still be built, but gofast couldn't do all the work for you", err)
	}

	fmt.Fprint(s.logExporter, string(buf))
}

func (s *Runner) Dir(dir string) {
	fmt.Fprintf(
		s.logExporter,
		dirColor.Sprint("  DIR   ")+" "+color.MagentaString(dir)+"\n",
	)

	err := os.MkdirAll(dir, s.permissionLevel)
	if err != nil {
		s.Fatal(err.Error())
	}
}

type TemplateArgs struct {
	FS         embed.FS
	TmplPath   string
	TargetPath string
	Args       any
}

func (s *Runner) Template(args *TemplateArgs) {
	fmt.Fprintf(
		s.logExporter,
		templateColor.Sprint("TEMPLATE")+" "+color.HiBlueString(args.TargetPath)+"\n",
	)

	buf, err := args.FS.ReadFile(args.TmplPath)
	if err != nil {
		s.Fatal(err.Error())
	}

	tmpl, err := template.New(args.TmplPath).Parse(string(buf))
	if err != nil {
		s.Fatal("template located at embedded filesystem for %s that targets path %s is invalid: %v", args.TmplPath, args.TargetPath, err)
	}

	intermediateBuffer := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(intermediateBuffer, args.Args)
	if err != nil {
		s.Fatal("failed executing template: %v", err)
	}

	output := bytes.TrimSpace(intermediateBuffer.Bytes())
	if len(output) == 0 {
		fmt.Fprintf(
			s.logExporter,
			templateColor.Sprint("  SKIP  ")+" You passed some input that resulted in this file being empty; skipping\n",
		)
	}

	err = os.WriteFile(args.TargetPath, append(output, '\n'), s.permissionLevel)
	if err != nil {
		s.Fatal("failed creating the output file %s: %v", args.TargetPath, err)
	}
}
