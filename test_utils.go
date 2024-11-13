package gdal

import (
	"os"
	"path"
)

func makeTempFilePath(name string) string {
	tmpPath := path.Join(os.TempDir(), name)
	println(tmpPath)
	return tmpPath
}
