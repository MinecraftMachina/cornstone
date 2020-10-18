package multimc

import (
	"cornstone/aliases/e"
	"cornstone/cmd/multimc/initialize"
	"cornstone/cmd/multimc/install"
	"cornstone/cmd/multimc/dev"
	"cornstone/cmd/multimc/run"
	"cornstone/multimc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var multimcPath string

var Cmd = &cobra.Command{
	Use:   "multimc",
	Short: "Operate on MultiMC installations",
}

func init() {
	Cmd.PersistentFlags().StringVarP(&multimcPath, "multimc-path", "m", "", "Path to MultiMC root directory")
	if err := Cmd.MarkPersistentFlagRequired("multimc-path"); err != nil {
		log.Fatalln(e.P(err))
	}
	if err := viper.BindPFlag("multimcPath", Cmd.PersistentFlags().Lookup("multimc-path")); err != nil {
		log.Fatalln(e.P(err))
	}
	profile, err := multimc.GetOSProfile()
	if err != nil {
		log.Fatalln(e.P(err))
	}
	viper.Set("profile", profile)

	Cmd.AddCommand(dev.Cmd)
	Cmd.AddCommand(initialize.Cmd)
	Cmd.AddCommand(install.Cmd)
	Cmd.AddCommand(run.Cmd)
}
