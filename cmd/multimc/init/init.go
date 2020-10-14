package init

import (
	"bytes"
	"cornstone/multimc"
	"cornstone/util"
	"fmt"
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

func validateMultiMCPath() {
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		if err := os.MkdirAll(destPath, 755); err != nil {
			return
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		files, err := ioutil.ReadDir(destPath)
		if err != nil {
			log.Fatal(err)
		}
		if len(files) > 0 {
			log.Fatal("input directory not empty")
		}
	}
}

func execute() error {
	validateMultiMCPath()

	var downloadUrl string
	if dev {
		downloadUrl = profile.DownloadDevUrl
	} else {
		downloadUrl = profile.DownloadUrl
	}
	log.Println("Downloading MultiMC...")
	bar := util.NewBar(1)
	var data = make([]byte, 0)
	if _, err := util.DefaultClient.New().Get(downloadUrl).ByteResponse().ReceiveSuccess(&data); err != nil {
		return err
	}
	bar.Add(1)

	log.Println("Extracting MultiMC...")
	if err := util.ExtractArchive(profile.NewReader(), util.ExtractConfig{
		Data:       bytes.NewReader(data),
		BasePath:   profile.BasePath,
		TargetPath: destPath,
		Unwrap:     false,
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
