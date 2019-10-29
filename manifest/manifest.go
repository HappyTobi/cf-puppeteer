package manifest

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

/*
Application yaml represents the complete yaml structure
*/
type Application struct {
	Name                    string              `yaml:"name"`
	Instances               int                 `yaml:"instances,omitempty"`
	Memory                  string              `yaml:"memory,omitempty"`
	DiskQuota               string              `yaml:"disk_quota,omitempty"`
	Routes                  []map[string]string `yaml:"routes,omitempty"`
	NoRoute                 bool                `yaml:"no-route,omitempty"`
	Buildpacks              []string            `yaml:"buildpacks,omitempty"`
	Command                 string              `yaml:"command,omitempty"`
	Env                     map[string]string   `yaml:"env,omitempty"`
	Services                []string            `yaml:"services,omitempty"`
	Stack                   string              `yaml:"stack,omitempty"`
	Timeout                 int                 `yaml:"timeout,omitempty"`
	HealthCheckType         string              `yaml:"health-check-type,omitempty"`
	HealthCheckHTTPEndpoint string              `yaml:"health-check-http-endpoint,omitempty"`
	AppPath					string				`yaml:"path,omitempty"`
}

// Manifest struct represents the application manifest.
type Manifest struct {
	ApplicationManifests []Application `yaml:"applications"`
}

// Parse parse application manifest files from provided path and return
// (right now) the app name of the first found application.
func Parse(manifestFilePath string) (manifest Manifest, err error) {
	document, err := loadYmlFile(manifestFilePath)

	if err != nil || document.ApplicationManifests == nil {
		return document, fmt.Errorf("could not parse file - file not valid")
	}

	return document, nil
}

func loadYmlFile(manifestFilePath string) (manifest Manifest, err error) {
	fileBytes, err := ioutil.ReadFile(manifestFilePath)
	if err != nil {
		return Manifest{}, fmt.Errorf("error while reading manifest: %s", manifestFilePath)
	}

	var document Manifest
	err = yaml.Unmarshal(fileBytes, &document)
	if err != nil {
		return Manifest{}, fmt.Errorf("error while parsing the manifest %s error: %v", manifestFilePath, err)
	}

	return document, nil
}

//WriteYmlFile write yml file to specified path and return them parsed
func WriteYmlFile(manifestFilePath string, manifest Manifest) (newManifest Manifest, err error) {
	mManifest, err := yaml.Marshal(&manifest)
	if err != nil {
		return Manifest{}, err
	}
	bManifest := []byte(string(mManifest))
	err = ioutil.WriteFile(manifestFilePath, bManifest, 0644)
	if err != nil {
		return Manifest{}, err
	}
	return Parse(manifestFilePath)

}
