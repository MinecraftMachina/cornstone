package multimc

import (
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
	"runtime"
)

type OSProfile struct {
	BasePath    string
	BinaryPath  string
	DownloadUrl string
	NewReader   func() archiver.Reader
}

var osProfiles = map[string]OSProfile{
	"windows": {
		BasePath:    "MultiMC/",
		BinaryPath:  "MultiMC.exe",
		DownloadUrl: "https://files.multimc.org/downloads/mmc-stable-win32.zip",
		NewReader: func() archiver.Reader {
			return archiver.NewZip()
		},
	}, "linux": {
		BasePath:    "MultiMC/",
		BinaryPath:  "MultiMC",
		DownloadUrl: "https://files.multimc.org/downloads/mmc-stable-lin64.tar.gz",
		NewReader: func() archiver.Reader {
			return archiver.NewTarGz()
		},
	},
	"darwin": {
		BasePath:    "MultiMC.app/",
		BinaryPath:  "Contents/MacOS/MultiMC",
		DownloadUrl: "https://files.multimc.org/downloads/mmc-stable-osx64.tar.gz",
		NewReader: func() archiver.Reader {
			return archiver.NewTarGz()
		},
	},
}

func GetOSProfile() (*OSProfile, error) {
	profile, ok := osProfiles[runtime.GOOS]
	if !ok {
		return &OSProfile{}, errors.New("unsupported operating system")
	}
	return &profile, nil
}
