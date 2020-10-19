package install

import (
	"cornstone/aliases/e"
	"cornstone/curseforge"
	"cornstone/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var destPath string
var input string
var concurrentCount int

var Cmd = &cobra.Command{
	Use:   "install",
	Short: "Install a Corn or Twitch modpack from file or URL into a server",
	PreRun: func(cmd *cobra.Command, args []string) {
		concurrentCount = viper.GetInt("concurrentCount")
		destPath = viper.GetString("serverPath")
	},
	Run: func(cmd *cobra.Command, args []string) {
		util.EnsureDirectoryExists(destPath, false, true)
		if err := execute(); err != nil {
			os.RemoveAll(destPath)
			log.Fatalln(e.P(err))
		}
	},
}

func execute() error {
	if err := curseforge.NewModpackInstaller(&curseforge.ModpackInstallerConfig{
		DestPath:        destPath,
		Input:           input,
		ConcurrentCount: concurrentCount,
		TargetType:      curseforge.TargetServer,
	}).Install(); err != nil {
		return err
	}
	destPathAbs, err := filepath.Abs(destPath)
	if err != nil {
		return err
	}
	log.Println("Done! Saved to:", destPathAbs)
	return nil
}

func init() {
	Cmd.Flags().StringVarP(&input, "input", "i", "", "File path or URL to corn-manifest modpack")
	if err := Cmd.MarkFlagRequired("input"); err != nil {
		log.Fatalln(e.P(err))
	}
}
