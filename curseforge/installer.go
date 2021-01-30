package curseforge

import (
	"context"
	"cornstone/aliases/e"
	"cornstone/archive"
	"cornstone/multimc"
	"cornstone/util"
	"encoding/json"
	"fmt"
	"github.com/ViRb3/go-parallel/downloader"
	"github.com/ViRb3/go-parallel/throttler"
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
	tempStagingPath, err := util.TempDir()
	if err != nil {
		return e.S(err)
	}
	defer os.RemoveAll(tempStagingPath)

	if err := i.stageModpack(tempStagingPath); err != nil {
		return e.S(err)
	}
	if err := i.fixWrappedModpack(tempStagingPath); err != nil {
		return e.S(err)
	}

	manifestFile := filepath.Join(tempStagingPath, "manifest.json")
	manifestBytes, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return e.S(err)
	}
	manifest := CornManifest{}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return e.S(err)
	}
	if err := i.processOverrides(&manifest, tempStagingPath); err != nil {
		return e.S(err)
	}
	if i.TargetType == TargetMultiMC {
		if err := i.ensureInstanceConfig(&manifest, tempStagingPath); err != nil {
			return e.S(err)
		}
		if err := i.ensurePackFile(&manifest, tempStagingPath); err != nil {
			return e.S(err)
		}
	}
	if i.TargetType == TargetServer {
		if err := i.processForgeServer(&manifest, tempStagingPath); err != nil {
			return e.S(err)
		}
	}

	if err := os.MkdirAll(i.DestPath, 0777); err != nil {
		return e.S(err)
	}
	if err := util.MergePaths(tempStagingPath, i.DestPath); err != nil {
		return e.S(err)
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
	return nil
}

// If the modpack content is wrapped in one root directory, unwrap it.
func (i *ModpackInstaller) fixWrappedModpack(modpackPath string) error {
	manifestFile := filepath.Join(modpackPath, "manifest.json")
	if _, err := os.Stat(manifestFile); os.IsNotExist(err) {
		subFiles, err := ioutil.ReadDir(modpackPath)
		if err != nil {
			return e.S(err)
		}
		if len(subFiles) == 1 && subFiles[0].IsDir() {
			rootDirPath := filepath.Join(modpackPath, subFiles[0].Name())
			if err := util.MergePaths(rootDirPath, modpackPath); err != nil {
				return e.S(err)
			}
		}
	} else if err != nil {
		return e.S(err)
	}
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

	log.Println("Downloading Forge installer...")
	job := downloader.Job{SaveFilePath: savePath, Url: downloadUrl}
	result, cancelFunc := util.MultiDownload(downloader.SharedConfig{
		ShowProgress:   true,
		SkipSameLength: true,
		Workers:        i.ConcurrentCount,
		Jobs:           []downloader.Job{job},
	})
	defer cancelFunc()
	for resp := range result {
		if err := resp.Err; err != nil {
			return err
		}
	}

	log.Println("Installing Forge...")
	cmd := exec.Command("java", "-jar", forgeName, "-installServer")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = destPath
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := os.Remove(savePath); err != nil {
		return err
	}

	removeList := []string{
		savePath + ".log",
		filepath.Join(filepath.Dir(savePath), "installer.log"),
	}
	for _, removeItem := range removeList {
		// depending on version not all files will exist, so ignore errors
		os.Remove(removeItem)
	}
	return nil
}

