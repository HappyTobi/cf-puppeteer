package v3

import "github.com/blang/semver"

//MinControllerVersion represents the v3 controller version
var MinControllerVersion = "3.27.0"

//GetMinSemVersion retun varsion as semver.Version
func GetMinSemVersion() (semver.Version, error) {
	return semver.Make(MinControllerVersion)
}
