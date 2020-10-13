package util

import (
	"encoding/json"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"path/filepath"
	"time"
)

func JsonMarshalPretty(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func EnsureDirectoryExists(path string, name string) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			return
		} else {
			log.Fatal("Path for ", name, " isn't a directory")
		}
	}
	if os.IsNotExist(err) {
		log.Fatal("Path for ", name, " doesn't exist")
	}
	log.Fatal("Error accessing path for ", name)
}

func EnsureFileExists(path string, name string) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.Mode().IsRegular() {
			return
		} else {
			log.Fatal("Path for ", name, " isn't a regular file")
		}
	}
	if os.IsNotExist(err) {
		log.Fatal("Path for ", name, " doesn't exist")
	}
	log.Fatal("Error accessing path for ", name)
}

func NewBar(max int, description ...string) *progressbar.ProgressBar {
	bar := progressbar.Default(int64(max), description...)
	// wait for throttle duration so initial render doesn't skip
	// TODO: Fix
	time.Sleep(65 * time.Millisecond)
	bar.RenderBlank()
	return bar
}

func SafeJoin(paths ...string) string {
	if len(paths) < 1 {
		return ""
	}
	if len(paths) < 2 {
		return paths[0]
	}
	unsafeJoinElems := append([]string{"/"}, paths[1:]...)
	safePath := filepath.Join(unsafeJoinElems...)
	return filepath.Join(paths[0], safePath)
}
