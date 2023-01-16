package file

import (
	"crypto/sha1"
	"io"
	"log"
	"os"
	"shelf/common"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var DuplicateCmd = &cobra.Command{
	Use:     "duplicates",
	Short:   "Rename a file or a directory of files using various utilities.",
	Example: "glow file rename --extensions \"mp4,png\" --startsWith \"abc\" --endsWith \"123\" --replace \"abc\" --to \"\"\nglow file rename --iterate number --to \"BOGUS VOLUME {}\" --toTitle",
	Long:    ``,
	Run:     runDuplicates,
}

func runDuplicates(cmd *cobra.Command, args []string) {
	var files []os.FileInfo
	var paths []string
	// Gets the pool of files to handle
	if search, _ := cmd.Flags().GetBool("search"); search {
		files, paths = common.ReadFilesRecursive(common.GetCwd())
	} else {
		files = common.ReadFiles(common.GetCwd())
	}

	// 1. Buildup a hash table of the files, where the filesize is the key.
	hash_size := make(map[int64][]string)
	for i := range files {
		hash_size[files[i].Size()] = append(hash_size[files[i].Size()], paths[i])
	}

	// 2. For files with the same size, create a hash table with the hash of their first 1024 bytes; non-colliding elements are unique
	hash_1k := make(map[string][]string)
	for _, files := range hash_size {
		if len(files) < 2 {
			continue
		}
		for _, filename := range files {
			hashSmall := getHash(filename, true)
			hash_1k[hashSmall] = append(hash_1k[hashSmall], filename)
		}
	}

	// 3. For files with the same hash on the first 1k bytes, calculate the hash on the full contents - files with matching ones are NOT unique.
	hash_full := make(map[string]string)
	for _, files := range hash_1k {
		if len(files) < 2 {
			continue // This hash is unique, no files has the same
		}

		for _, filename := range files {
			fullHash := getHash(filename, false)
			duplicate, ok := hash_full[fullHash]

			// If the duplicate exists
			if ok {
				color.Green("Duplicate found: %s and %s\n", filename, duplicate)
			} else {
				hash_full[fullHash] = filename
			}
		}
	}
}

// Initialize the command
func init() {
	// Constraints
	DuplicateCmd.Flags().BoolP("search", "s", false, "Search recursively within the current directory for duplicates.")

	// Fate of the duplicates
	DuplicateCmd.Flags().BoolP("quarantine", "q", true, "Quarantines the duplicates in a subdirectory to be manually handled.")
	DuplicateCmd.Flags().BoolP("remove", "r", false, color.RedString("Delete all duplicates (cannot be undone, be sure of what you're doing)"))

	// Security
	DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")

	// Methods of finding duplicates
	// DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")
	// DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")
}

// Get the hash of a file
func getHash(path string, firstChunk bool) string {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hash := sha1.New()
	if firstChunk {
		// Read the first 1024 bytes
		bytesSlice := make([]byte, 1024)
		_, err = file.Read(bytesSlice)
		if err != nil {
			log.Fatal(err)
		}
		hash.Write(bytesSlice)
	} else {
		// Get all the file contents and make a hash
		hash := sha1.New()
		io.Copy(hash, file)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Digest the hash
	sha1_hash := hash.Sum(nil)
	return string(sha1_hash)
}
