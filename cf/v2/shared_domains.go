package v2

import (
	"encoding/json"
	"fmt"
	"time"
)

type SharedDomainResponse struct {
	Entity struct {
		Internal        bool        `json:"internal"`
		Name            string      `json:"name"`
		RouterGroupGUID interface{} `json:"router_group_guid"`
		RouterGroupType interface{} `json:"router_group_type"`
	} `json:"entity"`
	Metadata struct {
		CreatedAt time.Time `json:"created_at"`
		GUID      string    `json:"guid"`
		UpdatedAt time.Time `json:"updated_at"`
		URL       string    `json:"url"`
	} `json:"metadata"`
}

//LoadSharedDomains
func (resource *ResourcesData) LoadSharedDomains(domainGUID string) (*SharedDomainResponse, error) {
	path := fmt.Sprintf(`/v2/shared_domains/%s`, domainGUID)
	/*if resource.TraceLogging {
		fmt.Printf("get routes: %s:\n", path)
	}*/

	jsonResponse, err := resource.cli.GetJSON(path)
	if err != nil {
		return nil, err
	}

	var response SharedDomainResponse
	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
