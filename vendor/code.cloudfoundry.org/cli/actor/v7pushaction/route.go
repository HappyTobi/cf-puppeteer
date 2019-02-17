package v7pushaction

import (
	"regexp"
	"strings"

	"code.cloudfoundry.org/cli/actor/actionerror"
	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/actor/v7action"
	log "github.com/sirupsen/logrus"
)

func (actor Actor) CreateAndMapDefaultApplicationRoute(orgGUID string, spaceGUID string, app v7action.Application) (Warnings, error) {
	log.Info("default route creation/mapping")
	var warnings Warnings
	defaultRoute, domainWarnings, err := actor.getDefaultRoute(orgGUID, spaceGUID, app.Name)
	warnings = append(warnings, domainWarnings...)
	if err != nil {
		log.Errorln("getting default route:", err)
		return warnings, err
	}
	log.WithField("defaultRoute", defaultRoute.String()).Debug("calculated default route")

	boundRoutes, appRouteWarnings, err := actor.V2Actor.GetApplicationRoutes(app.GUID)
	warnings = append(warnings, appRouteWarnings...)
	if err != nil {
		log.Errorln("getting application routes:", err)
		return warnings, err
	}
	log.WithField("boundRoutes", boundRoutes.Summary()).Debug("existing app routes")

	_, routeAlreadyBound := actor.routeInListBySettings(defaultRoute, boundRoutes)
	if routeAlreadyBound {
		return warnings, err
	}

	spaceRoute, spaceRouteWarnings, err := actor.V2Actor.FindRouteBoundToSpaceWithSettings(defaultRoute)
	warnings = append(warnings, spaceRouteWarnings...)
	routeAlreadyExists := true
	if _, ok := err.(actionerror.RouteNotFoundError); ok {
		routeAlreadyExists = false
	} else if err != nil {
		log.Errorln("checking if route is in space:", err)
		return warnings, err
	} else {
		log.Debug("route already exists")
	}

	if !routeAlreadyExists {
		log.Debug("creating default route")
		var createRouteWarning v2action.Warnings
		spaceRoute, createRouteWarning, err = actor.V2Actor.CreateRoute(defaultRoute, false)
		warnings = append(warnings, createRouteWarning...)
		if err != nil {
			log.Errorln("creating route:", err)
			return warnings, err
		}
	}

	log.Debug("mapping default route")
	mapWarnings, err := actor.V2Actor.MapRouteToApplication(spaceRoute.GUID, app.GUID)
	warnings = append(warnings, mapWarnings...)
	return warnings, err
}

func (actor Actor) getDefaultRoute(orgGUID string, spaceGUID string, appName string) (v2action.Route, Warnings, error) {
	defaultDomain, domainWarnings, err := actor.DefaultDomain(orgGUID)
	if err != nil {
		return v2action.Route{}, domainWarnings, err
	}

	return v2action.Route{
		Host:      appName,
		Domain:    defaultDomain,
		SpaceGUID: spaceGUID,
	}, domainWarnings, nil
}

func (Actor) routeInListBySettings(route v2action.Route, routes []v2action.Route) (v2action.Route, bool) {
	for _, r := range routes {
		if r.Host == route.Host && r.Path == route.Path && r.Port == route.Port &&
			r.SpaceGUID == route.SpaceGUID && r.Domain.GUID == route.Domain.GUID {
			return r, true
		}
	}

	return v2action.Route{}, false
}

func (Actor) sanitize(name string) string {
	name = strings.ToLower(name)

	re := regexp.MustCompile("\\s+")
	name = re.ReplaceAllString(name, "-")

	re = regexp.MustCompile("[^[:alnum:]\\-]")
	name = re.ReplaceAllString(name, "")

	return strings.TrimLeft(name, "-")
}
