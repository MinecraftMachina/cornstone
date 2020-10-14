package util

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"github.com/cavaliercoder/grab"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ExtractFileConfig struct {
	ArchivePath string
	Common      ExtractCommonConfig
}

type ExtractReaderConfig struct {
	Data   *bytes.Reader
	Common ExtractCommonConfig
}

type ExtractCommonConfig struct {
	BasePath string
	DestPath string
	Unwrap   bool
}

func getFileFullNameFromHeader(file archiver.File) (string, error) {
	switch file.Header.(type) {
	case zip.FileHeader:
		return file.Header.(zip.FileHeader).Name, nil
	case *tar.Header:
		return file.Header.(*tar.Header).Name, nil
	default:
		return "", errors.New("unsupported header type")
	}
}

func processFile(file archiver.File, config ExtractCommonConfig) error {
	fullName, err := getFileFullNameFromHeader(file)
	if err != nil {
		return err
	}
	fullName, err = filepath.Rel(config.BasePath, fullName)
	if err != nil {
		return nil
	}
	// skip files outside of basePath
	if strings.HasPrefix(fullName, "..") {
		return nil
	}
	if config.Unwrap {
		firstPathIndex := strings.Index(fullName, string(filepath.Separator))
		if firstPathIndex == -1 {
			return nil
		} else {
			fullName = fullName[firstPathIndex+1:]
		}
	}
	fullName = filepath.Join(config.DestPath, fullName)
	if file.IsDir() {
		if err := os.MkdirAll(fullName, file.Mode()); err != nil {
			return err
		}
	} else {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(fullName, data, file.Mode()); err != nil {
			return err
		}
		if err := os.Chtimes(fullName, time.Now(), file.ModTime()); err != nil {
			return err
		}
	}
	return nil
}

//ExtractArchiveFromReader: All files not a child of basePath will be skipped.
func ExtractArchiveFromFile(walker archiver.Walker, config ExtractFileConfig) error {
	if err := walker.Walk(config.ArchivePath, func(file archiver.File) error {
		return processFile(file, config.Common)
	}); err != nil {
		return err
	}
	return nil
}

//ExtractArchiveFromReader: All files not a child of basePath will be skipped.
func ExtractArchiveFromReader(reader archiver.Reader, config ExtractReaderConfig) error {
	if err := reader.Open(config.Data, int64(config.Data.Len())); err != nil {
		return err
	}
	for {
		file, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if err := processFile(file, config.Common); err != nil {
			return err
		}
	}
	return nil
}

func DownloadAndExtract(walker archiver.Walker, downloadUrl string, config ExtractCommonConfig) error {
	tempFile, err := ioutil.TempFile(os.TempDir(), "cornstone")
	if err != nil {
		return err
	}
	tempFilePath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFilePath)

	request, err := grab.NewRequest(tempFilePath, downloadUrl)
	if err != nil {
		return err
	}
	if err := NewMultiDownloader(1, request).Do(); err != nil {
		return err
	}
	if err := ExtractArchiveFromFile(walker, ExtractFileConfig{
		ArchivePath: tempFilePath,
		Common:      config,
	}); err != nil {
		return err
	}
	return nil
}
