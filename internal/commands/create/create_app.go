package create

import (
	"embed"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AnthonyHewins/gofast/internal/runner"
)

//go:embed templates/*
var f embed.FS

type Creator struct {
	r *runner.Runner
	*CreateArgs
}

type CreateArgs struct {
	BufID     string
	AppName   string
	Module    string
	GoVersion string
}

func NewCreator(r *runner.Runner, args *CreateArgs) *Creator {
	return &Creator{r: r, CreateArgs: args}
}

func (a *Creator) CreateNewApp() {
	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			a.r.Fatal(err.Error())
		}

		if path == "templates" {
			return nil
		}

		targetPath := filepath.Join(a.AppName, strings.TrimPrefix(path, "templates/"))
		if d.IsDir() {
			a.r.Dir(targetPath)
			return nil
		}

		dir, _ := filepath.Split(targetPath)
		a.r.Dir(dir)
		a.tmpl(path, strings.TrimSuffix(targetPath, ".tmpl"))
		return nil
	})

	if err != nil {
		a.r.Fatal(err.Error())
	}

	a.cmd("git", "init")
	a.cmd("go", "get", "./...")
	a.r.Newline()

	tab := func(s string) {
		a.r.Info("\t" + s)
	}

	a.r.InfoBold("Your app is built at ./%s", a.AppName)
	a.r.Info("Start off by removing what you don't need, if anything:")
	tab("If you don't need a CLI, rm -r cmd/cli")
	tab("If you don't need a server, rm -r cmd/server")
	tab("If you don't need a REST server, delete the serveGRPCGateway function, then the compiler will tell you what other things to delete")
	a.r.Info("Then make sure to do 'go mod tidy' to clean up the dependencies when done to thin the build out")
	a.r.Newline()

	a.r.Info("Using this application template")
	a.r.Newline()

	a.r.InfoBold("General usage")
	tab("You can get logging, tracing, metrics, and database connections using the cmdline package in your app, using App.<Dependency>()")
	tab("Add to this struct to suit your app's needs for bootstrapping dependencies, such as extra DB connections.")
	tab("Use App's methods to pass them into your application code. This is where the bulk of what you'll need to start creating")
	tab("a server/CLI will exist and you can bring it into your code with all the boilerplate essentially done already")
	a.r.Newline()

	a.r.InfoBold("Server (note that by default there is something you NEED to do before it can compile):")
	tab("Start by defining your gRPC interface(s) in api/proto/<serviceName>/v1/<service>.proto, using the provided example")
	tab("Use buf generate to generate your server")
	tab("Go to cmd/server/grpcserver and call <generated packagename>.Register<SvcName>Server(grpcServer, s) to bind the gRPC server")
	tab("If you're using gRPC gateway for a REST server as well, then you'll need to also register that inside the svcHandler slice in the serveGRPCGateway function")
	a.r.Newline()

	a.r.InfoBold("CLI")
	tab("For CLIs, the job is much easier, just make sure you have cobra-cli installed and begin using it.")
	tab("You have logging, tracing, and versioning already set up as persistent flags in the root module you can leverage with the App struct")
}

func (c *Creator) cmd(cmd string, args ...string) {
	command := exec.Command(cmd, args...)
	command.Dir = "./" + c.AppName
	c.r.Run(command)
}

func (c *Creator) tmpl(tmplPath, target string) {
	c.r.Template(&runner.TemplateArgs{
		FS:         f,
		TmplPath:   tmplPath,
		TargetPath: target,
		Args:       &c.CreateArgs,
	})
}
