package provider

import (
	"context"

	graphBetaAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/deviceandappmanagement/beta/assignmentFilter"
	graphCloudPcProvisioningPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/devicemanagement/v1.0/cloudPcProvisioningPolicy"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Resources returns a slice of functions that each return a resource.Resource.
// This function is a method of the M365Provider type and takes a context.Context as an argument.
// The returned slice is intended to hold the Microsoft 365 provider resources.
//
// Parameters:
//   - ctx: The context for controlling cancellation and timeout.
//
// Returns:
//
//	[]func() resource.Resource: A slice of functions, each returning a resource.Resource.
//
// Resources returns a slice of functions that each return a resource.Resource.
func (p *M365Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		graphBetaAssignmentFilter.NewAssignmentFilterResource,
		graphCloudPcProvisioningPolicy.NewCloudPcProvisioningPolicyResource,
	}
}
