/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/AnthonyHewins/gofast/internal/create"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create VALID-GO-MODULE-PATH",
	Short: "Create a new applcation",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 1:
			app, err := create.NewAppFromCobra(args[0], cmd.Flags())
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
