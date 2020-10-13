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
type ExternalFile struct {
	Url         string `json:"url"`
	InstallPath string `json:"installPath"`
	Required    bool   `json:"required"`
}
