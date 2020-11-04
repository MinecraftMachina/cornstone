package util

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"github.com/cavaliercoder/grab"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
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
	// All archive files not a child of BasePath will be skipped.
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

type extractor interface {
	archiver.Walker
	archiver.Reader
}

func createExtractor(header []byte) (extractor, error) {
	kind, err := filetype.Match(header)
	if err != nil {
		return nil, err
	}
	switch kind {
	case matchers.TypeZip:
		return archiver.NewZip(), nil
	case matchers.TypeGz:
		return archiver.NewTarGz(), nil // TODO: Don't assume .tar.gz
	case matchers.TypeTar:
		return archiver.NewTar(), nil
	default:
		return nil, errors.New("unsupported archive type: " + kind.Extension)
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
	fullName = SafeJoin(config.DestPath, fullName)
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

func ExtractArchiveFromFile(config ExtractFileConfig) error {
	file, err := os.Open(config.ArchivePath)
	if err != nil {
		return err
	}
	defer file.Close()
	header := make([]byte, 261)
	if _, err := file.Read(header); err != nil {
		return err
	}
	extractor, err := createExtractor(header)
	if err != nil {
		return err
	}
	if err := extractor.Walk(config.ArchivePath, func(file archiver.File) error {
		return processFile(file, config.Common)
	}); err != nil {
		return err
	}
	return nil
}

func ExtractArchiveFromReader(config ExtractReaderConfig) error {
	header := make([]byte, 261)
	if _, err := config.Data.Read(header); err != nil {
		return err
	}
	if _, err := config.Data.Seek(0, io.SeekStart); err != nil {
		return err
	}
	extractor, err := createExtractor(header)
	if err != nil {
		return err
	}
	if err := extractor.Open(config.Data, int64(config.Data.Len())); err != nil {
		return err
	}
	for {
		file, err := extractor.Read()
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

func DownloadAndExtract(downloadUrl string, logger *log.Logger, config ExtractCommonConfig) error {
	tempFile, err := TempFile()
	if err != nil {
		return err
	}
	tempFilePath := tempFile.Name()
	if err := tempFile.Close(); err != nil {
		return err
	}
	defer os.Remove(tempFilePath)

	request, err := grab.NewRequest(tempFilePath, downloadUrl)
	if err != nil {
		return err
	}
	logger.Println("Downloading...")
	result, cancelFunc := NewMultiDownloader(1, request).Do()
	defer cancelFunc()
	for resp := range result {
		if err := resp.Err(); err != nil {
			return err
		}
	}
	logger.Println("Extracting...")
	if err := ExtractArchiveFromFile(ExtractFileConfig{
		ArchivePath: tempFilePath,
		Common:      config,
	}); err != nil {
		return err
	}
	return nil
}
