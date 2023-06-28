package curseforge

import "github.com/mitchellh/mapstructure"

func ParseCurseExportData(from map[string]interface{}) (CurseExportData, error) {
	var exportData CurseExportData

	err := mapstructure.Decode(from, &exportData)

	return exportData, err
}
