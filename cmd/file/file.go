/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package file

import (
	"github.com/spf13/cobra"
)

var FileCmd = &cobra.Command{
	Use:   "file",
	Short: "Utility commands to work with files.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

func init() {
	// FileCmd.AddCommand(renameCmd)

	// VideoCmd.Flags().StringVarP(&urlPath, "url", "u", "", "The url to ping")
	// if err := VideoCmd.MarkFlagRequired("url"); err != nil {
	// 	fmt.Println(err)
	// }

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
