package curseforge

import (
	"archive/zip"
	"bufio"
	"strconv"

	"github.com/packwiz/packwiz/core"
)

func CreateCurseModlist(zw *zip.Writer, mods []*core.Mod) error {
	modlistFile, err := zw.Create("modlist.html")

	if err != nil {
		return err
	}

	w := bufio.NewWriter(modlistFile)

	_, err = w.WriteString("<ul>\r\n")

	if err != nil {
		return err
	}

	for _, mod := range mods {
		projectRaw, ok := mod.GetParsedUpdateData("curseforge")

		if !ok {
			_, err = w.WriteString("<li>" + mod.Name + "</li>\r\n")

			if err != nil {
				return err
			}

			continue
		}

		project := projectRaw.(CurseUpdateData)

		_, err = w.WriteString("<li><a href=\"https://www.curseforge.com/projects/" + strconv.FormatUint(uint64(project.ProjectID), 10) + "\">" + mod.Name + "</a></li>\r\n")

		if err != nil {
			return err
		}
	}

	_, err = w.WriteString("</ul>\r\n")

	if err != nil {
		return err
	}

	return w.Flush()
}
