package convert

import (
	"context"
	"cornstone/curseforge"
	"cornstone/util"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
)

var manifestInput string
var manifestOutput string
var force bool
var concurrentCount int

var Cmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert a Twitch modpack manifest to a Corn manifest (backwards-compatible)",
	PreRun: func(cmd *cobra.Command, args []string) {
		manifestInput = viper.GetString("manifestInput")
		concurrentCount = viper.GetInt("concurrentCount")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := execute(); err != nil {
			log.Fatal(err)
		}
	},
}

func execute() error {
	if _, err := os.Stat(manifestOutput); err == nil || !os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("output file %s already exists or access denied", manifestOutput))
	}
	log.Println("Loading manifest...")
	manifestBytes, err := ioutil.ReadFile(manifestInput)
	if err != nil {
		return err
	}
	manifest := curseforge.CornManifest{}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return err
	}

	var source []interface{}
	for i := range manifest.Files {
		source = append(source, &manifest.Files[i])
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	throttler := util.NewThrottler(util.ThrottlerConfig{
		Ctx:          ctx,
		ResultBuffer: 10,
		Workers:      concurrentCount,
		Source:       source,
		Operation: func(sourceItem interface{}) (interface{}, error) {
			file := sourceItem.(*curseforge.CornFile)
			if file.CornMetadata.Name != "" && !force {
				return nil, nil
			}
			addon, err := curseforge.QueryAddon(file.ProjectID)
			if err != nil {
				return nil, err
			}
			file.CornMetadata = curseforge.CornMetadata{
				Name:       addon.Name,
				Summary:    addon.Summary,
				WebsiteURL: addon.WebsiteURL,
			}
			return nil, nil
		},
	})

	log.Println("Querying addons...")
	for result := range throttler.Run() {
		if result.Error != nil {
			cancelFunc()
			return result.Error
		}
	}

	cornManifestBytes, err := util.JsonMarshalPretty(manifest)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(manifestOutput, cornManifestBytes, 644); err != nil {
		return err
	}
	log.Println("Done! Saved to: " + manifestOutput)
	return nil
}

func init() {
	Cmd.Flags().StringVarP(&manifestOutput, "output", "o", "manifest.new.json", "Path to output manifest.json")
	Cmd.Flags().BoolVarP(&force, "force", "f", false, "Force convert files even if they are already converted")
}
