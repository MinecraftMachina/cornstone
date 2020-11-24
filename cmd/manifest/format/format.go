package format

import (
	"context"
	"cornstone/aliases/e"
	"cornstone/curseforge"
	"cornstone/throttler"
	"cornstone/util"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"sort"
)

var manifestInput string
var force bool
var concurrentCount int

var Cmd = &cobra.Command{
	Use:   "format",
	Short: "Format a Curse/Corn modpack manifest and update its metadata",
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
	addonThrottler := throttler.NewThrottler(throttler.Config{
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
			manifestFile.Metadata = curseforge.Metadata{
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
	for result := range addonThrottler.Run() {
		if result.Error != nil {
			cancelFunc()
			return e.S(result.Error)
		}
	}

	sort.Slice(manifest.Files, func(i, j int) bool {
		return manifest.Files[i].Metadata.ProjectName < manifest.Files[j].Metadata.ProjectName
	})
	sort.Slice(manifest.ExternalFiles, func(i, j int) bool {
		return manifest.ExternalFiles[i].Name < manifest.ExternalFiles[j].Name
	})

	cornManifestBytes, err := util.JsonMarshalPretty(manifest)
	if err != nil {
		return e.S(err)
	}
	if err := ioutil.WriteFile(manifestInput, cornManifestBytes, 0666); err != nil {
		return e.S(err)
	}

	log.Println("Done!")
	return nil
}

func init() {
	Cmd.Flags().BoolVarP(&force, "force", "f", false, "Force annotate files from server even if they are already annotated")
}
