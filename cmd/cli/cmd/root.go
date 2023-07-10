/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/AnthonyHewins/gofast/internal/cmdline"
	"github.com/spf13/cobra"
)

// build vars
var (
	version string
)

//go:embed templates/*
var f embed.FS

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gofast",
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

		switch len(args) {
		case 1:
			app, err := cmdline.NewAppFromCobra("", cmd.Flags())
			if err != nil {
				return err
			}

			svc{
				logWriter: app.Logger(),
				Name:      "",
				GRPC:      false,
				Proto:     false,
				GoVersion: version,
			}

			fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				fmt.Printf("path=%q, isDir=%v\n", path, d.IsDir())
				return nil
			})
		default:
			return fmt.Errorf("wrong number of args")
		}
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
	f.BoolP("version", "v", false, "Print version")

	pf := rootCmd.PersistentFlags()

	pf.StringP("log-exporter", "", "If blank, log to stdout, else log to this file")
}
