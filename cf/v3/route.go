package v3

import (
	"fmt"
)

//MapRoute map route to application
func (resource *ResourcesData) MapRoute(appName string, host string, domain string) (err error) {
	args := []string{"map-route", appName, domain, "--hostname", host}
	fmt.Printf("map route %v", args)
	err = resource.Executor.Execute(args)
	if err != nil {
		return err
	}
	return nil
}

//UnMapRoute remove route from application
func (resource *ResourcesData) UnMapRoute(appName string, host string, domain string) (err error) {
	args := []string{"unmap-route", appName, domain, "--hostname", host}
	fmt.Printf("map route %v", args)
	err = resource.Executor.Execute(args)
	if err != nil {
		return err
	}
	return nil
}
