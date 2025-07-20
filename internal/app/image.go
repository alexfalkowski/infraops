package app

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func image(app *App) pulumi.String {
	if app.IsInternal() {
		return pulumi.String(fmt.Sprintf("docker.io/alexfalkowski/%s:v%s", app.Name, app.Version))
	}
	return pulumi.String(fmt.Sprintf("docker.io/alexfalkowski/%s:%s", app.Name, app.Version))
}
