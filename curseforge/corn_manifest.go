package curseforge

type CornManifest struct {
	Manifest
	Files         []CornFile     `json:"files"`
	ExternalFiles []ExternalFile `json:"_externalFiles"`
}
type CornFile struct {
	File
	CornMetadata
}
type CornMetadata struct {
	Name       string `json:"_name"`
	Summary    string `json:"_summary"`
	WebsiteURL string `json:"_websiteUrl"`
}
type ExtractConfig struct {
	Enable bool `json:"enable"`
	Unwrap bool `json:"unwrap"`
}
type ExternalFile struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	// Can be a file path or directory path.
	// In the case of a directory path, the file name will be inferred.
	// See: grab.Request#Filename
	InstallPath string        `json:"installPath"`
	Required    bool          `json:"required"`
	Extract     ExtractConfig `json:"extract"`
}
