package util

import (
	"encoding/json"
	"github.com/cavaliercoder/grab"
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

func SafeJoin(basePath string, unsafePath string) string {
	return filepath.Join(basePath, filepath.Join("/", unsafePath))
}

func DownloadFileWithProgress(displayName string, downloadFilePath string, downloadUrl string) error {
	client := grab.NewClient()
	client.UserAgent = DefaultUserAgent
	request, err := grab.NewRequest(downloadFilePath, downloadUrl)
	if err != nil {
		return err
	}

	log.Printf("Downloading %s...\n", displayName)
	response := client.Do(request)
	
	barSize := int(response.Size / 1000)
	if barSize < 1 {
		// spinner mode
		barSize = -1
	}
	bar := NewBar(barSize)
	defer bar.Finish()

	t := time.NewTicker(200 * time.Millisecond)
	defer t.Stop()
Loop:
	for {
		select {
		case <-t.C:
			bar.Set(int(response.BytesComplete() / 1000))
		case <-response.Done:
			break Loop
		}
	}
	if err := response.Err(); err != nil {
		return err
	}
	return nil
}
