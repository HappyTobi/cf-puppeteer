package v3

import (
	"github.com/happytobi/cf-puppeteer/ui"
)

func (resource *ResourcesData) AssignAppManifest(manifestPath string) (err error) {
	args := []string{"v3-apply-manifest", "-f", manifestPath}

	err = resource.Executor.Execute(args)
	if err != nil {
		ui.Failed("could not read manifest from path %s error: %s", manifestPath, err)
		return err
	}
	return nil
}
