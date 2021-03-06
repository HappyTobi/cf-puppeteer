package v2

import (
	"fmt"
)

//MapRoute map route to application
func (resource *LegacyResourcesData) MapRoute(appName string, host string, domain string, path string) (err error) {
	args := []string{"map-route", appName, domain, "--hostname", host}

	if len(path) > 0 {
		args = append(args, "--path", path)
	}

	fmt.Printf("map route %v", args)
	err = resource.Executor.Execute(args)
	if err != nil {
		return err
	}
	return nil
}

//UnMapRoute remove route from application
func (resource *LegacyResourcesData) UnMapRoute(appName string, host string, domain string, path string) (err error) {
	args := []string{"unmap-route", appName, domain, "--hostname", host}

	if len(path) > 0 {
		args = append(args, "--path", path)
	}

	fmt.Printf("map route %v", args)
	err = resource.Executor.Execute(args)
	if err != nil {
		return err
	}
	return nil
}
