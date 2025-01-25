/*
Copyright © 2023 Breno Nogueira breno.s.nogueira@gmail.com
*/
package cmd

import (
	"os"

	"shelf/cmd/diff"
	"shelf/cmd/duplicate"
	"shelf/cmd/file"

	"shelf/cmd/singles"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shelf",
	Short: "A nifty CLI tool for the file system power user",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Finished Commands
	rootCmd.AddCommand(singles.WhoamiCmd)
	rootCmd.AddCommand(file.RenameCmd)
	rootCmd.AddCommand(duplicate.DuplicateCmd)
	rootCmd.AddCommand(diff.DiffCmd)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
