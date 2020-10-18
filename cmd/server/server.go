package server

import (
	"cornstone/aliases/e"
	"cornstone/cmd/server/install"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var serverPath string

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Operate on server installations",
}

func init() {
	Cmd.PersistentFlags().StringVarP(&serverPath, "server-path", "s", "", "Path to server root directory")
	if err := Cmd.MarkPersistentFlagRequired("server-path"); err != nil {
		log.Fatalln(e.P(err))
	}
	if err := viper.BindPFlag("serverPath", Cmd.PersistentFlags().Lookup("server-path")); err != nil {
		log.Fatalln(e.P(err))
	}

	Cmd.AddCommand(install.Cmd)
}
