package util

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ExtractConfig struct {
	Data       *bytes.Reader
	BasePath   string
	TargetPath string
	Unwrap     bool
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

//ExtractArchive: All files not a child of basePath will be skipped.
func ExtractArchive(reader archiver.Reader, config ExtractConfig) error {
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
		fullName, err := getFileFullNameFromHeader(file)
		if err != nil {
			return err
		}
		fullName, err = filepath.Rel(config.BasePath, fullName)
		if err != nil {
			continue
		}
		// skip files outside of basePath
		if strings.HasPrefix(fullName, "..") {
			continue
		}
		if config.Unwrap {
			firstPathIndex := strings.Index(fullName, string(filepath.Separator))
			if firstPathIndex == -1 {
				continue
			} else {
				fullName = fullName[firstPathIndex+1:]
			}
		}
		fullName = filepath.Join(config.TargetPath, fullName)
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
	}
	return nil
}
