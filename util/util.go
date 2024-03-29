package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func JsonMarshalPretty(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func EnsureFileExists(path string) {
	ensureExists(path, false, false, false)
}

func EnsureDirectoryExists(path string, mustEmpty bool, canCreate bool) {
	ensureExists(path, true, mustEmpty, canCreate)
}

func ensureExists(path string, mustDirectory bool, mustEmpty bool, canCreate bool) {
	fileType := "file"
	if mustDirectory {
		fileType = "directory"
	}
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() && mustDirectory {
			if mustEmpty {
				files, err := ioutil.ReadDir(path)
				if err != nil {
					log.Fatalln(err)
				}
				if len(files) > 0 {
					log.Fatalf("path '%s' is not empty\n", path)
				}
			}
			return
		} else if stat.Mode().IsRegular() && !mustDirectory {
			return
		} else {
			log.Fatalf("Path '%s' isn't a %s\n", path, fileType)
		}
	} else if os.IsNotExist(err) {
		if canCreate && mustDirectory {
			if err := os.MkdirAll(path, 0777); err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatalf("Path '%s' doesn't exist\n", path)
		}
	} else {
		log.Fatalf("Error accessing path '%s'\n", path)
	}
}

func SafeJoin(basePath string, unsafePath string) string {
	return filepath.Join(basePath, filepath.Clean("/"+unsafePath))
}

func MergePaths(sourcePath string, destPath string) error {
	err := filepath.Walk(sourcePath, func(sourceFile string, sourceStat os.FileInfo, err error) error {
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
			if err := os.Rename(sourceFile, destPath); err != nil {
				return err
			}
			if sourceStat.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		} else if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return os.RemoveAll(sourcePath)
}

func TempFile() (*os.File, error) {
	return ioutil.TempFile(".", "cornstone-temp-")
}

func TempDir() (string, error) {
	return ioutil.TempDir(".", "cornstone-temp-")
}
