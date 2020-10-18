package run

import (
	"cornstone/aliases/e"
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
		multimcPath := viper.GetString("multimcPath")
		profile = viper.Get("profile").(*multimc.OSProfile)
		binaryPath = filepath.Join(multimcPath, profile.BinaryPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		util.EnsureFileExists(binaryPath)
		if err := execute(); err != nil {
			log.Fatalln(e.P(err))
		}
	},
}

func execute() error {
	cmd := exec.Command(binaryPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = filepath.Dir(binaryPath)
	cmd.Path = filepath.Base(binaryPath)
	if runtime.GOOS == "windows" {
		cmd.Env = os.Environ()
		// run with NVIDIA GPU if Optimus is present
		cmd.Env = append(cmd.Env, "SHIM_MCCOMPAT=0x800000001")
	}
	return e.S(cmd.Start())
}
