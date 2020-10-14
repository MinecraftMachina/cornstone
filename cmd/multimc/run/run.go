package run

import (
	"cornstone/multimc"
	"cornstone/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var binaryPath string
var profile *multimc.OSProfile

var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Run an existing MultiMC",
	PreRun: func(cmd *cobra.Command, args []string) {
		binaryPath = viper.GetString("multimcPath")
		profile = viper.Get("profile").(*multimc.OSProfile)
		binaryPath = filepath.Join(binaryPath, profile.BinaryPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := execute(); err != nil {
			log.Fatal(err)
		}
	},
}

func execute() error {
	util.EnsureFileExists(binaryPath, "MultiMC")
	cmd := exec.Command(binaryPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if runtime.GOOS == "windows" {
		cmd.Env = os.Environ()
		// run with NVIDIA GPU if Optimus is present
		cmd.Env = append(cmd.Env, "SHIM_MCCOMPAT=0x800000001")
	}
	return cmd.Start()
}
