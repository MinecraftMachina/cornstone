package install

import (
	"context"
	"cornstone/aliases/e"
	"cornstone/curseforge"
	"cornstone/multimc"
	"cornstone/util"
	"encoding/json"
	"github.com/cavaliercoder/grab"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var destPath string
var profile *multimc.OSProfile
var input string
var name string
var unwrap bool
var concurrentCount int

var Cmd = &cobra.Command{
	Use:   "install",
	Short: "Install a Corn or Twitch modpack from file or URL into an existing MultiMC",
	PreRun: func(cmd *cobra.Command, args []string) {
		concurrentCount = viper.GetInt("concurrentCount")
		destPath = viper.GetString("multimcPath")
		profile = viper.Get("profile").(*multimc.OSProfile)
		destPath = filepath.Join(destPath, filepath.Dir(profile.BinaryPath))
	},
	Run: func(cmd *cobra.Command, args []string) {
		util.EnsureDirectoryExists(destPath, "MultiMC")
		if err := execute(); err != nil {
			log.Fatalln(e.P(err))
		}
	},
}

var forgeMap = map[string]string{
	"1.2.5": "3.4.9.171",
	"1.4.2": "6.0.1.355",
	"1.4.7": "6.6.2.534",
	"1.5.2": "7.8.1.737",
}

// ref: https://github.com/MultiMC/MultiMC5/blob/develop/api/logic/InstanceImportTask.cpp
func execute() error {
	instancePath := filepath.Join(destPath, "instances", name)
	if err := os.MkdirAll(instancePath, 0777); err != nil {
		return e.S(err)
	}

	if err := stageModpack(instancePath); err != nil {
		return e.S(err)
	}

	manifestFile := filepath.Join(instancePath, "manifest.json")
	manifestBytes, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return e.S(err)
	}
	manifest := curseforge.CornManifest{}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return e.S(err)
	}
	if err := processOverrides(&manifest, instancePath); err != nil {
		return e.S(err)
	}
	if err := createInstanceConfig(&manifest, instancePath); err != nil {
		return e.S(err)
	}
	if err := createPackFile(&manifest, instancePath); err != nil {
		return e.S(err)
	}
	if err := downloadMods(&manifest, instancePath); err != nil {
		return e.S(err)
	}

	log.Println("Done!")
	return nil
}