func (i *ModpackInstaller) processMods(manifest *CornManifest, destPath string) error {
	modsPath := filepath.Join(destPath, "mods")
	if err := os.MkdirAll(modsPath, 0777); err != nil {
		return err
	}

	var source []interface{}
	for i2 := range manifest.Files {
		file := &manifest.Files[i2]
		if i.TargetType == TargetServer && file.ServerIgnored {
			continue
		}
		source = append(source, file)
	}

	type OpResult struct {
		file        *CornFile
		downloadUrl string
		Err         error
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	addonThrottler := throttler.NewThrottler(throttler.Config{
		Ctx:          ctx,
		ResultBuffer: 10,
		Workers:      i.ConcurrentCount,
		Source:       source,
		Operation: func(sourceItem interface{}) interface{} {
			file := sourceItem.(*CornFile)
			url, err := GetAddonFileDownloadUrl(file.ProjectID, file.FileID)
			return OpResult{file, url, err}
		},
	})

	log.Println("Building addon download URLs...")
	downloadPaths := map[string]bool{}
	var jobs []downloader.Job
	for result := range addonThrottler.Run() {
		opResult := result.(OpResult)
		if opResult.Err != nil {
			cancelFunc()
			return opResult.Err
		}

		downloadPath := util.SafeJoin(modsPath, path.Base(opResult.downloadUrl))
		if !opResult.file.Required {
			downloadPath += ".disabled"
		}
		downloadPaths[downloadPath] = true
		request := downloader.Job{SaveFilePath: downloadPath, Url: opResult.downloadUrl}
		jobs = append(jobs, request)
	}
	cancelFunc()

	for _, file := range manifest.ExternalFiles {
		if i.TargetType == TargetServer && file.ServerIgnored {
			continue
		}
		var downloadPath string
		if file.Extract.Enable {
			tempFile, err := util.TempFile()
			if err != nil {
				return err
			}
			tempFilePath := tempFile.Name()
			if err := tempFile.Close(); err != nil {
				return err
			}
			defer os.Remove(tempFilePath)
			downloadPath = tempFilePath
		} else {
			downloadPath = util.SafeJoin(destPath, file.InstallPath)
			if !file.Required {
				downloadPath += ".disabled"
			}
			downloadPaths[downloadPath] = true
		}
		request := downloader.Job{SaveFilePath: downloadPath, Url: file.Url, Tag: file}
		jobs = append(jobs, request)
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

	log.Println("Obtaining files...")
	results, cancelFunc := util.MultiDownload(downloader.SharedConfig{
		ShowProgress:   true,
		SkipSameLength: true,
		Workers:        i.ConcurrentCount,
		Jobs:           jobs,
	})
	defer cancelFunc()
	for result := range results {
		job := result.Job
		if err := result.Err; err != nil {
			if file, ok := job.Tag.(ExternalFile); ok {
				log.Println("error downloading external file: " + file.Name)
				if file.Required {
					return err
				}
			} else if file, ok := job.Tag.(*CornFile); ok {
				log.Printf("error downloading addon: %d\n", file.ProjectID)
				if file.Required {
					return err
				}
			} else {
				return err
			}
		}
		if file, ok := job.Tag.(ExternalFile); ok && file.Extract.Enable && file.Required {
			extractPath := util.SafeJoin(destPath, file.InstallPath)
			log.Printf("Extracting external file '%s'...\n", file.Name)
			if err := archive.ExtractArchiveFromFile(archive.ExtractFileConfig{
				ArchivePath: result.Job.SaveFilePath,
				Common: archive.ExtractCommonConfig{
					BasePath: "",
					DestPath: extractPath,
					Unwrap:   file.Extract.Unwrap,
				},
			}); err != nil {
				return err
			}
			return nil
		}
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
		var minecraftPath string
		if i.TargetType == TargetServer {
			minecraftPath = stagingPath
		} else {
			minecraftPath = filepath.Join(stagingPath, "minecraft")
		}
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
		if err := archive.DownloadAndExtract(i.Input, logger, archive.ExtractCommonConfig{
			BasePath: "",
			DestPath: stagingPath,
			Unwrap:   false,
		}); err != nil {
			return err
		}
	} else {
		log.Println("Extracting modpack...")
		if err := archive.ExtractArchiveFromFile(archive.ExtractFileConfig{
			ArchivePath: i.Input,
			Common: archive.ExtractCommonConfig{
				BasePath: "",
				DestPath: stagingPath,
				Unwrap:   false,
			},
		}); err != nil {
			return err
		}
	}
	return nil
}
