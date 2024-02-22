/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/AnthonyHewins/gofast/internal/commands/create"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create VALID-GO-MODULE-PATH",
	Short: "Create a new applcation",
	RunE: func(cmd *cobra.Command, args []string) error {
		var moduleParts []string
		switch len(args) {
		case 1:
			moduleParts = strings.Split(args[0], "/")
		default:
			return fmt.Errorf("wrong number of args")
		}

		n := len(moduleParts)
		if len(moduleParts) < 3 {
			return fmt.Errorf("invalid module, not a valid domain: %v", args[0])
		}

		a, err := newApp(cmd.Flags())
		if err != nil {
			return err
		}

		c := create.NewCreator(a.r, &create.CreateArgs{
			AppName:   moduleParts[n-1],
			BufID:     strings.Join(moduleParts[n-2:], "/"),
			Module:    args[0],
			GoVersion: a.r.GoVersion(),
		})

		c.CreateNewApp()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
