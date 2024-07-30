package lean

import (
	"github.com/alexfalkowski/infraops/area/apps/pulumi"
)

// Fns to create resources.
var Fns = pulumi.CreateFns{
	createKonfig,
	createStandort,
	createBezeichner,
	createWeb,
}
