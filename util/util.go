package util

import (
	"encoding/json"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"path/filepath"
)

func JsonMarshalPretty(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func EnsureDirectoryExists(path string, printString string) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			return
		} else {
			log.Fatalln("Path for ", printString, " isn't a directory")
		}
	}
	if os.IsNotExist(err) {
		log.Fatalln("Path for ", printString, " doesn't exist")
	}
	log.Fatalln("Error accessing path for ", printString)
}

func EnsureFileExists(path string, printString string) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.Mode().IsRegular() {
			return
		} else {
			log.Fatalln("Path for ", printString, " isn't a regular file")
		}
	}
	if os.IsNotExist(err) {
		log.Fatalln("Path for ", printString, " doesn't exist")
	}
	log.Fatalln("Error accessing path for ", printString)
}

func NewBar(max int, description ...string) *progressbar.ProgressBar {
	bar := progressbar.Default(int64(max), description...)
	return bar
}

func SafeJoin(basePath string, unsafePath string) string {
	return filepath.Join(basePath, filepath.Join("/", unsafePath))
}

func MergePaths(sourcePath string, destPath string) error {
	return filepath.Walk(sourcePath, func(sourceFile string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(sourcePath, sourceFile)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		destPath := filepath.Join(destPath, rel)
		destStat, err := os.Stat(destPath)
		if os.IsNotExist(err) || (err == nil && destStat.Mode().IsRegular()) {
			return os.Rename(sourceFile, destPath)
		} else if err != nil {
			return err
		}
		return nil
	})
}
