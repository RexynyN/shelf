package duplicate

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"shelf/common"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func removeIndexes(items []common.FileStats, idxs []int) []common.FileStats {
	if len(idxs) == 0 {
		return items
	}

	// Put the last item in the position of the deleted item, so it's faster
	items[idxs[0]] = items[len(items)-1]
	items = items[:len(items)-1]
	return removeIndexes(items, idxs[1:])
}

func detectDupNumbering(filename string) (bool, string) {
	name := common.GetPureFilename(filename)
	if !strings.Contains(name, ")") || !strings.Contains(name, "(") {
		return false, filename
	}

	// Get the content of the last parenthesis to confirm if it is a number
	parStart, parEnd := strings.LastIndex(name, "("), strings.LastIndex(name, ")")
	content := name[parStart+1 : parEnd]
	if _, err := strconv.Atoi(content); err != nil {
		return false, filename
	}

	// Returns the filename without the numbered part
	return true, strings.TrimSpace(strings.ReplaceAll(filename, "("+content+")", ""))
}

func sameNameDups(files []common.FileStats) map[string][]NamedDuplicate {
	namedDups := make(map[string][]NamedDuplicate)
	for len(files) != 0 {
		numbered, original := detectDupNumbering(files[0].Filename)
		dup := NamedDuplicate{
			Path:       files[0].Path,
			Filename:   files[0].Filename,
			IsNumbered: numbered,
		}
		namedDups[original] = append(namedDups[original], dup)
		// Get the first file out of the way
		files = files[1:]

		removeIdxs := make([]int, 0)
		for idx, file := range files {
			if strings.HasPrefix(file.Filename, original) {
				numbered, _ := detectDupNumbering(file.Filename)
				dup := NamedDuplicate{
					Path:       file.Path,
					Filename:   file.Filename,
					IsNumbered: numbered,
				}
				namedDups[original] = append(namedDups[original], dup)
				removeIdxs = append(removeIdxs, idx)
			}
		}
		fmt.Println(namedDups[original])
		files = removeIndexes(files, removeIdxs)
	}
	return namedDups
}

func printFate(dups []NamedDuplicate, spared NamedDuplicate) {
	sparedFile := spared.Path
	for _, stats := range dups {
		fmt.Println("	- ", stats.Path)
	}
	color.Yellow("Spared: %s \n", sparedFile)
}

// func quarantineFate(dups []NamedDuplicate, spared NamedDuplicate) {
// 	dupPath := filepath.Join(common.GetCwd(), "__duplicates__")
// 	intraPath := filepath.Join(dupPath, fmt.Sprint(index))
// 	common.CreatePath(intraPath)
// 	for _, stats := range dups {
// 		// Ignore the error because if we hit a miss, it's a duplicated entry in the slice
// 		// (We hate having to manually write a set to prevent dups)
// 		_ = os.Rename(stats.Path, filepath.Join(intraPath, common.GetPureFilename(stats.Filename)))
// 	}
// }

func removeFate(dups []NamedDuplicate, spared NamedDuplicate) {
	color.Red("Getting ready to delete the duplicatres, I hope you know what you're doing...")
	magicWord := deleteMessages[rand.Intn(len(deleteMessages))]
	color.Yellow("Type the following word to guarantee that you wanna PERMANENTLY DELETE these files: '%s'", magicWord)
	typed := ""
	for {
		fmt.Scanf("%s", &typed)
		if typed == magicWord {
			break
		}
		color.Red("That's not the right word, try again. (If want to cancel the deletion, press CTRL+C)")
	}

	sparedFile := spared.Filename
	for _, stats := range dups {
		if stats.Path == sparedFile {
			continue
		}

		err := os.Remove(stats.Path)
		if err != nil {
			log.Fatal(err)
		}
	}
	color.Yellow("Spared: %s", strings.ReplaceAll(sparedFile, CWD, ""))
}

func pickSpareDup(dups []NamedDuplicate) NamedDuplicate {
	strat, _ := flags.GetString("spare")

	switch strings.ToLower(strat) {
	// IT USES THE MODIFICATION TIME, NOT THE CREATED AT TIME!!
	case "oldest", "old":
		oldest, date := dups[0], time.Now().Unix()
		for _, file := range dups {
			fi, err := os.Stat(file.Path)
			if err != nil {
				color.Red("The file %s couldn't be opened to compare its modification date, skipping...", file.Filename)
				continue
			}

			mod := fi.ModTime().Unix()
			if mod < date {
				oldest = file
				date = mod
			}
		}
		return oldest

	// IT USES THE MODIFICATION TIME, NOT THE CREATED AT TIME!!²
	case "newest", "new":
		newest, date := dups[0], int64(0)
		for _, file := range dups {
			fi, err := os.Stat(file.Path)
			if err != nil {
				color.Red("The file %s couldn't be opened to compare its modification date, skipping...", file.Filename)
				continue
			}

			mod := fi.ModTime().Unix()
			if mod > date {
				newest = file
				date = mod
			}
		}
		return newest

	// Should be treated as random, because we can't guarantee what file is gonna be the first
	case "first":
		return dups[0]

	case "random":
		return dups[rand.Intn(len(dups))]

	case "biggest", "big":
		biggest, bSize := dups[0], int64(0)
		for _, file := range dups {
			fi, err := os.Stat(file.Path)
			if err != nil {
				color.Red("The file %s couldn't be opened to compare its size, skipping...", file.Filename)
				continue
			}

			size := fi.Size()
			if size > bSize {
				biggest = file
				bSize = size
			}
		}

		return biggest

	case "smallest", "small":
		smallest, sSize := dups[0], int64(math.MaxInt64)
		for _, file := range dups {
			fi, err := os.Stat(file.Path)
			if err != nil {
				color.Red("The file '%s' couldn't be opened to compare its size, skipping...", file.Filename)
				continue
			}

			size := fi.Size()
			if size > sSize {
				smallest = file
				sSize = size
			}
		}
		return smallest

	default:
		color.Red("The given spare option '%s' is not valid! Give a valid option and try again.", strat)
		os.Exit(1)
	}

	return dups[0]
}

func applyFate(dups []NamedDuplicate, spared NamedDuplicate) {
	if remove, _ := flags.GetBool("remove"); remove {
		removeFate(dups, spared)
	} else if quan, _ := flags.GetBool("remove"); quan {
		color.Red("Nuh uh")
		// quarantineFate(dups, spared)
	} else {
		printFate(dups, spared)
	}

}

// TODO: Use this to create a "canonical" strat (which the unique not-numbered dup is the spared, if more than one exists, fallback to another strat)
func isNumbered(dup NamedDuplicate) bool {
	return dup.IsNumbered
}

func isNotNumbered(dup NamedDuplicate) bool {
	return !dup.IsNumbered
}

func searchNamedDups(files []common.FileStats) {
	// Not the best place to put this, but I really do not care
	if quan, _ := flags.GetBool("quarantine"); quan {
		common.CreatePath("__duplicates__")
	}

	namedDups := sameNameDups(files)
	for _, dups := range namedDups {
		// No named duplicate
		if len(dups) == 1 {
			continue
		}

		// Find out what file to spare using the given strat
		applyFate(dups, pickSpareDup(dups))
	}
}
