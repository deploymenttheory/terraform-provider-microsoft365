package provider

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/deviceandappmanagement/beta/assignmentFilter"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	tflog.Info(ctx, "Initializing Resources")

	if p.clients == nil {
		tflog.Warn(ctx, "Provider clients are not initialized.")
		return []func() resource.Resource{}
	}
	if p.clients.BetaClient == nil {
		tflog.Warn(ctx, "BetaClient is not initialized.")
		return []func() resource.Resource{}
	}

	tflog.Info(ctx, "Provider is initialized successfully.")

	return []func() resource.Resource{
		assignmentFilter.NewAssignmentFilterResource(p.clients.BetaClient),
		// Register other resources here
	}
}
