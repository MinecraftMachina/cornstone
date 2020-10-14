package initialize

import (
	"cornstone/multimc"
	"cornstone/util"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
)

var destPath string
var profile *multimc.OSProfile
var dev bool

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "Download and initialize MultiMC",
	PreRun: func(cmd *cobra.Command, args []string) {
		destPath = viper.GetString("multimcPath")
		profile = viper.Get("profile").(*multimc.OSProfile)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := execute(); err != nil {
			log.Fatal(err)
		}
	},
}

func validateMultiMCPath() error {
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		if err := os.MkdirAll(destPath, 755); err != nil {
			return nil
		}
	} else if err != nil {
		return err
	} else {
		files, err := ioutil.ReadDir(destPath)
		if err != nil {
			return err
		}
		if len(files) > 0 {
			return errors.New("input directory not empty")
		}
	}
	return nil
}

func execute() error {
	if err := validateMultiMCPath(); err != nil {
		return err
	}

	var downloadUrl string
	if dev {
		downloadUrl = profile.DownloadDevUrl
	} else {
		downloadUrl = profile.DownloadUrl
	}

	tempFile, err := ioutil.TempFile(os.TempDir(), "multimc")
	if err != nil {
		return err
	}
	tempFilePath := tempFile.Name()
	tempFile.Close()
	defer func() {
		os.Remove(tempFilePath)
	}()

	if err := util.DownloadFileWithProgress("MultiMC", tempFilePath, downloadUrl); err != nil {
		return err
	}

	log.Println("Extracting MultiMC...")
	if err := util.ExtractArchiveFromFile(profile.NewWalker(), util.ExtractFileConfig{
		FilePath: tempFilePath,
		Common: util.ExtractCommonConfig{
			BasePath:   profile.BasePath,
			TargetPath: destPath,
			Unwrap:     false,
		},
	}); err != nil {
		return err
	}

	if runtime.GOOS == "darwin" {
		log.Println("Applying post-fixes...")
		// Remove all files from quarantine or user will have to approve each bin and lib individually
		cmd := exec.Command("sh", "-c", fmt.Sprintf("sudo xattr -r -d com.apple.quarantine \"%s\"", destPath))
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	log.Println("Done!")
	return nil
}

func init() {
	Cmd.Flags().BoolVar(&dev, "dev", false, "Download the development version instead of stable")
}
