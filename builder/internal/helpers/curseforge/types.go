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

type GameApiStatus uint8
type GameStatus uint8

type CurseCategory struct {
	ID      uint32 `json:"id"`
	Slug    string `json:"slug"`
	IsClass bool   `json:"isClass"`
	ClassID uint32 `json:"classId"`
}

const (
	GameApiStatusPrivate GameApiStatus = iota + 1
	GameApiStatusPublic
)

const (
	GameStatusDraft GameStatus = iota + 1
	GameStatusTest
	GameStatusPendingReview
	GameStatusRejected
	GameStatusApproved
	GameStatusLive
)

const (
	ModloaderTypeAny ModloaderType = iota
	ModloaderTypeForge
	ModloaderTypeCauldron
	ModloaderTypeLiteloader
	ModloaderTypeFabric
	ModloaderTypeQuilt
	ModloaderTypeNeoForge
)

const (
	DependencyTypeEmbedded DependencyType = iota + 1
	DependencyTypeOptional
	DependencyTypeRequired
	DependencyTypeTool
	DependencyTypeIncompatible
	DependencyTypeInclude
)

type CurseGame struct {
	ID        uint32        `json:"id"`
	Name      string        `json:"name"`
	Slug      string        `json:"slug"`
	Status    GameStatus    `json:"status"`
	APIStatus GameApiStatus `json:"apiStatus"`
}
