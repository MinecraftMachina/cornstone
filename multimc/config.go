package multimc

import (
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
	"runtime"
)

type OSProfile struct {
	BasePath       string
	BinaryPath     string
	DownloadUrl    string
	DownloadDevUrl string
	NewWalker      func() archiver.Walker
}

var osProfiles = map[string]OSProfile{
	"windows": {
		BasePath:       "MultiMC/",
		BinaryPath:     "MultiMC.exe",
		DownloadUrl:    "https://files.multimc.org/downloads/mmc-stable-win32.zip",
		DownloadDevUrl: "https://files.multimc.org/downloads/mmc-develop-win32.zip",
		NewWalker: func() archiver.Walker {
			return archiver.NewZip()
		},
	}, "linux": {
		BasePath:       "MultiMC/",
		BinaryPath:     "MultiMC",
		DownloadUrl:    "https://files.multimc.org/downloads/mmc-stable-lin64.tar.gz",
		DownloadDevUrl: "https://files.multimc.org/downloads/mmc-develop-lin64.tar.gz",
		NewWalker: func() archiver.Walker {
			return archiver.NewTarGz()
		},
	},
	"darwin": {
		BasePath:       "MultiMC.app/",
		BinaryPath:     "Contents/MacOS/MultiMC",
		DownloadUrl:    "https://files.multimc.org/downloads/mmc-stable-osx64.tar.gz",
		DownloadDevUrl: "https://files.multimc.org/downloads/mmc-develop-osx64.tar.gz",
		NewWalker: func() archiver.Walker {
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
