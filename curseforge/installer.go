package curseforge

import (
	"context"
	"cornstone/aliases/e"
	"cornstone/multimc"
	"cornstone/util"
	"encoding/json"
	"fmt"
	"github.com/cavaliercoder/grab"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type ModpackInstaller struct {
	*ModpackInstallerConfig
}

type ModpackInstallerConfig struct {
	DestPath        string
	Input           string
	Unwrap          bool
	ConcurrentCount int
	TargetType      int
}

const (
	TargetMultiMC = iota
	TargetServer  = iota
)

func NewModpackInstaller(config *ModpackInstallerConfig) *ModpackInstaller {
	return &ModpackInstaller{config}
}

// ref: https://github.com/MultiMC/MultiMC5/blob/develop/api/logic/InstanceImportTask.cpp
func (i *ModpackInstaller) Install() error {
	if err := os.MkdirAll(i.DestPath, 0777); err != nil {
		return e.S(err)
	}

	if err := i.stageModpack(i.DestPath); err != nil {
		return e.S(err)
	}

	manifestFile := filepath.Join(i.DestPath, "manifest.json")
	manifestBytes, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return e.S(err)
	}
	manifest := CornManifest{}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return e.S(err)
	}
	if i.TargetType != TargetServer {
		if err := i.processOverrides(&manifest, i.DestPath); err != nil {
			return e.S(err)
		}
	}
	if i.TargetType == TargetMultiMC {
		if err := i.ensureInstanceConfig(&manifest, i.DestPath); err != nil {
			return e.S(err)
		}
		if err := i.ensurePackFile(&manifest, i.DestPath); err != nil {
			return e.S(err)
		}
	}
	if i.TargetType == TargetServer {
		if err := i.processForgeServer(&manifest, i.DestPath); err != nil {
			return err
		}
	}
	var modsDestPath string
	if i.TargetType == TargetServer {
		modsDestPath = i.DestPath
	} else {
		modsDestPath = filepath.Join(i.DestPath, "minecraft")
	}
	if err := i.processMods(&manifest, modsDestPath); err != nil {
		return e.S(err)
	}

	log.Println("Done!")
	return nil
}

func (i *ModpackInstaller) processForgeServer(manifest *CornManifest, destPath string) error {
	forgeVersion := i.getForgeVersion(manifest)
	if forgeVersion == "" {
		return nil
	}
	minecraftVersion := manifest.Minecraft.Version
	fullVersion := fmt.Sprintf("%s-%s", minecraftVersion, forgeVersion)
	forgeName := fmt.Sprintf("forge-%s-installer.jar", fullVersion)
	downloadUrl := fmt.Sprintf("https://files.minecraftforge.net/maven/net/minecraftforge/forge/%s/%s", fullVersion, forgeName)
	savePath := filepath.Join(destPath, forgeName)
	request, err := grab.NewRequest(savePath, downloadUrl)
	if err != nil {
		return err
	}
	log.Println("Downloading Forge installer...")
	if err := util.NewMultiDownloader(i.ConcurrentCount, request).Do(); err != nil {
		return err
	}
	log.Println("Installing Forge...")
	cmd := exec.Command("java", "-jar", forgeName, "-installServer")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = destPath
	if err := cmd.Run(); err != nil {
		return err
	}

	removeList := []string{savePath, savePath + ".log"}
	for _, removeItem := range removeList {
		if err := os.Remove(removeItem); err != nil {
			return err
		}
	}
	return nil
}

func (i *ModpackInstaller) processMods(manifest *CornManifest, destPath string) error {
	modsPath := filepath.Join(destPath, "mods")
	if err := os.MkdirAll(modsPath, 0777); err != nil {
		return err
	}

	var source []interface{}
	for i := range manifest.Files {
		source = append(source, &manifest.Files[i])
	}

	type OpResult struct {
		file        *CornFile
		downloadUrl string
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	throttler := util.NewThrottler(util.ThrottlerConfig{
		Ctx:          ctx,
		ResultBuffer: 10,
		Workers:      i.ConcurrentCount,
		Source:       source,
		Operation: func(sourceItem interface{}) (interface{}, error) {
			file := sourceItem.(*CornFile)
			url, err := GetAddonFileDownloadUrl(file.ProjectID, file.FileID)
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
		downloadPath := util.SafeJoin(destPath, file.InstallPath)
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
	if err := util.NewMultiDownloader(i.ConcurrentCount, requests...).Do(); err != nil {
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

func (i *ModpackInstaller) ensurePackFile(manifest *CornManifest, stagingPath string) error {
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

	pack.Components = append(pack.Components, multimc.Component{
		Important: true,
		UID:       "net.minecraft",
		Version:   manifest.Minecraft.Version,
	})

	forgeVersion := i.getForgeVersion(manifest)
	if forgeVersion != "" {
		pack.Components = append(pack.Components, multimc.Component{
			UID:     "net.minecraftforge",
			Version: forgeVersion,
		})
	}

	packBytes, err := util.JsonMarshalPretty(pack)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(packPath, packBytes, 0666); err != nil {
		return err
	}
	return nil
}

func (i *ModpackInstaller) getForgeVersion(manifest *CornManifest) string {
	for _, loader := range manifest.Minecraft.ModLoaders {
		if strings.HasPrefix(loader.ID, "forge-") {
			return strings.TrimPrefix(loader.ID, "forge-")
		}
	}
	return ""
}

func (i *ModpackInstaller) processOverrides(manifest *CornManifest, stagingPath string) error {
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

func (i *ModpackInstaller) ensureInstanceConfig(manifest *CornManifest, stagingPath string) error {
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

func (i *ModpackInstaller) stageModpack(stagingPath string) error {
	if _, err := os.Stat(i.Input); err != nil {
		log.Println("Obtaining modpack...")
		logger := log.New(os.Stderr, "", log.LstdFlags)
		if err := util.DownloadAndExtract(i.Input, logger, util.ExtractCommonConfig{
			BasePath: "",
			DestPath: stagingPath,
			Unwrap:   i.Unwrap,
		}); err != nil {
			return e.S(err)
		}
	} else {
		log.Println("Extracting modpack...")
		if err := util.ExtractArchiveFromFile(util.ExtractFileConfig{
			ArchivePath: i.Input,
			Common: util.ExtractCommonConfig{
				BasePath: "",
				DestPath: stagingPath,
				Unwrap:   i.Unwrap,
			},
		}); err != nil {
			return e.S(err)
		}
	}
	return nil
}
