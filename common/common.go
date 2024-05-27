package common

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"shelf/config"
	"strings"
)

type FileStats struct {
	Info     os.FileInfo
	Path     string
	Filename string
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func Filter[T any](array []T, funct func(T) bool) []T {
	filtered := make([]T, 1)
	for _, item := range array {
		if funct(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func GetFilename(path string) string {
	sep := string(os.PathSeparator)
	tokens := strings.Split(path, sep)
	return tokens[len(tokens)-1]
}

func ReadFiles(path string) (files []FileStats) {
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Panic(err)
	}

	for _, file := range dirFiles {
		if !file.IsDir() {
			files = append(files, FileStats{
				Info:     file,
				Path:     filepath.Join(path, file.Name()),
				Filename: file.Name(),
			})
		}
	}

	return files
}

func ReadDir(path string) (files []os.FileInfo) {
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	return dirFiles
}

func ReadFilesByExtension(path string, extensions []string) (files []os.FileInfo) {
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dirFiles {
		if !file.IsDir() && checkExtension(file.Name(), extensions) {
			files = append(files, file)
		}
	}

	return files
}

func checkExtension(name string, extensions []string) (sentinel bool) {
	sentinel = false
	for _, extension := range extensions {
		if strings.HasSuffix(name, extension) {
			sentinel = true
			break
		}
	}
	return
}

func ReadFilesRecursive(path string) (files []FileStats) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				files = append(files, FileStats{
					Info:     info,
					Path:     path,
					Filename: GetFilename(path),
				})
			}
			return nil
		})

	if err != nil {
		log.Panic(err)
	}

	return files
}

func GetCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get the current working directory.")
	}

	return cwd
}

func GetExePath() string {
	cwd, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return strings.ReplaceAll(cwd, config.AppConfig.AppName+".exe", "")
}

func GetFileExtension(filename string) string {
	return "." + strings.Split(filename, ".")[1]
}

func GetPureFilename(filename string) string {
	return strings.Split(filename, ".")[0]
}

func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func CreatePath(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatal("Couldn't create the specified path")
		}
	}
}

func RemoveDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
