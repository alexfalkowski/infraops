package main

import (
	"github.com/alexfalkowski/infraops/internal/pulumi"
)

// Fns to create resources.
var Fns = pulumi.CreateFns{
	createSite, createAppConfig,
	createInfraOps, createDocker, createBin,
	createNonnative, createGoHealth, createGoService,
	createGoServiceTemplate, createGoClientTemplate,
	createStatus, createStandort, createKonfig,
	createMigrieren, createBezeichner, createWeb,
	createKonfigCtl,
}
