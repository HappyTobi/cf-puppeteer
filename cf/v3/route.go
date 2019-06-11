package v3

import (
	"encoding/json"
	"fmt"
)

//RouteMapping
type RouteMapping struct {
	Relationships struct {
		App struct {
			GUID string `json:"guid"`
		} `json:"app"`
		Route struct {
			GUID string `json:"guid"`
		} `json:"route"`
	} `json:"relationships"`
}

//RouteMapping map route to application REMOVE?
func (resource *ResourcesData) RouteMapping(appGUID string, routeGUID string) error {
	path := fmt.Sprintf(`v3/route_mappings`)

	var v3RouteMapping RouteMapping
	v3RouteMapping.Relationships.App.GUID = appGUID
	v3RouteMapping.Relationships.Route.GUID = routeGUID

	appJSON, err := json.Marshal(v3RouteMapping)
	if err != nil {
		return err
	}

	_, err = resource.cli.PostJSON(path, string(appJSON))

	return err
}
