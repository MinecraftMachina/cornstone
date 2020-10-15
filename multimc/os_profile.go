package multimc

import (
	"github.com/pkg/errors"
	"runtime"
)

type OSProfile struct {
	BinaryPath     string
	DownloadUrl    string
	DownloadDevUrl string
	JavaUrl        string
	JavaBinaryPath string
}

var osProfiles = map[string]OSProfile{
	"windows": {
		BinaryPath:     "MultiMC.exe",
		DownloadUrl:    "https://files.multimc.org/downloads/mmc-stable-win32.zip",
		DownloadDevUrl: "https://files.multimc.org/downloads/mmc-develop-win32.zip",
		JavaUrl:        "https://corretto.aws/downloads/latest/amazon-corretto-8-x64-windows-jre.zip",
		JavaBinaryPath: "bin/javaw.exe",
	}, "linux": {
		BinaryPath:     "MultiMC",
		DownloadUrl:    "https://files.multimc.org/downloads/mmc-stable-lin64.tar.gz",
		DownloadDevUrl: "https://files.multimc.org/downloads/mmc-develop-lin64.tar.gz",
		JavaUrl:        "https://corretto.aws/downloads/latest/amazon-corretto-8-x64-linux-jdk.tar.gz",
		JavaBinaryPath: "bin/javaw",
	},
	"darwin": {
		BinaryPath:     "Contents/MacOS/MultiMC",
		DownloadUrl:    "https://files.multimc.org/downloads/mmc-stable-osx64.tar.gz",
		DownloadDevUrl: "https://files.multimc.org/downloads/mmc-develop-osx64.tar.gz",
		JavaUrl:        "https://corretto.aws/downloads/latest/amazon-corretto-8-x64-macos-jdk.tar.gz",
		JavaBinaryPath: "Contents/Home/bin/java",
	},
}

func GetOSProfile() (*OSProfile, error) {
	profile, ok := osProfiles[runtime.GOOS]
	if !ok {
		return &OSProfile{}, errors.New("unsupported operating system")
	}
	return &profile, nil
}
