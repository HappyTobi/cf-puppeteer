package v3

import (
	"github.com/happytobi/cf-puppeteer/arguments"
	"github.com/happytobi/cf-puppeteer/ui"
)

func (resource *ResourcesData) AssignAppManifest(parsedArguments *arguments.ParserArguments) (err error) {
	args := []string{"v3-apply-manifest", "-f", parsedArguments.ManifestPath}

	ui.Say("apply manifest file")
	err = resource.Executor.Execute(args)
	if err != nil {
		ui.Failed("could not read manifest from path %s error: %s", parsedArguments.ManifestPath, err)
		return err
	}

	ui.Ok()

	return nil
}
