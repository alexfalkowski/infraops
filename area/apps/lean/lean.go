package lean

import (
	"github.com/alexfalkowski/infraops/pulumi"
)

// Fns to create resources.
var Fns = pulumi.CreateFns{
	createKonfig,
	createStandort,
	createBezeichner,
	createWeb,
}
