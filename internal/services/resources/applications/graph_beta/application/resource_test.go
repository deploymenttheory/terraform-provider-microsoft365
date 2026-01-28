package graphBetaApplication_test

import (
	graphBetaApplication "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/application"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaApplication.ResourceName

	// testResource is the test resource implementation for applications
	testResource = graphBetaApplication.ApplicationTestResource{}
)
