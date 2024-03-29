/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"{{ .Module }}/internal/bootstrap"
)

// build vars
var (
	version string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .AppName }}",
	Short: "Generate CLIs and servers quickly for Golang",
	Long: `Generate CLIs/servers with a lot of good but opinionated best practices already done for you.
Features:
	- cobra CLI for CLIs
	- slog for logging
	- sqlx, OTEL for database
	- Prometheus metrics`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		if v, _ := cmd.Flags().GetBool("version"); v {
			fmt.Println(version)
			return nil
		}

		return fmt.Errorf("missing args")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	f := rootCmd.Flags()
	f.BoolP("toggle", "t", false, "Help message for toggle")
	f.BoolP("version", "v", false, "Print version")

	pf := rootCmd.PersistentFlags()

	pf.String(bootstrap.LogLevel, "", "Log level to use. None for no logs, or debug, warn/warning, info, error/err")
	pf.String(bootstrap.LogExporter, "", "Log exporter to use. By default, it goes off log level: info/debug go to STDOUT, warn/error to STDERR. Specify 'stderr' to write to stderr, and anything else opens a file")
	pf.String(bootstrap.LogFmt, "", "Log format to use. Blank or 'json' will create a json logger, or you can use logfmt/text")
	pf.Bool(bootstrap.LogSource, false, "Make all logging show where the log occurred")

	pf.String(bootstrap.Host, "localhost", "The database host to use")
	pf.String(bootstrap.Name, "", "The database name to use")
	pf.String(bootstrap.User, "", "The database user to use")
	pf.String(bootstrap.Password, "", "The database password to use")
	pf.Uint16(bootstrap.Port, 5432, "The database port to use")

	pf.String("trace-exporter", "", "Export data using this exporter. Options are stdout (can be configured to go to a file using trace-exporter-arg), otlp with gRPC, jaegar. Use 'none' or leave blank to skip tracing")
	pf.String("trace-exporter-arg", "", "Export data using this URI. For otlp and jaegar this will point to the collector of tracing, for stdout this will point to a file rather than stdout")
	pf.Duration("trace-exporter-timeout", time.Second*5, "How long the tracer will try to export before it abandons the whole process (not supported for all trace exporters)")
}
