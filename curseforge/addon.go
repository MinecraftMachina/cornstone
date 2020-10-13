package curseforge

import "time"

type Addon struct {
	ID                     int                      `json:"id"`
	Name                   string                   `json:"name"`
	Authors                []Authors                `json:"authors"`
	Attachments            []Attachments            `json:"attachments"`
	WebsiteURL             string                   `json:"websiteUrl"`
	GameID                 int                      `json:"gameId"`
	Summary                string                   `json:"summary"`
	DefaultFileID          int                      `json:"defaultFileId"`
	DownloadCount          float64                  `json:"downloadCount"`
	LatestFiles            []LatestFile             `json:"latestFiles"`
	Categories             []Categories             `json:"categories"`
	Status                 int                      `json:"status"`
	PrimaryCategoryID      int                      `json:"primaryCategoryId"`
	CategorySection        CategorySection          `json:"categorySection"`
	Slug                   string                   `json:"slug"`
	GameVersionLatestFiles []GameVersionLatestFiles `json:"gameVersionLatestFiles"`
	IsFeatured             bool                     `json:"isFeatured"`
	PopularityScore        float64                  `json:"popularityScore"`
	GamePopularityRank     int                      `json:"gamePopularityRank"`
	PrimaryLanguage        string                   `json:"primaryLanguage"`
	GameSlug               string                   `json:"gameSlug"`
	GameName               string                   `json:"gameName"`
	PortalName             string                   `json:"portalName"`
	DateModified           time.Time                `json:"dateModified"`
	DateCreated            time.Time                `json:"dateCreated"`
	DateReleased           time.Time                `json:"dateReleased"`
	IsAvailable            bool                     `json:"isAvailable"`
	IsExperiemental        bool                     `json:"isExperiemental"`
}
type Authors struct {
	Name              string  `json:"name"`
	URL               string  `json:"url"`
	ProjectID         int     `json:"projectId"`
	ID                int     `json:"id"`
	ProjectTitleID    *int    `json:"projectTitleId"`
	ProjectTitleTitle *string `json:"projectTitleTitle"`
	UserID            int     `json:"userId"`
	TwitchID          int     `json:"twitchId"`
}
type Attachments struct {
	ID           int    `json:"id"`
	ProjectID    int    `json:"projectId"`
	Description  string `json:"description"`
	IsDefault    bool   `json:"isDefault"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Status       int    `json:"status"`
}
type Dependencies struct {
	ID      int `json:"id"`
	AddonID int `json:"addonId"`
	Type    int `json:"type"`
	FileID  int `json:"fileId"`
}
type Modules struct {
	Foldername  string `json:"foldername"`
	Fingerprint int    `json:"fingerprint"`
	Type        int    `json:"type"`
}
type SortableGameVersion struct {
	GameVersionPadded      string    `json:"gameVersionPadded"`
	GameVersion            string    `json:"gameVersion"`
	GameVersionReleaseDate time.Time `json:"gameVersionReleaseDate"`
	GameVersionName        string    `json:"gameVersionName"`
}
type LatestFile struct {
	ID                         int                   `json:"id"`
	DisplayName                string                `json:"displayName"`
	FileName                   string                `json:"fileName"`
	FileDate                   time.Time             `json:"fileDate"`
	FileLength                 int                   `json:"fileLength"`
	ReleaseType                int                   `json:"releaseType"`
	FileStatus                 int                   `json:"fileStatus"`
	DownloadURL                string                `json:"downloadUrl"`
	IsAlternate                bool                  `json:"isAlternate"`
	AlternateFileID            int                   `json:"alternateFileId"`
	Dependencies               []Dependencies        `json:"dependencies"`
	IsAvailable                bool                  `json:"isAvailable"`
	Modules                    []Modules             `json:"modules"`
	PackageFingerprint         int64                 `json:"packageFingerprint"`
	GameVersion                []string              `json:"gameVersion"`
	SortableGameVersion        []SortableGameVersion `json:"sortableGameVersion"`
	InstallMetadata            interface{}           `json:"installMetadata"`
	Changelog                  interface{}           `json:"changelog"`
	HasInstallScript           bool                  `json:"hasInstallScript"`
	IsCompatibleWithClient     bool                  `json:"isCompatibleWithClient"`
	CategorySectionPackageType int                   `json:"categorySectionPackageType"`
	RestrictProjectFileAccess  int                   `json:"restrictProjectFileAccess"`
	ProjectStatus              int                   `json:"projectStatus"`
	RenderCacheID              int                   `json:"renderCacheId"`
	FileLegacyMappingID        interface{}           `json:"fileLegacyMappingId"`
	ProjectID                  int                   `json:"projectId"`
	ParentProjectFileID        interface{}           `json:"parentProjectFileId"`
	ParentFileLegacyMappingID  interface{}           `json:"parentFileLegacyMappingId"`
	FileTypeID                 interface{}           `json:"fileTypeId"`
	ExposeAsAlternative        interface{}           `json:"exposeAsAlternative"`
	PackageFingerprintID       int                   `json:"packageFingerprintId"`
	GameVersionDateReleased    time.Time             `json:"gameVersionDateReleased"`
	GameVersionMappingID       int                   `json:"gameVersionMappingId"`
	GameVersionID              int                   `json:"gameVersionId"`
	GameID                     int                   `json:"gameId"`
	IsServerPack               bool                  `json:"isServerPack"`
	ServerPackFileID           *int                  `json:"serverPackFileId"`
	GameVersionFlavor          interface{}           `json:"gameVersionFlavor"`
}
type Categories struct {
	CategoryID int    `json:"categoryId"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	AvatarURL  string `json:"avatarUrl"`
	ParentID   int    `json:"parentId"`
	RootID     int    `json:"rootId"`
	ProjectID  int    `json:"projectId"`
	AvatarID   int    `json:"avatarId"`
	GameID     int    `json:"gameId"`
}
type CategorySection struct {
	ID                      int         `json:"id"`
	GameID                  int         `json:"gameId"`
	Name                    string      `json:"name"`
	PackageType             int         `json:"packageType"`
	Path                    string      `json:"path"`
	InitialInclusionPattern string      `json:"initialInclusionPattern"`
	ExtraIncludePattern     interface{} `json:"extraIncludePattern"`
	GameCategoryID          int         `json:"gameCategoryId"`
}
type GameVersionLatestFiles struct {
	GameVersion       string      `json:"gameVersion"`
	ProjectFileID     int         `json:"projectFileId"`
	ProjectFileName   string      `json:"projectFileName"`
	FileType          int         `json:"fileType"`
	GameVersionFlavor interface{} `json:"gameVersionFlavor"`
}
