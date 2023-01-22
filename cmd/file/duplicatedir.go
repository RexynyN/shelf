package file

import (
	"github.com/spf13/cobra"
)

// var DuplicateDirCmd = &cobra.Command{
// 	Use:     "duplicate",
// 	Short:   "Rename a file or a directory of files using various utilities.",
// 	Example: "glow file rename --extensions \"mp4,png\" --startsWith \"abc\" --endsWith \"123\" --replace \"abc\" --to \"\"\nglow file rename --iterate number --to \"BOGUS VOLUME {}\" --toTitle",
// 	Long:    ``,
// 	Run:     runDuplicateDir,
// }

func runDuplicateDir(cmd *cobra.Command, args []string) {
	return
}

// Initialize the command
func init() {
	// // Directories
	// DuplicateDirCmd.Flags().StringP("dir", "d", "", "Give the target directory to compare.")

	// // Compare methods
	// DuplicateDirCmd.Flags().BoolP("name", "n", false, "Check for duplicate filenames in the target directory")
}
