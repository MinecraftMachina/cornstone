package manifest

import (
	"cornstone/cmd/manifest/convert"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var manifest string

var Cmd = &cobra.Command{
	Use:   "manifest",
	Short: "Operates on Twitch modpack manifest",
}

func init() {
	Cmd.PersistentFlags().StringVarP(&manifest, "input", "i", "", "Path to input manifest.json")
	if err := Cmd.MarkPersistentFlagRequired("input"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("manifestInput", Cmd.PersistentFlags().Lookup("input")); err != nil {
		log.Fatal(err)
	}

	Cmd.AddCommand(convert.Cmd)
}
