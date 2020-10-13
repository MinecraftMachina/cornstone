package curseforge

type Manifest struct {
	Minecraft       Minecraft `json:"minecraft"`
	ManifestType    string    `json:"manifestType"`
	ManifestVersion int       `json:"manifestVersion"`
	Name            string    `json:"name"`
	Version         string    `json:"version"`
	Author          string    `json:"author"`
	ProjectID       int       `json:"projectID"`
	Files           []File    `json:"files"`
	Overrides       string    `json:"overrides"`
}
type ModLoader struct {
	ID      string `json:"id"`
	Primary bool   `json:"primary"`
}
type Minecraft struct {
	Version    string      `json:"version"`
	ModLoaders []ModLoader `json:"modLoaders"`
}
type File struct {
	ProjectID int  `json:"projectID"`
	FileID    int  `json:"fileID"`
	Required  bool `json:"required"`
}
