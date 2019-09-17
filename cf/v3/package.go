package v3

import (
	"fmt"
	"os"
	"strings"
)

//GenerateTempFile generate path in temp directory bases on appName
func (resource *ResourcesData) GenerateTempFile(appName string, fileExtension string) (zipFile string) {
	tempDir := strings.TrimSuffix(os.TempDir(), "/")
	pathFormat := "%s/%s.%s"
	if strings.HasPrefix(appName, "/") {
		pathFormat = "%s%s.%s"
	}
	return fmt.Sprintf(pathFormat, tempDir, appName, fileExtension)
}
