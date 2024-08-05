package main

import (
	"github.com/alexfalkowski/infraops/pulumi"
)

// Fns to create resources.
var Fns = pulumi.CreateFns{
	createSite, createAppConfig,
	createInfraOps, createDocker, createBin,
	createNonnative, createGoHealth, createGoService,
	createGoServiceTemplate, createGoClientTemplate,
	createStatus, createStandort, createAuth, createKonfig,
	createMigrieren, createBezeichner, createWeb,
	createServiceControl, createKonfigControl,
}
