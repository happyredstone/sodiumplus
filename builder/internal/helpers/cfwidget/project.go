package cfwidget

type ProjectUrls struct {
	CurseForge string `json:"curseforge"`
	Project    string `json:"project"`
}

type DownloadCount struct {
	Monthly int `json:"monthly"`
	Total   int `json:"total"`
}

type ProjectMember struct {
	Title    string `json:"title"`
	Username string `json:"username"`
	Id       int    `json:"id"`
}

type ProjectFile struct {
	Id         int      `json:"id"`
	Url        string   `json:"url"`
	Display    string   `json:"display"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Version    string   `json:"version"`
	Filesize   int      `json:"filesize"`
	Versions   []string `json:"versions"`
	Downloads  int      `json:"downloads"`
	UploadedAt string   `json:"uploaded_at"`
}

type Project struct {
	Id          int                      `json:"id"`
	Title       string                   `json:"title"`
	Summary     string                   `json:"summary"`
	Description string                   `json:"description"`
	Game        string                   `json:"game"`
	Type        string                   `json:"type"`
	Urls        ProjectUrls              `json:"urls"`
	Thumbnail   string                   `json:"thumbnail"`
	CreatedAt   string                   `json:"created_at"`
	Downloads   DownloadCount            `json:"downloads"`
	License     string                   `json:"license"`
	Donate      string                   `json:"donate"`
	Categories  []string                 `json:"categories"`
	Members     []ProjectMember          `json:"members"`
	Links       []string                 `json:"links"`
	Files       []ProjectFile            `json:"files"`
	Versions    map[string][]ProjectFile `json:"versions"`
	Download    ProjectFile              `json:"download"`
}
