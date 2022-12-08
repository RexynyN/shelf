/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"glow/cmd/file"
	"glow/cmd/singles"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shelf",
	Short: "A nifty CLI tool for file system power user",
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
	rootCmd.AddCommand(file.DuplicateCmd)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
