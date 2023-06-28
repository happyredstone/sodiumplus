package curseforge

type HashAlgo uint8
type FileType uint8
type ModloaderType uint8
type DependencyType uint8

type CurseExportData struct {
	ProjectID uint32 `mapstructure:"project-id"`
}

type ModInfo struct {
	Name              string        `json:"name"`
	Summary           string        `json:"summary"`
	Slug              string        `json:"slug"`
	ID                uint32        `json:"id"`
	GameID            uint32        `json:"gameId"`
	PrimaryCategoryID uint32        `json:"primaryCategoryId"`
	ClassID           uint32        `json:"classId"`
	LatestFiles       []ModFileInfo `json:"latestFiles"`

	GameVersionLatestFiles []struct {
		GameVersion string        `json:"gameVersion"`
		ID          uint32        `json:"fileId"`
		Name        string        `json:"filename"`
		FileType    FileType      `json:"releaseType"`
		Modloader   ModloaderType `json:"modLoader"`
	} `json:"latestFilesIndexes"`

	ModLoaders []string `json:"modLoaders"`

	Links struct {
		WebsiteURL string `json:"websiteUrl"`
	} `json:"links"`
}
