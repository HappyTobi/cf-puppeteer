package v7action

import (
	"code.cloudfoundry.org/cli/actor/actionerror"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccerror"
)

//go:generate counterfeiter . ManifestParser

type ManifestParser interface {
	AppNames() []string
	RawManifest(name string) ([]byte, error)
}

// ApplyApplicationManifest reads in the manifest from the path and provides it
// to the cloud controller.
func (actor Actor) ApplyApplicationManifest(parser ManifestParser, spaceGUID string) (Warnings, error) {
	var allWarnings Warnings

	for _, appName := range parser.AppNames() {
		rawManifest, err := parser.RawManifest(appName)
		if err != nil {
			return allWarnings, err
		}

		app, getAppWarnings, err := actor.GetApplicationByNameAndSpace(appName, spaceGUID)

		allWarnings = append(allWarnings, getAppWarnings...)
		if err != nil {
			return allWarnings, err
		}

		applyManifestWarnings, err := actor.SetApplicationManifest(app.GUID, rawManifest)
		allWarnings = append(allWarnings, applyManifestWarnings...)
		if err != nil {
			return allWarnings, err
		}
	}

	return allWarnings, nil
}

func (actor Actor) SetApplicationManifest(appGUID string, rawManifest []byte) (Warnings, error) {
	var allWarnings Warnings
	jobURL, applyManifestWarnings, err := actor.CloudControllerClient.UpdateApplicationApplyManifest(appGUID, rawManifest)
	allWarnings = append(allWarnings, applyManifestWarnings...)
	if err != nil {
		return allWarnings, err
	}

	pollWarnings, err := actor.CloudControllerClient.PollJob(jobURL)
	allWarnings = append(allWarnings, pollWarnings...)
	if err != nil {
		if newErr, ok := err.(ccerror.JobFailedError); ok {
			return allWarnings, actionerror.ApplicationManifestError{Message: newErr.Message}
		}
		return allWarnings, err
	}
	return allWarnings, nil
}

func (actor Actor) GetRawApplicationManifestByNameAndSpace(appName string, spaceGUID string) ([]byte, Warnings, error) {
	app, warnings, err := actor.GetApplicationByNameAndSpace(appName, spaceGUID)
	if err != nil {
		return nil, warnings, err
	}

	rawManifest, manifestWarnings, err := actor.CloudControllerClient.GetApplicationManifest(app.GUID)
	return rawManifest, append(warnings, manifestWarnings...), err
}
