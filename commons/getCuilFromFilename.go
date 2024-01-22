package commons

import (
	"path/filepath"
	"strings"
)

// Función para obtener el cuil del nombre del archivo sin la extensión .csv
func GetCuilFromFilename(filename string) string {
	base := filepath.Base(filename)
	cuil := strings.TrimSuffix(base, filepath.Ext(base))
	return cuil
}
