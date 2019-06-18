package v3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

//V3Package represents post model of V3Package body
type Package struct {
	PackageType   string `json:"type"`
	Relationships struct {
		App struct {
			Data struct {
				GUID string `json:"guid"`
			} `json:"data"`
		} `json:"app"`
	} `json:"relationships"`
}

//V3PackageResponse create package response payload
type PackageResponse struct {
	GUID  string `json:"guid"`
	State string `json:"state"`
	Links struct {
		Upload struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"upload"`
	} `json:"links"`
}

//CreatePackage create a package with v3 cf api
func (resource *ResourcesData) CreatePackage(appGUID string) (*PackageResponse, error) {
	path := fmt.Sprintf(`/v3/packages`)
	var v3Package Package
	v3Package.PackageType = "bits"
	v3Package.Relationships.App.Data.GUID = appGUID

	//TODO move to function
	appJSON, err := json.Marshal(v3Package)
	if err != nil {
		return nil, err
	}

	jsonResult, err := resource.Cli.PostJSON(path, string(appJSON))
	if err != nil {
		return nil, err
	}

	var response PackageResponse
	err = json.Unmarshal([]byte(jsonResult), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//CheckPackageState create a package with v3 cf api
func (resource *ResourcesData) CheckPackageState(packageGUID string) (*PackageResponse, error) {
	path := fmt.Sprintf(`/v3/packages/%s`, packageGUID)

	jsonResult, err := resource.Cli.GetJSON(path)

	if err != nil {
		return nil, err
	}

	var response PackageResponse
	err = json.Unmarshal([]byte(jsonResult), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//UploadApplication upload a zip file to the created package
func (resource *ResourcesData) UploadApplication(appName string, applicationFiles string, targetURL string) (*PackageResponse, error) {
	if !resource.zipper.IsZipFile(applicationFiles) {
		zipFileName := fmt.Sprintf("%s%s.zip", os.TempDir(), appName)
		newZipFile, err := os.Create(zipFileName)

		if err != nil {
			return nil, err
		}
		defer newZipFile.Close()

		err = resource.zipper.Zip(applicationFiles, newZipFile)
		applicationFiles = zipFileName
	}

	file, err := os.Open(applicationFiles)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("bits", filepath.Base(applicationFiles))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	jsonResult, err := resource.httpClient.PostFormData(targetURL, body.Bytes(), writer.FormDataContentType())

	var response PackageResponse
	err = json.Unmarshal([]byte(jsonResult), &response)
	if err != nil {
		return nil, err
	}

	return &response, err
}
