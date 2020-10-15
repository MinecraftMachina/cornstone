package manifest

import (
	"cornstone/aliases/e"
	"cornstone/cmd/manifest/convert"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var manifest string

var Cmd = &cobra.Command{
	Use:   "manifest",
	Short: "Operate on modpack manifests",
}

func init() {
	Cmd.PersistentFlags().StringVarP(&manifest, "input", "i", "", "Path to input manifest.json")
	if err := Cmd.MarkPersistentFlagRequired("input"); err != nil {
		log.Fatalln(e.P(err))
	}
	if err := viper.BindPFlag("manifestInput", Cmd.PersistentFlags().Lookup("input")); err != nil {
		log.Fatalln(e.P(err))
	}

	Cmd.AddCommand(convert.Cmd)
}
