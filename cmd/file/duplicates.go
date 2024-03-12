package file

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"shelf/common"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Duplicate struct {
	Stats        common.FileStats
	OriginalPath string
}

var DuplicateCmd = &cobra.Command{
	Use:     "duplicates",
	Short:   "Find duplicated files using a set of utilities for more precise or loose parameters.",
	Example: "// TODO",
	Long:    ``,
	Run:     runDuplicates,
}

var CWD string = common.GetCwd()
var files []common.FileStats

// Initialize the command
func init() {
	// Constraints
	DuplicateCmd.Flags().BoolP("search", "s", false, "Search recursively within the current directory for duplicates.")
	DuplicateCmd.Flags().Bool("quiet", false, "Hides all logs of found duplicates, just prints essencial information")
	DuplicateCmd.Flags().BoolP("name", "n", false, "Search for same-name files (homonymous) within the directory, including files with a number suffix. Eg. 'file (1).jpg'")

	// Fate of the duplicates
	DuplicateCmd.Flags().BoolP("quarantine", "q", false, "Quarantines the duplicates in a subdirectory to be manually handled.")
	DuplicateCmd.Flags().BoolP("remove", "r", false, color.RedString("Deletes all duplicates (cannot be undone, be sure of what you're doing)"))

	// Security
	DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")

	// Methods of finding duplicates
	// DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")
}

func detectDupNumbering(filename string) (bool, string) {
	name := common.GetPureFilename(filename)
	if !strings.HasSuffix(name, ")") || !strings.Contains(name, "(") {
		return false, filename
	}

	// Get the content of the last parenthesis to confirm if it is a number
	parStart, parEnd := strings.LastIndex(name, "("), strings.LastIndex(name, ")")
	content := name[parStart+1 : parEnd]
	if _, err := strconv.Atoi(content); err != nil {
		return false, filename
	}
	return true, strings.TrimSpace(strings.ReplaceAll(filename, "("+content+")", ""))
}

func sameNameDup(files []common.FileStats) {
	namedDups := make(map[string][]string)
	for len(files) != 0 {
		_, original := detectDupNumbering(files[0].Filename)
		namedDups[original] = append(namedDups[original], files[0].Path)
		for _, file := range files {
			if strings.HasPrefix(file.Filename, original) {
				namedDups[original] = append(namedDups[original], file.Path)
			}
		}
		files = files[1:]
	}

}

func searchDups() {

}

func runDuplicates(cmd *cobra.Command, args []string) {
	var full int = 0
	var partial int = 0
	// Gets the pool of files to handle
	color.Cyan("Reading files...")
	if search, _ := cmd.Flags().GetBool("search"); search {
		files = common.ReadFilesRecursive(CWD)
	} else {
		files = common.ReadFiles(CWD)
	}

	if name, _ := cmd.Flags().GetBool("name"); name {
		sameNameDup(files)
		return
	}

	color.Cyan("Size hashing...")
	// 1. Buildup a hash table of the files, where the filesize is the key.
	hash_size := make(map[int64][]common.FileStats)
	for i := range files {
		hash_size[files[i].Info.Size()] = append(hash_size[files[i].Info.Size()], files[i])
	}

	color.Cyan("Byte hashing...")
	// 2. For files with the same size, create a hash table with the hash of their first 1024 bytes; non-colliding elements are unique
	hash_1k := make(map[string][]common.FileStats)
	for _, files := range hash_size {
		if len(files) < 2 {
			continue
		}

		partial += len(files)
		for _, stats := range files {
			hashSmall := getHash(stats.Path, true)
			hash_1k[hashSmall] = append(hash_1k[hashSmall], stats)
		}
	}

	printDups, _ := cmd.Flags().GetBool("quiet")
	printDups = !printDups

	// 3. For files with the same hash on the first 1k bytes, calculate the hash on the full contents - files with matching ones are NOT unique.
	duplicates := make(map[string][]common.FileStats)
	hash_full := make(map[string]common.FileStats)
	color.Cyan("Searching for duplicates through Hashes...")
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
				if printDups {
					full++
					original := strings.ReplaceAll(stats.Path, CWD, "")
					dup := strings.ReplaceAll(duplicate.Path, CWD, "")
					color.Green("Duplicate found: %s and %s\n", dup, original)
				}
				// Append to the map
				duplicates[fullHash] = append(duplicates[fullHash], duplicate, stats)
			} else {
				hash_full[fullHash] = stats
			}
		}
	}

	fmt.Println("Partial ", partial)
	fmt.Println("Full ", full)

	// Fate of the duplicates
	// Quarantine them
	if quarantine, _ := cmd.Flags().GetBool("quarantine"); quarantine {
		common.CreatePath("__duplicates__")

		dupPath := filepath.Join(common.GetCwd(), "__duplicates__")
		index := 0
		for _, array := range duplicates {
			intraPath := filepath.Join(dupPath, fmt.Sprint(index))
			common.CreatePath(intraPath)
			for _, stats := range array {
				// Ignore the error because if we hit a miss, it's a duplicated entry in the slice
				// (We hate having to manually write a set to prevent dups)
				_ = os.Rename(stats.Path, filepath.Join(intraPath, stats.Info.Name()))
			}
			index++
		}
	}

	// REMOVE THEM
	if remove, _ := cmd.Flags().GetBool("remove"); remove {
		spared := ""
		for _, array := range duplicates {
			for index, stats := range array {
				// The first element of the array is always spared, whichever it is
				if index != 0 {
					if stats.Path == spared {
						continue
					}

					err := os.Remove(stats.Path)
					if err != nil {
						log.Fatal(err)
					}
				} else {
					spared = stats.Path
					printer := strings.ReplaceAll(stats.Path, CWD, "")
					color.Yellow("Spared: %s", printer)
				}
			}
		}
	}
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
		read, err := file.Read(bytesSlice)

		// 1024 bytes is an arbitrary number, but a file may have less than 1024 bytes
		// so make an extra case for that.
		if err != nil && read < 1024 {
			bytesSlice = make([]byte, read)
			file.Read(bytesSlice)
		} else if err != nil {
			log.Fatal(err) // Log if any other error occured
		}

		hash.Write(bytesSlice)
	} else {
		// Get all the file contents and make a hash
		io.Copy(hash, file)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Digest the hash
	sha1_hash := hash.Sum(nil)
	return string(sha1_hash)
}
