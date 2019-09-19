package v3

import (
	"fmt"
	"github.com/pkg/errors"
)

func (resource *ResourcesData) AssignAppManifest(manifestPath string) (err error) {
	args := []string{"v3-apply-manifest", "-f", manifestPath}

	err = resource.Executor.Execute(args)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error while assigning manifest to application %s", manifestPath))
	}
	return nil
}
