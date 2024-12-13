package mcmeta

type ShortVersionInfo struct {
	Id              string `json:"id"`
	Kind            string `json:"type"`
	Url             string `json:"url"`
	Time            string `json:"time"`
	ReleaseTime     string `json:"releaseTime"`
	Sha1            string `json:"sha1"`
	ComplianceLevel int    `json:"complianceLevel"`
}

type VersionManifest struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`

	Versions []ShortVersionInfo `json:"versions"`
}
