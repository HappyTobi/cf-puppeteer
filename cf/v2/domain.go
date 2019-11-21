package v2

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Routes struct {
	Host   string
	Domain string
}

//DomainResponse reponse while loading domains
type DomainResponse struct {
	Pagination struct {
		TotalResults int    `json:"total_results"`
		TotalPages   int    `json:"total_pages"`
		NextUrl      string `json:"next_url"`
		PrevUrl      string `json:"prev_url"`
	} `json:"pagination"`
	Resources []struct {
		Metadata struct {
			GUID string `json:"guid"`
			Url  string `json:"url"`
		} `json:"metadata"`
		Entity struct {
			Name     string `json:"name"`
			Internal bool   `json:"internal"`
		} `json:"entity"`
	} `json:"resources"`
}

func (resource *LegacyResourcesData) GetDomain(domains []map[string]string) (*[]Routes, error) {
	//default order asc.
	path := fmt.Sprintf(`/v2/domains`)

	response, err := resource.getDomain(path)
	if err != nil {
		return nil, err
	}

	domainGUID := make(map[string]Routes)
	for _, domainRes := range response.Resources {
		for _, routes := range domains {
			domain := routes["route"]
			hostName := strings.ReplaceAll(domain, domainRes.Entity.Name, "")

			_, exists := domainGUID[domain]
			if exists {
				exists = len(domainGUID[domain].Host) < len(hostName)
			}

			//question ist when route matches 2 time what kind of your we are using?
			if strings.Contains(domain, domainRes.Entity.Name) && len(hostName) > 0 && !exists {
				hostName = strings.TrimRight(hostName, ".")
				newRoute := &Routes{
					Host:   hostName,
					Domain: domainRes.Entity.Name,
				}
				domainGUID[domain] = *newRoute
			}
		}
	}

	//move to func and all recursive
	for response.Pagination.NextUrl != "" && len(domainGUID) <= 0 {
		response, err := resource.getDomain(response.Pagination.NextUrl)
		if err != nil {
			return nil, err
		}

		for _, domainRes := range response.Resources {
			for _, routes := range domains {
				for _, domain := range routes {
					_, exists := domainGUID[domain]
					if strings.Contains(domain, domainRes.Entity.Name) && !exists {
						hostName := strings.ReplaceAll(domain, domainRes.Entity.Name, "")
						hostName = strings.TrimRight(hostName, ".")
						newRoute := &Routes{
							Host:   hostName,
							Domain: domainRes.Entity.Name,
						}
						domainGUID[domain] = *newRoute
					}
				}
			}
		}
	}

	var domainsFound []Routes
	for _, v := range domainGUID {
		domainsFound = append(domainsFound, v)
	}

	return &domainsFound, err
}

func (resource *LegacyResourcesData) getDomain(path string) (*DomainResponse, error) {
	var response DomainResponse

	jsonResult, err := resource.Cli.GetJSON(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonResult), &response)
	if err != nil {
		return nil, err
	}

	sort.Slice(response.Resources, func(i, j int) bool {
		return response.Resources[i].Entity.Name < response.Resources[j].Entity.Name
	})
	return &response, nil
}
