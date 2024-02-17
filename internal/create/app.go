package create

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/pflag"
)

//go:embed templates/*
var f embed.FS

type App struct {
	logExporter     io.Writer
	errExporter     io.Writer
	permissionLevel fs.FileMode

	BufID     string
	AppName   string
	Module    string
	GoVersion string
}

func NewAppFromCobra(module string, cmd *pflag.FlagSet) (*App, error) {
	domainParts := strings.Split(module, "/")
	n := len(domainParts)
	if len(domainParts) < 3 {
		return nil, fmt.Errorf("invalid module, not a valid domain: %v", module)
	}

	a := &App{
		logExporter:     os.Stdout,
		errExporter:     os.Stderr,
		permissionLevel: 0755,
		AppName:         domainParts[n-1],
		BufID:           strings.Join(domainParts[n-2:], "/"),
		Module:          module,
		GoVersion:       "",
	}

	semanticVersionPieces := strings.Split(runtime.Version()[2:], ".")
	switch len(semanticVersionPieces) {
	case 2, 3:
		a.GoVersion = strings.Join(semanticVersionPieces[0:2], ".")
	default:
		a.info("unable to determine how to label your go version with this string: %s (got it from calling runtime.Version()). Expected semver: goXX.XX.XX", semanticVersionPieces)
		a.info("going to default to 1.24")
		a.GoVersion = "1.24"
	}

	return a, nil
}

func (a *App) CreateNewApp() {
	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			a.fatal(err.Error())
		}

		if path == "templates" {
			return nil
		}

		targetPath := filepath.Join(a.AppName, strings.TrimPrefix(path, "templates/"))
		if d.IsDir() {
			a.dir(targetPath)
			return nil
		}

		dir, _ := filepath.Split(targetPath)
		a.dir(dir)
		a.template(path, targetPath)
		return nil
	})

	if err != nil {
		a.fatal(err.Error())
	}

	a.run("git", "init")
	a.run("go", "get", "./...")
	fmt.Fprintln(a.logExporter)

	tab := func(s string) {
		a.info("\t" + s)
	}

	a.infoBold("Your app is built at ./%s", a.AppName)
	a.info("Start off by removing what you don't need, if anything:")
	tab("If you don't need a CLI, rm -r cmd/cli")
	tab("If you don't need a server, rm -r cmd/server")
	tab("If you don't need a REST server, delete the serveGRPCGateway function, then the compiler will tell you what other things to delete")
	a.info("Then make sure to do 'go mod tidy' to clean up the dependencies when done to thin the build out")
	fmt.Fprintln(a.logExporter)

	a.info("Using this application template")
	fmt.Fprintln(a.logExporter)

	a.infoBold("General usage")
	tab("You can get logging, tracing, metrics, and database connections using the cmdline package in your app, using App.<Dependency>()")
	tab("Add to this struct to suit your app's needs for bootstrapping dependencies, such as extra DB connections.")
	tab("Use App's methods to pass them into your application code. This is where the bulk of what you'll need to start creating")
	tab("a server/CLI will exist and you can bring it into your code with all the boilerplate essentially done already")
	fmt.Fprintln(a.logExporter)

	a.infoBold("Server (note that by default there is something you NEED to do before it can compile):")
	tab("Start by defining your gRPC interface(s) in api/proto/<serviceName>/v1/<service>.proto, using the provided example")
	tab("Use buf generate to generate your server")
	tab("Go to cmd/server/grpcserver and call <generated packagename>.Register<SvcName>Server(grpcServer, s) to bind the gRPC server")
	tab("If you're using gRPC gateway for a REST server as well, then you'll need to also register that inside the svcHandler slice in the serveGRPCGateway function")
	fmt.Fprintln(a.logExporter)

	a.infoBold("CLI")
	tab("For CLIs, the job is much easier, just make sure you have cobra-cli installed and begin using it.")
	tab("You have logging, tracing, and versioning already set up as persistent flags in the root module you can leverage with the App struct")
}
