package initialize

import (
	"cornstone/aliases/e"
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
	"path/filepath"
	"runtime"
)

var multimcPath string
var destPath string
var profile *multimc.OSProfile
var dev bool
var analytics bool
var noJava bool

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "Download and initialize MultiMC with Java bundled",
	PreRun: func(cmd *cobra.Command, args []string) {
		multimcPath = viper.GetString("multimcPath")
		profile = viper.Get("profile").(*multimc.OSProfile)
		destPath = filepath.Join(multimcPath, filepath.Dir(profile.BinaryPath))
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateMultiMCPath(); err != nil {
			log.Fatalln(e.P(err))
		}
		if err := execute(); err != nil {
			os.RemoveAll(multimcPath)
			log.Fatalln(e.P(err))
		}
	},
}

func validateMultiMCPath() error {
	if _, err := os.Stat(multimcPath); os.IsNotExist(err) {
		if err := os.MkdirAll(multimcPath, 0777); err != nil {
			return nil
		}
	} else if err != nil {
		return err
	} else {
		files, err := ioutil.ReadDir(multimcPath)
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
	var downloadUrl string
	if dev {
		downloadUrl = profile.DownloadDevUrl
	} else {
		downloadUrl = profile.DownloadUrl
	}

	log.Println("Downloading MultiMC...")
	eventChan := make(chan bool)
	go func() {
		<-eventChan
		log.Println("Extracting MultiMC...")
	}()
	if err := util.DownloadAndExtract(downloadUrl, eventChan, util.ExtractCommonConfig{
		BasePath: "",
		DestPath: multimcPath,
		Unwrap:   true,
	}); err != nil {
		return e.S(err)
	}

	if !noJava {
		log.Println("Downloading Java...")
		javaPath := filepath.Join(destPath, "java")
		if err := os.MkdirAll(javaPath, 0777); err != nil {
			return e.S(err)
		}
		eventChan := make(chan bool)
		go func() {
			<-eventChan
			log.Println("Extracting Java...")
		}()
		if err := util.DownloadAndExtract(profile.JavaUrl, eventChan, util.ExtractCommonConfig{
			BasePath: "",
			DestPath: javaPath,
			Unwrap:   true,
		}); err != nil {
			return e.S(err)
		}
	}

	var javaPath string
	if noJava {
		javaPath = "java"
	} else {
		javaPath = filepath.Join("java", profile.JavaBinaryPath)
	}

	log.Println("Configuring MultiMC...")
	hostname, err := os.Hostname()
	if err != nil {
		return e.S(err)
	}
	config, err := multimc.GenerateMainConfig(&multimc.MainConfigData{
		JavaPath:     javaPath,
		Analytics:    analytics,
		LastHostname: hostname,
	})
	if err != nil {
		return e.S(err)
	}
	if err := ioutil.WriteFile(filepath.Join(destPath, "multimc.cfg"), []byte(config), 0666); err != nil {
		return e.S(err)
	}

	if runtime.GOOS == "darwin" {
		log.Println("Applying post-fixes...")
		// Remove all files from quarantine or user will have to approve each bin and lib individually
		cmd := exec.Command("sh", "-c", fmt.Sprintf("sudo xattr -r -d com.apple.quarantine \"%s\"", multimcPath))
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return e.S(err)
		}
	}
	log.Println("Done!")
	return nil
}

func init() {
	Cmd.Flags().BoolVar(&dev, "dev", false, "Download the development version instead of stable")
	Cmd.Flags().BoolVar(&analytics, "analytics", false, "Enable MultiMC analytics")
	Cmd.Flags().BoolVar(&noJava, "no-java", false, "Don't download Java")
}
