package lint

import (
	"context"
	"cornstone/aliases/e"
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
	"path/filepath"
)

var manifestInput string
var manifestOutput string
var force bool
var concurrentCount int

var Cmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint a Corn or Twitch modpack manifest, updating annotations",
	PreRun: func(cmd *cobra.Command, args []string) {
		manifestInput = viper.GetString("manifestInput")
		concurrentCount = viper.GetInt("concurrentCount")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := execute(); err != nil {
			log.Fatalln(e.P(err))
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
		return e.S(err)
	}
	manifest := curseforge.CornManifest{}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return e.S(err)
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
			manifestFile := sourceItem.(*curseforge.CornFile)
			hash := fmt.Sprintf("%d%d", manifestFile.ProjectID, manifestFile.FileID)
			if hash == manifestFile.Metadata.Hash && !force {
				return nil, nil
			}
			addon, err := curseforge.QueryAddon(manifestFile.ProjectID)
			if err != nil {
				return nil, err
			}
			files, err := curseforge.QueryAddonFiles(manifestFile.ProjectID)
			if err != nil {
				return nil, err
			}
			var fileName string
			for _, file := range files {
				if file.ID == manifestFile.FileID {
					fileName = file.FileName
					break
				}
			}
			manifestFile.Metadata = curseforge.CornMetadata{
				ProjectName: addon.Name,
				FileName:    fileName,
				Summary:     addon.Summary,
				WebsiteURL:  addon.WebsiteURL,
				Hash:        hash,
			}
			return nil, nil
		},
	})

	log.Println("Querying addons...")
	for result := range throttler.Run() {
		if result.Error != nil {
			cancelFunc()
			return e.S(result.Error)
		}
	}

	cornManifestBytes, err := util.JsonMarshalPretty(manifest)
	if err != nil {
		return e.S(err)
	}
	if err := ioutil.WriteFile(manifestOutput, cornManifestBytes, 0666); err != nil {
		return e.S(err)
	}

	manifestOutputAbs, err := filepath.Abs(manifestOutput)
	if err != nil {
		return err
	}
	log.Println("Done! Saved to:" + manifestOutputAbs)
	return nil
}

func init() {
	Cmd.Flags().StringVarP(&manifestOutput, "output", "o", "manifest.new.json", "Path to output manifest.json")
	Cmd.Flags().BoolVarP(&force, "force", "f", false, "Force annotate files from server even if they are already annotated")
}
