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
			log.Fatal("Path for ", printString, " isn't a directory")
		}
	}
	if os.IsNotExist(err) {
		log.Fatal("Path for ", printString, " doesn't exist")
	}
	log.Fatal("Error accessing path for ", printString)
}

func EnsureFileExists(path string, printString string) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.Mode().IsRegular() {
			return
		} else {
			log.Fatal("Path for ", printString, " isn't a regular file")
		}
	}
	if os.IsNotExist(err) {
		log.Fatal("Path for ", printString, " doesn't exist")
	}
	log.Fatal("Error accessing path for ", printString)
}

func NewBar(max int, description ...string) *progressbar.ProgressBar {
	bar := progressbar.Default(int64(max), description...)
	return bar
}

func SafeJoin(basePath string, unsafePath string) string {
	return filepath.Join(basePath, filepath.Join("/", unsafePath))
}
