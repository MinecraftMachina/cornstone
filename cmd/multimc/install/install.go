package install

import (
	"cornstone/aliases/e"
	"cornstone/curseforge"
	"cornstone/multimc"
	"cornstone/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

var destPath string
var profile *multimc.OSProfile
var input string
var name string
var concurrentCount int

var Cmd = &cobra.Command{
	Use:   "install",
	Short: "Install a Corn or Twitch modpack from file or URL into an existing MultiMC",
	PreRun: func(cmd *cobra.Command, args []string) {
		concurrentCount = viper.GetInt("concurrentCount")
		destPath = viper.GetString("multimcPath")
		profile = viper.Get("profile").(*multimc.OSProfile)
		destPath = filepath.Join(destPath, filepath.Dir(profile.BinaryPath))
	},
	Run: func(cmd *cobra.Command, args []string) {
		util.EnsureDirectoryExists(destPath, false, false)
		if err := execute(); err != nil {
			log.Fatalln(e.P(err))
		}
	},
}

func execute() error {
	if err := curseforge.NewModpackInstaller(&curseforge.ModpackInstallerConfig{
		DestPath:        filepath.Join(destPath, "instances", name),
		Input:           input,
		ConcurrentCount: concurrentCount,
		TargetType:      curseforge.TargetMultiMC,
	}).Install(); err != nil {
		return err
	}
	destPathAbs, err := filepath.Abs(destPath)
	if err != nil {
		return err
	}
	log.Println("Done! Saved to: ", destPathAbs)
	return nil
}

func init() {
	Cmd.Flags().StringVarP(&input, "input", "i", "", "File path or URL to corn-manifest modpack")
	if err := Cmd.MarkFlagRequired("input"); err != nil {
		log.Fatalln(e.P(err))
	}
	Cmd.Flags().StringVarP(&name, "name", "n", "", "Name to use for modpack when importing to MultiMC")
	if err := Cmd.MarkFlagRequired("name"); err != nil {
		log.Fatalln(e.P(err))
	}
}