func downloadMods(manifest *curseforge.CornManifest, stagingPath string) error {
	minecraftPath := filepath.Join(stagingPath, "minecraft")
	modsPath := filepath.Join(stagingPath, "minecraft", "mods")
	if err := os.MkdirAll(modsPath, 0777); err != nil {
		return err
	}

	var source []interface{}
	for i := range manifest.Files {
		source = append(source, &manifest.Files[i])
	}

	type OpResult struct {
		file        *curseforge.CornFile
		downloadUrl string
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	throttler := util.NewThrottler(util.ThrottlerConfig{
		Ctx:          ctx,
		ResultBuffer: 10,
		Workers:      concurrentCount,
		Source:       source,
		Operation: func(sourceItem interface{}) (interface{}, error) {
			file := sourceItem.(*curseforge.CornFile)
			url, err := curseforge.GetAddonFileDownloadUrl(file.ProjectID, file.FileID)
			return OpResult{file, url}, err
		},
	})

	log.Println("Building addon download URLs...")
	downloadPaths := map[string]bool{}
	var requests []*grab.Request
	for result := range throttler.Run() {
		if result.Error != nil {
			cancelFunc()
			return result.Error
		}
		opResult := result.Data.(OpResult)

		downloadPath := util.SafeJoin(modsPath, path.Base(opResult.downloadUrl))
		if !opResult.file.Required {
			downloadPath += ".disabled"
		}
		downloadPaths[downloadPath] = true
		request, err := grab.NewRequest(downloadPath, opResult.downloadUrl)
		if err != nil {
			return err
		}
		requests = append(requests, request)
	}

	for _, file := range manifest.ExternalFiles {
		downloadPath := util.SafeJoin(minecraftPath, file.InstallPath)
		if !file.Required {
			downloadPath += ".disabled"
		}
		downloadPaths[downloadPath] = true
		request, err := grab.NewRequest(downloadPath, file.Url)
		if err != nil {
			return err
		}
		requests = append(requests, request)
	}

	log.Println("Downloading files...")
	if err := util.NewMultiDownloader(concurrentCount, requests...).Do(); err != nil {
		return err
	}

	log.Println("Removing old mods...")
	if err := filepath.Walk(modsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(modsPath, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if _, ok := downloadPaths[path]; !ok {
			log.Println(path)
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func createPackFile(manifest *curseforge.CornManifest, stagingPath string) error {
	packPath := filepath.Join(stagingPath, "mmc-pack.json")
	if _, err := os.Stat(packPath); err == nil {
		// if mmc-pack.json already exists in modpack, use it
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	pack := multimc.Pack{
		Components:    []multimc.Component{},
		FormatVersion: 1,
	}

	var forgeVersion string
	for _, loader := range manifest.Minecraft.ModLoaders {
		if strings.HasPrefix(loader.ID, "forge-") {
			forgeVersion = strings.TrimPrefix(loader.ID, "forge-")
		} else {
			log.Printf("Unknown mod loader in manifest: %s\n", loader.ID)
		}
	}

	mcVersion := manifest.Minecraft.Version

	pack.Components = append(pack.Components, multimc.Component{
		Important: true,
		UID:       "net.minecraft",
		Version:   mcVersion,
	})

	if forgeVersion != "" {
		// TODO: Proper resolution
		if forgeVersion == "recommended" {
			forgeVersion = forgeMap[mcVersion]
		}
	}
	pack.Components = append(pack.Components, multimc.Component{
		UID:     "net.minecraftforge",
		Version: forgeVersion,
	})

	packBytes, err := util.JsonMarshalPretty(pack)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(packPath, packBytes, 0666); err != nil {
		return err
	}
	return nil
}

func processOverrides(manifest *curseforge.CornManifest, stagingPath string) error {
	if manifest.Overrides != "" {
		overridePath := util.SafeJoin(stagingPath, manifest.Overrides)
		minecraftPath := filepath.Join(stagingPath, "minecraft")
		if _, err := os.Stat(minecraftPath); err == nil {
			return util.MergePaths(overridePath, minecraftPath)
		} else if os.IsNotExist(err) {
			return os.Rename(overridePath, minecraftPath)
		} else {
			return err
		}
	}
	return nil
}

func createInstanceConfig(manifest *curseforge.CornManifest, stagingPath string) error {
	instanceConfigPath := filepath.Join(stagingPath, "instance.cfg")
	if _, err := os.Stat(instanceConfigPath); err == nil {
		// if instance.cfg already exists in modpack, use it
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	instanceConfig, err := multimc.GenerateInstanceConfig(&multimc.InstanceConfigData{Name: manifest.Name})
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(instanceConfigPath, []byte(instanceConfig), 0666); err != nil {
		return err
	}
	return nil
}

func stageModpack(stagingPath string) error {
	if _, err := os.Stat(input); err != nil {
		log.Println("Downloading modpack...")
		eventChan := make(chan bool)
		go func() {
			<-eventChan
			log.Println("Extracting modpack...")
		}()
		if err := util.DownloadAndExtract(input, eventChan, util.ExtractCommonConfig{
			BasePath: "",
			DestPath: stagingPath,
			Unwrap:   unwrap,
		}); err != nil {
			return e.S(err)
		}
	} else {
		log.Println("Extracting modpack...")
		if err := util.ExtractArchiveFromFile(util.ExtractFileConfig{
			ArchivePath: input,
			Common: util.ExtractCommonConfig{
				BasePath: "",
				DestPath: stagingPath,
				Unwrap:   unwrap,
			},
		}); err != nil {
			return e.S(err)
		}
	}
	return nil
}

func init() {
	Cmd.Flags().StringVarP(&input, "input", "i", "", "File path or URL to corn-manifest modpack")
	if err := Cmd.MarkFlagRequired("input"); err != nil {
		log.Fatalln(e.P(err))
	}
	Cmd.Flags().BoolVarP(&unwrap, "unwrap", "u", false, "Discard the root directory of the archive wwhen extracting")
	Cmd.Flags().StringVarP(&name, "name", "n", "", "Name to use for modpack when importing to MultiMC")
	if err := Cmd.MarkFlagRequired("name"); err != nil {
		log.Fatalln(e.P(err))
	}
}
