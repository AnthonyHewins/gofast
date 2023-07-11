/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/AnthonyHewins/gofast/internal/cmdline"
	"github.com/spf13/cobra"
)

// build vars
var (
	version string
)

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
			app, err := cmdline.NewAppFromCobra(args[0], cmd.Flags())
			if err != nil {
				return err
			}

			app.CreateNewApp()
		default:
			return fmt.Errorf("wrong number of args")
		}

		return nil
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
	f := rootCmd.Flags()
	f.BoolP("version", "v", false, "Print version")
}
