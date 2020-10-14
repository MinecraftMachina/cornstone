package multimc

import (
	"cornstone/cmd/multimc/init"
	"cornstone/cmd/multimc/install"
	"cornstone/cmd/multimc/dev"
	"cornstone/cmd/multimc/run"
	"cornstone/multimc"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var multimcPath string

var Cmd = &cobra.Command{
	Use:   "multimc",
	Short: "Configures MultiMC",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		validateMultiMCPath()
	},
}

func validateMultiMCPath() {
	if multimcPath == "" {
		return // let cobra prompt for the flag
	}
	stat, err := os.Stat(multimcPath)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		log.Fatal("supplied path is not a directory")
	}
}

func init() {
	Cmd.PersistentFlags().StringVarP(&multimcPath, "multimc-path", "m", "", "Path to MultiMC root directory")
	if err := Cmd.MarkPersistentFlagRequired("multimc-path"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("multimcPath", Cmd.PersistentFlags().Lookup("multimc-path")); err != nil {
		log.Fatal(err)
	}
	profile, err := multimc.GetOSProfile()
	if err != nil {
		log.Fatal(err)
	}
	viper.Set("profile", profile)

	Cmd.AddCommand(dev.Cmd)
	Cmd.AddCommand(init.Cmd)
	Cmd.AddCommand(install.Cmd)
	Cmd.AddCommand(run.Cmd)
}
