package cmd

import (
	"cornstone/aliases/e"
	"cornstone/cmd/manifest"
	"cornstone/cmd/multimc"
	"cornstone/cmd/server"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var concurrentCount int

var cmd = &cobra.Command{
	Use:   "cornstone",
	Short: "The single utility for all your modded Minecraft needs",
}

func Execute() error {
	return cmd.Execute()
}

func init() {
	cmd.PersistentFlags().IntVarP(&concurrentCount, "concurrent-count", "c", 5, "Concurrent download count")
	if err := viper.BindPFlag("concurrentCount", cmd.PersistentFlags().Lookup("concurrent-count")); err != nil {
		log.Fatalln(e.P(err))
	}

	cmd.AddCommand(manifest.Cmd)
	cmd.AddCommand(multimc.Cmd)
	cmd.AddCommand(server.Cmd)
}
