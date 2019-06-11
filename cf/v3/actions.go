package v3

import (
	"fmt"
	"github.com/happytobi/cf-puppeteer/ui"
	"os"
)

//AssignAppManifest assign an appManifest
func (resource *ResourcesData) AssignAppManifest(appLink string, manifestPath string) error {
	path := fmt.Sprintf(`%s/actions/apply_manifest`, appLink)

	file, err := os.Open(manifestPath)
	if err != nil {
		ui.Failed("could not read manifest from path %s error: %s", manifestPath, err)
		panic(err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Printf("manifest file stat error %s", err)
		return err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	file.Read(buffer)

	ui.Say("start uploading manifest file")
	_, err = resource.httpClient.PostJSON(path, buffer)
	if err != nil {
		return err
	}

	ui.Ok()

	return nil
}
