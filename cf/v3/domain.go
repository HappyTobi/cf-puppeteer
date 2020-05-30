package v3

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Routes struct {
	Host   string
	Domain string
	Path   string
}

//DomainResponse reponse while loading domains
type DomainResponse struct {
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
		GUID          string    `json:"guid"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Name          string    `json:"name"`
		Internal      bool      `json:"internal"`
		Relationships struct {
			Organization struct {
				Data interface{} `json:"data"`
			} `json:"organization"`
			SharedOrganizations struct {
				Data []interface{} `json:"data"`
			} `json:"shared_organizations"`
		} `json:"relationships"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"resources"`
}

func (resource *ResourcesData) GetDomain(domains []map[string]string) (*[]Routes, error) {
	path := fmt.Sprintf(`/v3/domains`)

	response, err := resource.getDomain(path)
	if err != nil {
		return nil, err
	}

	domainGUID := make(map[string]Routes)

	for _, domainRes := range response.Resources {
		for _, routes := range domains {
			for _, domain := range routes {
				dhp := strings.Split(domain, "/")
				path := ""
				if len(dhp) > 1 {
					path = dhp[1]
				}
				domain = dhp[0]

				_, exists := domainGUID[domain]
				if strings.Contains(domain, domainRes.Name) && !exists {
					hostName := strings.ReplaceAll(domain, domainRes.Name, "")
					hostName = strings.TrimRight(hostName, ".")
					newRoute := &Routes{
						Host:   hostName,
						Domain: domainRes.Name,
						Path:   path,
					}
					domainGUID[domain] = *newRoute
				}
			}
		}
	}

	for response.Pagination.Next.Href != "" && len(domainGUID) <= 0 {
		response, err := resource.getDomain(response.Pagination.Next.Href)
		if err != nil {
			return nil, err
		}

		for _, domainRes := range response.Resources {
			for _, routes := range domains {
				for _, domain := range routes {
					dhp := strings.Split(domain, "/")
					path := ""
					if len(dhp) > 1 {
						path = dhp[1]
					}
					domain = dhp[0]

					_, exists := domainGUID[domain]
					if strings.Contains(domain, domainRes.Name) && !exists {
						hostName := strings.ReplaceAll(domain, domainRes.Name, "")
						hostName = strings.TrimRight(hostName, ".")
						newRoute := &Routes{
							Host:   hostName,
							Domain: domainRes.Name,
							Path:   path,
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

func (resource *ResourcesData) getDomain(path string) (*DomainResponse, error) {
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
		return response.Resources[i].Name < response.Resources[j].Name
	})
	return &response, nil
}
