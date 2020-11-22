package dev

import (
	"bufio"
	"cornstone/aliases/e"
	"cornstone/multimc"
	"cornstone/util"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var playerName string
var destPath string
var profile *multimc.OSProfile
var playerNameRegex = regexp.MustCompile(`^[A-z]{3,16}$`)

var Cmd = &cobra.Command{
	Use:   "dev",
	Short: "Set up MultiMC for development",
	PreRun: func(cmd *cobra.Command, args []string) {
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
	for !playerNameRegex.MatchString(playerName) {
		fmt.Print("Player name (3-16 alphanumeric characters): ")
		reader := bufio.NewReader(os.Stdin)
		var err error
		playerName, err = reader.ReadString('\n')
		if err != nil {
			return e.S(err)
		}
		playerName = strings.TrimRight(playerName, "\r\n")
	}

	accountsJson, err := multimc.MakeNewAccountsJson(playerName)
	if err != nil {
		return e.S(err)
	}
	if err := ioutil.WriteFile(filepath.Join(destPath, "accounts.json"), accountsJson, 0666); err != nil {
		return e.S(err)
	}
	log.Println("Done!")
	return nil
}

func init() {
	Cmd.Flags().StringVarP(&playerName, "player-name", "n", "", "Player name to use")
}
