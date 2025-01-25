package diff

import (
	"shelf/common"

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
	DiffCmd = &cobra.Command{
		Use:     "diff",
		Args:    cobra.ExactArgs(2),
		Short:   "Find duplicated files using a set of utilities for more precise or loose parameters.",
		Example: "shelf file rename --extensions \"mp4,png\" --startsWith \"abc\" --endsWith \"123\" --replace \"abc\" --to \"\"\nshelf file rename --iterate number --to \"BOGUS VOLUME {}\" --toTitle",
		Long:    "",
		Run:     runDiff,
	}
	currentDir     = ""
	targetDir      = ""
	files          []common.FileStats
	flags          *pflag.FlagSet
	deleteMessages = []string{"danger", "permanent", "delete", "loop", "fallback", "backup", "oxymoron", "responsibility", "deletion"}
)

func init() {
	DiffCmd.Flags().BoolP("search", "s", false, "Search recursively within the current directory for duplicates.")
	// DuplicateCmd.Flags().Bool("quiet", false, "Hides all logs of found duplicates, just prints essential information.")
	// DuplicateCmd.Flags().BoolP("name", "n", false, "Search for same-name files (homonymous) within the directory, including files with a number suffix. Eg. 'file (1).jpg'.")
	// DuplicateCmd.Flags().BoolP("quarantine", "q", false, "Quarantines the duplicates in a subdirectory to be manually handled.")
	// DuplicateCmd.Flags().BoolP("remove", "r", false, color.RedString("Deletes all duplicates (cannot be undone, be sure of what you're doing)."))
	// DuplicateCmd.Flags().String("spare", "oldest", "Strategy for sparing duplicates. Options ['oldest' (Default), 'newest', 'random', 'first', 'biggest', 'smallest'].")
	// DuplicateCmd.Flags().BoolP("enforce", "e", false, "Enforces the files are down-to-the-byte clones to apply its fate.")
}

func runDiff(cmd *cobra.Command, args []string) {
	currentDir, targetDir = args[0], args[1]

	var dups []common.FileStats
	if search, _ := cmd.Flags().GetBool("search"); search {
		curFiles, targetFiles := common.ReadFilesRecursive(currentDir), common.ReadFilesRecursive(targetDir)
		dups = diffFiles(curFiles, targetFiles)
	} else {
		curFiles, targetFiles := common.ReadFiles(currentDir), common.ReadFiles(targetDir)
		dups = diffFiles(curFiles, targetFiles)
	}

	for _, dup := range dups {
		color.Yellow("%s is duplicated!", dup.Filename)
	}

}

func diffFiles(base, target []common.FileStats) (dups []common.FileStats) {
	fileMap := make(map[string]common.FileStats)
	for _, file := range base {
		fileMap[file.Filename] = file
	}

	for _, file := range target {
		search, ok := fileMap[file.Filename]
		if ok {
			dups = append(dups, search)
		}
	}
	return
}
