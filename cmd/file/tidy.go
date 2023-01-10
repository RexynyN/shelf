package file

import "github.com/spf13/cobra"

var TidyCmd = &cobra.Command{
	Use:     "rename",
	Short:   "Rename a file or a directory of files using various utilities.",
	Example: "shelf file rename --extensions \"mp4,png\" --startsWith \"abc\" --endsWith \"123\" --replace \"abc\" --to \"\"\nglow file rename --iterate number --to \"BOGUS VOLUME {}\" --toTitle",
	Long:    ``,
	Run:     runTidy,
}

func runTidy(cmd *cobra.Command, args []string) {
	
}

// Initialize the command
// func init() {
// 	// Selectors
// 	RenameCmd.Flags().String("contains", "", "Selects all files which contains the given literal.")
// 	RenameCmd.Flags().String("startsWith", "", "Selects all files which starts with the given literal.")
// 	RenameCmd.Flags().String("endsWith", "", "Selects all files which ends with the given literal (excluding the file extension).")
// 	RenameCmd.Flags().String("extensions", "", "Selects files by the given pool of file extensions. (separated by comma)")

// 	// Operations
// 	RenameCmd.Flags().String("iterate", "", "Type of value to append to '--to' flag (number, letter, mixed), '--to' must have {} to be replaced by the value.")
// 	RenameCmd.Flags().BoolP("random", "r", false, "Renames all selected files to a random string of characters and numbers.")
// 	RenameCmd.Flags().String("replace", "", "Replace all instances of the given expression, if found. (--to flag is required)")
// 	RenameCmd.Flags().String("replaceOnce", "", "Replace first instance of the given expression, if found. (--to flag is required)")
// 	RenameCmd.Flags().String("to", "", "The value to replace, or the name to be set.")

// 	// String Cases
// 	RenameCmd.Flags().Bool("toUpper", false, "Flips all selected files to Upper Case (after all replace and rename operations)")
// 	RenameCmd.Flags().Bool("toLower", false, "Flips all selected files to Lower Case (after all replace and rename operations)")
// 	RenameCmd.Flags().Bool("toTitle", false, "Flips all selected files to Title Case (after all replace and rename operations)")

// 	// Tools
// 	RenameCmd.Flags().Bool("revert", false, "Revert the last rename operation in the current folder, if any.")
// }
