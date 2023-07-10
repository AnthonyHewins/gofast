package cmdline

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type App struct {
	AppName   string
	Module    string
	GoVersion string
}

func NewAppFromCobra(module string, cmd *cobra.Command) (*App, error) {
	domainParts := strings.Split(module, "/")
	if len(domainParts) < 2 {
		return nil, fmt.Errorf("invalid module, not a valid domain: %v", module)
	}

	version := runtime.Version()[2:]

	return &App{
		AppName:   domainParts[len(domainParts)-1],
		Module:    module,
		GoVersion: "",
	}, nil
}

func (s *App) run(cmd string, args ...string) []byte {
	fmt.Fprintf(
		s.logWriter,
		commandColor.Sprint("  CMD   ")+fmt.Sprintf(" %s %s\n", cmd, strings.Join(args, " ")),
	)

	command := exec.Command(cmd, args...)
	command.Dir = "./" + s.Name
	buf, err := command.Output()
	if err != nil {
		s.fatal("failed running command: %v", err)
	}

	return buf
}

// denotes a group of steps in the server creation process
func (s *App) info(str string, args ...any) {
	fmt.Fprintf(
		s.logWriter,
		infoColor.Sprint("  INFO  ")+" "+fmt.Sprintf(str, args...)+"\n",
	)
}

func (s *App) fatal(str string, args ...any) {
	fmt.Fprintf(
		s.logWriter,
		errColor.Sprint("  FATAL ")+" "+fmt.Sprintf(str, args...)+"\n",
	)

	os.Exit(1)
}

func (s *App) dir(dir string) {
	dir = path.Join(s.Name, dir)

	fmt.Fprintf(
		s.logWriter,
		"\t"+dirColor.Sprint("  DIR   ")+" "+color.MagentaString(dir)+"\n",
	)

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		s.fatal(err.Error())
	}
}

func (s *App) template(tmplFileWithoutExtension string) {
	outDir, base := path.Dir(tmplFileWithoutExtension), path.Base(tmplFileWithoutExtension)

	if outDir != "." {
		s.dir(outDir)
	}

	fmt.Fprintf(
		s.logWriter,
		"\t"+templateColor.Sprint("TEMPLATE")+" "+color.HiBlueString(tmplFileWithoutExtension)+"\n",
	)

	tmplFile := path.Join("templates", outDir, base+".tmpl")
	buf, err := f.ReadFile(tmplFile)
	if err != nil {
		s.fatal("missing template file %s: %v", tmplFile, err)
	}

	tmpl, err := template.New(tmplFile).Parse(string(buf))
	if err != nil {
		s.fatal("template for %s is invalid: %v", tmplFile, err)
	}

	intermediateBuffer := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(intermediateBuffer, s)
	if err != nil {
		s.fatal("failed executing template: %v", err)
	}

	if intermediateBuffer.Len() == 0 {
		s.info("you passed some args that resulted in this file having no output; skipping")
	}

	outPath := fmt.Sprintf("%s/%s", s.Name, tmplFileWithoutExtension)
	file, err := os.Create(outPath)
	if err != nil {
		s.fatal("failed creating the output file %s: %v", outPath, err)
	}

	_, err = fmt.Fprint(file, intermediateBuffer.Bytes())
	if err != nil {
		s.fatal("failed writing: %v", err)
	}
}
