package provider

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/deviceandappmanagement/beta/assignmentFilter"
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
func (p *M365Provider) Resources(ctx context.Context) []func() resource.Resource {
	clients, ok := p.getClients(ctx)
	if !ok {
		return nil
	}

	return []func() resource.Resource{
		assignmentFilter.NewAssignmentFilterResource(clients.BetaClient),
		// Add other Microsoft 365 provider resources here, using the appropriate client.
	}
}
