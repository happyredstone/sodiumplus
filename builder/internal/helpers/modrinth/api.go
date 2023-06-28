package modrinth

import (
	"net/http"

	modrinthApi "codeberg.org/jmansfield/go-modrinth/modrinth"
)

var ModrinthDefaultClient = modrinthApi.NewClient(&http.Client{})
