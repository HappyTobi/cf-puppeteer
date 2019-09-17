package v3

import (
	"encoding/json"
	"fmt"
	"github.com/happytobi/cf-puppeteer/arguments"
	"github.com/happytobi/cf-puppeteer/ui"
	"strings"
	"time"
)

//RouteMappingResponse
type RouteMappingResponse struct {
	Pagination struct {
		TotalResults int `json:"total_results"`
		TotalPages   int `json:"total_pages"`
		First        struct {
			Href string `json:"href"`
		} `json:"first"`
		Last struct {
			Href string `json:"href"`
		} `json:"last"`
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
		Previous interface{} `json:"previous"`
	} `json:"pagination"`
	Resources []struct {
		GUID      string    `json:"guid"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Weight    int       `json:"weight"`
		Links     struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			App struct {
				Href string `json:"href"`
			} `json:"app"`
			Route struct {
				Href string `json:"href"`
			} `json:"route"`
			Process struct {
				Href string `json:"href"`
			} `json:"process"`
		} `json:"links"`
	} `json:"resources"`
}

/* TODO add type to docker */
func (resource *ResourcesData) CreateApp(parsedArguments *arguments.ParserArguments) (err error) {
	args := []string{"v3-create-app", parsedArguments.AppName}
	ui.Say("create application %s", args)
	err = resource.Executor.Execute(args)
	if err != nil {
		return err
	}
	return nil
}

func (resource *ResourcesData) PushApp(parsedArguments *arguments.ParserArguments) (err error) {
	args := []string{"v3-push", parsedArguments.AppName, "--no-start"}
	if parsedArguments.AppPath != "" {
		args = append(args, "-p", parsedArguments.AppPath)
	}

	if parsedArguments.NoRoute == true {
		args = append(args, "--no-route")
	}

	ui.Say("start pushing application with arguments %s", args)
	err = resource.Executor.Execute(args)
	if err != nil {
		return err
	}

	ui.Say("apply environment variables")
	for _, env := range parsedArguments.MergedEnvs {
		args = []string{"v3-set-env", parsedArguments.AppName, env}
		err = resource.Executor.Execute(args)
		if err != nil {
			ui.Failed("could not set environment variable %s to application %s", env, parsedArguments.AppName)
		}
	}

	return nil
}

//AssignApp to a created droplet guid
func (resource *ResourcesData) GetRoutesApp(appGUID string) ([]string, error) {
	path := fmt.Sprintf(`/v3/apps/%s/route_mappings`, appGUID)

	jsonResult, err := resource.Cli.GetJSON(path)

	var response RouteMappingResponse

	err = json.Unmarshal([]byte(jsonResult), &response)
	var routeGUIDs []string
	for _, mapResource := range response.Resources {
		count := strings.LastIndex(mapResource.Links.Route.Href, "/")
		if count > 0 {
			routeGUIDs = append(routeGUIDs, mapResource.Links.Route.Href[count+1:])
		}
	}

	return routeGUIDs, err
}
