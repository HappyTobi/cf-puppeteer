package v3action

import (
	"code.cloudfoundry.org/cli/actor/actionerror"
	"code.cloudfoundry.org/cli/actor/versioncheck"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccversion"
)

// Organization represents a V3 actor organization.
type Organization ccv3.Organization

// GetOrganizationByName returns the organization with the given name.
func (actor Actor) GetOrganizationByName(name string) (Organization, Warnings, error) {
	orgs, warnings, err := actor.CloudControllerClient.GetOrganizations(
		ccv3.Query{Key: ccv3.NameFilter, Values: []string{name}},
	)
	if err != nil {
		return Organization{}, Warnings(warnings), err
	}

	if len(orgs) == 0 {
		return Organization{}, Warnings(warnings), actionerror.OrganizationNotFoundError{Name: name}
	}

	return Organization(orgs[0]), Warnings(warnings), nil
}

func (actor Actor) GetOrganizationsByGUIDs(guids ...string) ([]Organization, Warnings, error) {
	currentV3Ver := actor.CloudControllerClient.CloudControllerAPIVersion()

	guidsSupport, err := versioncheck.IsMinimumAPIVersionMet(currentV3Ver, ccversion.MinVersionSpacesGUIDsParamV3)
	if err != nil {
		guidsSupport = false
	}

	queries := []ccv3.Query{}
	if guidsSupport {
		queries = []ccv3.Query{ccv3.Query{Key: ccv3.GUIDFilter, Values: guids}}
	}

	orgs, warnings, err := actor.CloudControllerClient.GetOrganizations(queries...)
	if err != nil {
		return []Organization{}, Warnings(warnings), err
	}

	guidToOrg := make(map[string]ccv3.Organization)
	for _, org := range orgs {
		guidToOrg[org.GUID] = org
	}

	filteredOrgs := make([]ccv3.Organization, 0)
	for _, guid := range guids {
		filteredOrgs = append(filteredOrgs, guidToOrg[guid])
	}
	orgs = filteredOrgs

	return actor.convertCCToActorOrganizations(orgs), Warnings(warnings), nil
}

func (actor Actor) convertCCToActorOrganizations(v3orgs []ccv3.Organization) []Organization {
	orgs := make([]Organization, len(v3orgs))
	for i := range v3orgs {
		orgs[i] = Organization{
			GUID: v3orgs[i].GUID,
			Name: v3orgs[i].Name,
		}
	}
	return orgs
}
