package duplicate

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"shelf/common"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type NamedDuplicate struct {
	Path       string
	Filename   string
	IsNumbered bool
}

type Duplicate struct {
	Stats        common.FileStats
	OriginalPath string
}

var (
	DuplicateCmd = &cobra.Command{
		Use:     "duplicates",
		Short:   "Find duplicated files using a set of utilities for more precise or loose parameters.",
		Example: "// TODO",
		Long:    "",
		Run:     runDuplicates,
	}
	CWD            = common.GetCwd()
	files          []common.FileStats
	flags          *pflag.FlagSet
	deleteMessages = []string{"danger", "permanent", "delete", "loop", "fallback", "backup", "oxymoron", "responsibility", "deletion"}
)

func init() {
	DuplicateCmd.Flags().BoolP("search", "s", false, "Search recursively within the current directory for duplicates.")
	DuplicateCmd.Flags().Bool("quiet", false, "Hides all logs of found duplicates, just prints essential information.")
	DuplicateCmd.Flags().BoolP("name", "n", false, "Search for same-name files (homonymous) within the directory, including files with a number suffix. Eg. 'file (1).jpg'.")
	DuplicateCmd.Flags().BoolP("quarantine", "q", false, "Quarantines the duplicates in a subdirectory to be manually handled.")
	DuplicateCmd.Flags().BoolP("remove", "r", false, color.RedString("Deletes all duplicates (cannot be undone, be sure of what you're doing)."))
	DuplicateCmd.Flags().String("spare", "oldest", "Strategy for sparing duplicates. Options ['oldest' (Default), 'newest', 'random', 'first', 'biggest', 'smallest'].")
	DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")
}

func runDuplicates(cmd *cobra.Command, args []string) {
	flags = cmd.Flags()
	color.Cyan("Reading files...")

	if search, _ := flags.GetBool("search"); search {
		files = common.ReadFilesRecursive(CWD)
	} else {
		files = common.ReadFiles(CWD)
	}

	if name, _ := flags.GetBool("name"); name {
		searchNamedDuplicates(files)
		return
	}

	sizeHash := groupByFileSize(files)
	partialCount := len(sizeHash)

	firstChunkHash := hashFirstChunks(sizeHash)
	duplicates, fullCount := findFullDuplicates(firstChunkHash)

	printResults(partialCount, fullCount, duplicates)
	handleDuplicates(duplicates)
}

func groupByFileSize(files []common.FileStats) map[int64][]common.FileStats {
	color.Cyan("Grouping files by size...")
	sizeHash := make(map[int64][]common.FileStats)
	for _, file := range files {
		size := file.Info.Size()
		sizeHash[size] = append(sizeHash[size], file)
	}
	return sizeHash
}

func hashFirstChunks(sizeHash map[int64][]common.FileStats) map[string][]common.FileStats {
	color.Cyan("Hashing first 1024 bytes of files...")

	chunkHash := make(map[string][]common.FileStats)
	for _, group := range sizeHash {
		if len(group) < 2 {
			continue
		}

		for _, file := range group {
			hash := getHash(file.Path, true)
			chunkHash[hash] = append(chunkHash[hash], file)
		}
	}
	return chunkHash
}

func findFullDuplicates(chunkHash map[string][]common.FileStats) (map[string][]common.FileStats, int) {
	color.Cyan("Finding full duplicates by hashing entire files...")

	duplicates := make(map[string][]common.FileStats)
	fullHashes := make(map[string]common.FileStats)
	fullCount := 0

	for _, group := range chunkHash {
		if len(group) < 2 {
			continue
		}

		for _, file := range group {
			hash := getHash(file.Path, false)
			if original, exists := fullHashes[hash]; exists {
				duplicates[hash] = append(duplicates[hash], original, file)
				fullCount++
			} else {
				fullHashes[hash] = file
			}
		}
	}
	return duplicates, fullCount
}

func printResults(partialCount, fullCount int, duplicates map[string][]common.FileStats) {
	if quiet, _ := flags.GetBool("quiet"); quiet {
		return
	}

	color.Cyan("Partial matches: %d", partialCount)
	color.Cyan("Full matches: %d", fullCount)
	for hash, group := range duplicates {
		for _, file := range group {
			path := strings.ReplaceAll(file.Path, CWD, "")
			color.Green("Duplicate: %s [Hash: %s]", path, hash)
		}
	}
}

func handleDuplicates(duplicates map[string][]common.FileStats) {
	if remove, _ := flags.GetBool("remove"); remove {
		deleteDuplicates(duplicates)
	} else if quarantine, _ := flags.GetBool("quarantine"); quarantine {
		quarantineDuplicates(duplicates)
	}
}

func deleteDuplicates(duplicates map[string][]common.FileStats) {
	color.Red("Preparing to delete duplicates. Make sure you know what you're doing.")
	magicWord := deleteMessages[rand.Intn(len(deleteMessages))]
	color.Yellow("Type '%s' to confirm deletion:", magicWord)

	var typed string
	for {
		fmt.Scanf("%s", &typed)
		if typed == magicWord {
			break
		}
		color.Red("Incorrect word. Try again.")
	}

	for _, group := range duplicates {
		for i, file := range group {
			if i > 0 {
				if err := os.Remove(file.Path); err != nil {
					log.Printf("Failed to delete %s: %v", file.Path, err)
				}
			}
		}
	}
}

func quarantineDuplicates(duplicates map[string][]common.FileStats) {
	color.Cyan("Quarantining duplicates...")

	quarantineDir := filepath.Join(CWD, "__duplicates__")
	common.CreatePath(quarantineDir)

	for index, group := range duplicates {
		groupDir := filepath.Join(quarantineDir, fmt.Sprintf("group_%d", index))
		common.CreatePath(groupDir)
		for _, file := range group {
			newPath := filepath.Join(groupDir, file.Info.Name())
			if err := os.Rename(file.Path, newPath); err != nil {
				log.Printf("Failed to quarantine %s: %v", file.Path, err)
			}
		}
	}
}

func getHash(path string, firstChunk bool) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", path, err)
	}
	defer file.Close()

	hash := sha1.New()
	if firstChunk {
		buffer := make([]byte, 1024)
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatalf("Failed to read file %s: %v", path, err)
		}
		hash.Write(buffer[:bytesRead])
	} else {
		if _, err := io.Copy(hash, file); err != nil {
			log.Fatalf("Failed to hash file %s: %v", path, err)
		}
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
