package file

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var DuplicateCmd = &cobra.Command{
	Use:     "duplicate",
	Short:   "Rename a file or a directory of files using various utilities.",
	Example: "glow file rename --extensions \"mp4,png\" --startsWith \"abc\" --endsWith \"123\" --replace \"abc\" --to \"\"\nglow file rename --iterate number --to \"BOGUS VOLUME {}\" --toTitle",
	Long:    ``,
	Run:     runTidy,
}

func runDuplicate(cmd *cobra.Command, args []string) {

}

// Initialize the command
func init() {
	// Fate of the duplicates
	DuplicateCmd.Flags().BoolP("quarantine", "q", true, "Quarantines the duplicates in a subdirectory to be manually handled.")
	DuplicateCmd.Flags().BoolP("remove", "r", false, color.RedString("Delete all duplicates (cannot be undone, be sure of what you're doing)"))

	// Security
	DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")

	// Methods of finding duplicates
	// DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")
	// DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")
}
