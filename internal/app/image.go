package app

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func image(app *App) pulumi.String {
	return pulumi.String(fmt.Sprintf("docker.io/alexfalkowski/%s:v%s", app.Name, app.Version))
}
