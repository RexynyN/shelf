package duplicate

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"shelf/common"
	"strings"
	"time"

	"github.com/fatih/color"
)

func removeIndexes(items []common.FileStats, idxs []int) []common.FileStats {
	if len(idxs) == 0 {
		return items
	}
	items[idxs[0]] = items[len(items)-1] // Substitui o item removido pelo último
	items = items[:len(items)-1]
	return removeIndexes(items, idxs[1:])
}

// IsNamedDuplicate verifica se o nome do arquivo é um named duplicate, incluindo múltiplas numerações
func isNamedDuplicate(filename string) (bool, string) {
	baseFilename := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Verifica o padrão com múltiplas numerações " (N)" no final do nome antes da extensão
	re := regexp.MustCompile(`^(.*?)( \(\d+\))+$`)
	matches := re.FindStringSubmatch(baseFilename)

	if len(matches) > 1 {
		return true, strings.TrimSpace(matches[1]) + filepath.Ext(filename)
	}
	return false, filename
}

func sameNameDups(files []common.FileStats) map[string][]NamedDuplicate {
	namedDups := make(map[string][]NamedDuplicate)
	for len(files) != 0 {
		numbered, original := isNamedDuplicate(files[0].Filename)
		dup := NamedDuplicate{
			Path:       files[0].Path,
			Filename:   files[0].Filename,
			IsNumbered: numbered,
		}
		namedDups[original] = append(namedDups[original], dup)
		files = files[1:]

		removeIdxs := make([]int, 0)
		for idx, file := range files {
			if strings.HasPrefix(file.Filename, original) {
				numbered, _ := isNamedDuplicate(file.Filename)
				dup := NamedDuplicate{
					Path:       file.Path,
					Filename:   file.Filename,
					IsNumbered: numbered,
				}
				namedDups[original] = append(namedDups[original], dup)
				removeIdxs = append(removeIdxs, idx)
			}
		}
		files = removeIndexes(files, removeIdxs)
	}
	return namedDups
}

func printFate(dups []NamedDuplicate, spared NamedDuplicate) {
	for _, stats := range dups {
		color.Cyan("\t- %s", stats.Path)
	}
	color.Yellow("Spared: %s\n", spared.Path)
}

func quarantineFate(dups []NamedDuplicate, spared NamedDuplicate) {
	dupPath := filepath.Join(common.GetCwd(), "__duplicates__")
	common.CreatePath(dupPath)

	for _, stats := range dups {
		if stats.Path == spared.Path {
			continue
		}
		newPath := filepath.Join(dupPath, stats.Filename)
		if err := os.Rename(stats.Path, newPath); err != nil {
			color.Red("Failed to move file %s: %v", stats.Path, err)
		} else {
			color.Green("Quarantined: %s -> %s", stats.Path, newPath)
		}
	}
}

func removeFate(dups []NamedDuplicate, spared NamedDuplicate) {
	color.Red("Getting ready to delete the duplicates. Be cautious!")
	magicWord := deleteMessages[rand.Intn(len(deleteMessages))]
	color.Yellow("Type '%s' to confirm deletion:", magicWord)
	var typed string
	for {
		fmt.Scanf("%s", &typed)
		if typed == magicWord {
			break
		}
		color.Red("Incorrect word. Try again or press CTRL+C to cancel.")
	}

	for _, stats := range dups {
		if stats.Path == spared.Path {
			continue
		}
		if err := os.Remove(stats.Path); err != nil {
			color.Red("Failed to delete file %s: %v", stats.Path, err)
		} else {
			color.Green("Deleted: %s", stats.Path)
		}
	}
	color.Yellow("Spared: %s", spared.Path)
}

func pickSpareDup(dups []NamedDuplicate) NamedDuplicate {
	strat, _ := flags.GetString("spare")
	switch strings.ToLower(strat) {
	case "oldest":
		return pickOldest(dups)
	case "newest":
		return pickNewest(dups)
	case "random":
		return dups[rand.Intn(len(dups))]
	case "biggest":
		return pickBiggest(dups)
	case "smallest":
		return pickSmallest(dups)
	default:
		color.Red("Invalid spare option: %s", strat)
		os.Exit(1)
	}
	return dups[0]
}

func applyFate(dups []NamedDuplicate, spared NamedDuplicate) {
	if remove, _ := flags.GetBool("remove"); remove {
		removeFate(dups, spared)
	} else if quarantine, _ := flags.GetBool("quarantine"); quarantine {
		quarantineFate(dups, spared)
	} else {
		printFate(dups, spared)
	}
}

func searchNamedDups(files []common.FileStats) {
	namedDups := sameNameDups(files)
	for _, dups := range namedDups {
		if len(dups) <= 1 {
			continue
		}
		applyFate(dups, pickSpareDup(dups))
	}
}

func pickOldest(dups []NamedDuplicate) NamedDuplicate {
	oldest, minTime := dups[0], time.Now().Unix()
	for _, file := range dups {
		if fi, err := os.Stat(file.Path); err == nil {
			if modTime := fi.ModTime().Unix(); modTime < minTime {
				oldest, minTime = file, modTime
			}
		}
	}
	return oldest
}

func pickNewest(dups []NamedDuplicate) NamedDuplicate {
	newest, maxTime := dups[0], int64(0)
	for _, file := range dups {
		if fi, err := os.Stat(file.Path); err == nil {
			if modTime := fi.ModTime().Unix(); modTime > maxTime {
				newest, maxTime = file, modTime
			}
		}
	}
	return newest
}

func pickBiggest(dups []NamedDuplicate) NamedDuplicate {
	biggest, maxSize := dups[0], int64(0)
	for _, file := range dups {
		if fi, err := os.Stat(file.Path); err == nil {
			if size := fi.Size(); size > maxSize {
				biggest, maxSize = file, size
			}
		}
	}
	return biggest
}

func pickSmallest(dups []NamedDuplicate) NamedDuplicate {
	smallest, minSize := dups[0], int64(math.MaxInt64)
	for _, file := range dups {
		if fi, err := os.Stat(file.Path); err == nil {
			if size := fi.Size(); size < minSize {
				smallest, minSize = file, size
			}
		}
	}
	return smallest
}
