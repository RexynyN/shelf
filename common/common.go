package common

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func ReadFiles(path string) (files []os.FileInfo) {
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dirFiles {
		if !file.IsDir() {
			files = append(files, file)
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

	return strings.ReplaceAll(cwd, "glow.exe", "")
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
