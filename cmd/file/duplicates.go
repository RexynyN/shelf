package file

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"shelf/common"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Duplicate struct {
	Original  string
	Duplicate string
}

var DuplicateCmd = &cobra.Command{
	Use:     "duplicates",
	Short:   "Find duplicated files using a set of utilities for more precise or loose parameters.",
	Example: "// TODO",
	Long:    ``,
	Run:     runDuplicates,
}

func runDuplicates(cmd *cobra.Command, args []string) {
	var CWD string = common.GetCwd()
	var files []common.FileStats

	// Gets the pool of files to handle
	if search, _ := cmd.Flags().GetBool("search"); search {
		files = common.ReadFilesRecursive(CWD)
	} else {
		files = common.ReadFiles(CWD)
	}

	// 1. Buildup a hash table of the files, where the filesize is the key.
	hash_size := make(map[int64][]common.FileStats)
	for i := range files {
		hash_size[files[i].Info.Size()] = append(hash_size[files[i].Info.Size()], files[i])
	}

	// 2. For files with the same size, create a hash table with the hash of their first 1024 bytes; non-colliding elements are unique
	hash_1k := make(map[string][]common.FileStats)
	for _, files := range hash_size {
		if len(files) < 2 {
			continue
		}
		for _, stats := range files {
			hashSmall := getHash(stats.Path, true)
			hash_1k[hashSmall] = append(hash_1k[hashSmall], stats)
		}
	}

	// 3. For files with the same hash on the first 1k bytes, calculate the hash on the full contents - files with matching ones are NOT unique.
	duplicates := make(map[int][]common.FileStats)
	hash_full := make(map[string]common.FileStats)
	index := 0
	for _, files := range hash_1k {
		if len(files) < 2 {
			continue // This hash is unique, no files in the map has the same
		}

		// Iterate through the underlying array in the hashmap entry
		for _, stats := range files {
			// Get the hash for the entire file
			fullHash := getHash(stats.Path, false)
			duplicate, ok := hash_full[fullHash]

			// If the duplicate exists
			if ok {
				// Print the result
				original := strings.ReplaceAll(stats.Path, CWD, "")
				dup := strings.ReplaceAll(duplicate.Path, CWD, "")
				color.Green("Duplicate found: %s and %s\n", original, dup)

				// Append to the map the duplicate
				duplicates[index] = append(duplicates[index], duplicate, stats)
				index++
			} else {
				hash_full[fullHash] = stats
			}
		}
	}

	// Fate of the duplicates

	// Quarantine them
	if quarantine, _ := cmd.Flags().GetBool("quarantine"); quarantine {
		common.CreatePath("__duplicates__")

		// dupPath := filepath.Join(common.GetCwd(), "__duplicates__")
		for index, array := range duplicates {
			// intraPath := filepath.Join(dupPath, fmt.Sprint(index))
			// common.CreatePath(intraPath)
			for _, stats := range array {
				// 	err := os.Rename(stats.Path, filepath.Join(intraPath, stats.Info.Name()))
				// 	if err != nil {
				// 		log.Fatal(err)
				// 	}
				fmt.Println(index, array, stats)

			}
		}
	}

	// REMOVE THEM
}

// Initialize the command
func init() {
	// Constraints
	DuplicateCmd.Flags().BoolP("search", "s", false, "Search recursively within the current directory for duplicates.")

	// Fate of the duplicates
	DuplicateCmd.Flags().BoolP("quarantine", "q", false, "Quarantines the duplicates in a subdirectory to be manually handled.")
	DuplicateCmd.Flags().BoolP("remove", "r", false, color.RedString("Deletes all duplicates (cannot be undone, be sure of what you're doing)"))
	DuplicateCmd.Flags().BoolP("print", "p", true, "Prints the names of the found duplicates on screen.")

	// Security
	DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")

	// Methods of finding duplicates
	// DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")
	// DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")
}

// Get the hash of a file, with a flag to return the fist 1024 bytes chunk
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
