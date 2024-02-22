package cmd

import (
	"os"

	"github.com/AnthonyHewins/gofast/internal/runner"
	"github.com/spf13/pflag"
)

type app struct {
	r *runner.Runner
}

func newApp(cmd *pflag.FlagSet) (*app, error) {
	r := runner.NewRunner(os.Stdout, os.Stderr, 0755)

	a := &app{
		r: r,
	}

	return a, nil
}
