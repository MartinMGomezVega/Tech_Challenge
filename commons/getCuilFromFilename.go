package commons

import (
	"path/filepath"
	"strings"
)

// GetCuilFromFilename: Gets the filename quantile without the .csv file extension.
func GetCuilFromFilename(filename string) string {
	base := filepath.Base(filename)
	cuil := strings.TrimSuffix(base, filepath.Ext(base))
	return cuil
}
