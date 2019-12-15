package manifest

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
)

/*
Application yaml represents the complete yaml structure
type is always string because you can use a vars placeholder in all attributes.
*/
type Application struct {
	Name                    string              `yaml:"name"`
	Instances               string              `yaml:"instances,omitempty"`
	Memory                  string              `yaml:"memory,omitempty"`
	DiskQuota               string              `yaml:"disk_quota,omitempty"`
	Routes                  []map[string]string `yaml:"routes,omitempty"`
	NoRoute                 string              `yaml:"no-route,omitempty"`
	Buildpacks              []string            `yaml:"buildpacks,omitempty"`
	Command                 string              `yaml:"command,omitempty"`
	Env                     map[string]string   `yaml:"env,omitempty"`
	Services                []string            `yaml:"services,omitempty"`
	Stack                   string              `yaml:"stack,omitempty"`
	Path                    string              `yaml:"path,omitempty"`
	Timeout                 string              `yaml:"timeout,omitempty"`
	HealthCheckType         string              `yaml:"health-check-type,omitempty"`
	HealthCheckHTTPEndpoint string              `yaml:"health-check-http-endpoint,omitempty"`
}

// Manifest struct represents the application manifest.
type Manifest struct {
	ApplicationManifests []Application `yaml:"applications"`
}

//VarsFile
type Variables map[string]interface{}

//regex pattern - @see cf cli code
var (
	variablePlaceholderPattern = `\(\((!?[-/\.\w\pL]+)\)\)`
	interpolationRegex         = regexp.MustCompile(`\(\((!?[-/\.\w\pL]+)\)\)`)
)

//ParseAndReplaceWithVars parse a manifest and vars file.
// get all values from vars file and put them into the manifest file so there will be a returned new manifest without
// placeholders
func ParseApplicationManifest(manifestFilePath string, varsFilePath string) (manifest Manifest, err error) {
	document, err := loadYmlFile(manifestFilePath)

	if err != nil || document.ApplicationManifests == nil {
		return Manifest{}, fmt.Errorf("could not parse file, file not valid")
	}

	//if there's no vars file, we can return the parsed manifest direct
	if len(varsFilePath) <= 0 {
		return document, nil
	}

	varsFile, err := loadVarsFile(varsFilePath)
	if err != nil {
		return Manifest{}, fmt.Errorf("could not parse vars file, file not valid")
	}

	//iterate through all the applications an check if vars are existing
	for aI, app := range document.ApplicationManifests {
		appManifestElement := reflect.ValueOf(&app).Elem()
		for i := 0; i < appManifestElement.NumField(); i++ {
			appElementField := appManifestElement.Field(i)
			fieldValue := fmt.Sprintf("%v", appElementField.Interface())
			for mIndex, match := range interpolationRegex.FindAllSubmatch([]byte(fieldValue), -1) {
				//get variable name
				matchedVar := strings.TrimPrefix(string(match[1]), "!")
				varsValue := varsFile[matchedVar]
				//change string fields
				if appElementField.Kind() == reflect.String {
					replacedVar := strings.ReplaceAll(fieldValue, fmt.Sprintf("((%v))", matchedVar), fmt.Sprintf("%v", varsValue))
					appElementField.SetString(replacedVar)
				}
				//change slices
				if appElementField.Kind() == reflect.Slice {
					sliceElementField := appElementField.Index(mIndex)
					if sliceElementField.Kind() == reflect.Map {
						for _, mv := range sliceElementField.MapKeys() {
							x := sliceElementField.MapIndex(mv)
							nx := strings.ReplaceAll(x.String(), fmt.Sprintf("((%v))", matchedVar), fmt.Sprintf("%v", varsValue))
							sliceElementField.SetMapIndex(mv, reflect.ValueOf(nx))
						}
					}
				}
			}
		}
		document.ApplicationManifests[aI] = app
	}

	return document, nil
}

//load the vars file an throw errors then there is a issue
func loadVarsFile(varsFilePath string) (variables Variables, err error) {
	fileBytes, err := ioutil.ReadFile(varsFilePath)
	if err != nil {
		return variables, fmt.Errorf("error while reading varsfile: %s", varsFilePath)
	}

	err = yaml.Unmarshal(fileBytes, &variables)
	if err != nil {
		return variables, fmt.Errorf("error while parsing the varsfile %s error: %v", varsFilePath, err)
	}

	return variables, nil
}

//load the application yml file an throw errors then there is a issue
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
	return ParseApplicationManifest(manifestFilePath, "")
}
