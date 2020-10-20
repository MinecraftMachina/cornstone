package curseforge

type CornManifest struct {
	Manifest
	Files         []CornFile     `json:"files"`
	ExternalFiles []ExternalFile `json:"_externalFiles"`
}
type CornFile struct {
	File
	Metadata      Metadata `json:"_metadata"`
	ServerIgnored bool     `json:"_serverIgnored"`
}
type Metadata struct {
	ProjectName string `json:"projectName"`
	FileName    string `json:"fileName"`
	Summary     string `json:"summary"`
	WebsiteURL  string `json:"websiteUrl"`
	Hash        string `json:"hash"`
}
type ExtractConfig struct {
	Enable bool `json:"enable"`
	Unwrap bool `json:"unwrap"`
}
type ExternalFile struct {
	Name          string        `json:"name"`
	Url           string        `json:"url"`
	InstallPath   string        `json:"installPath"`
	Required      bool          `json:"required"`
	Extract       ExtractConfig `json:"extract"`
	ServerIgnored bool          `json:"serverIgnored"`
}
