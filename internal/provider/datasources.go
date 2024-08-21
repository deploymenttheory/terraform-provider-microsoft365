package provider

import (
	"context"

	graphBetaAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/deviceandappmanagement/beta/assignmentFilter"
	graphBetaDeviceManagementScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/deviceandappmanagement/beta/deviceManagementScript"
	graphCloudPcDeviceImage "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/devicemanagement/v1.0/cloudPcDeviceImage"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSources returns a slice of functions that each return a datasource.DataSource.
// This function is a method of the M365Provider type and takes a context.Context as an argument.
// The returned slice is intended to hold the Microsoft 365 provider datasources.
//
// Parameters:
//   - ctx: The context for controlling cancellation and timeout.
//
// Returns:
//
//	[]func() datasource.DataSource: A slice of functions, each returning a datasource.DataSource.
func (p *M365Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		graphBetaAssignmentFilter.NewAssignmentFilterDataSource,
		graphBetaDeviceManagementScript.NewDeviceManagementScriptDataSource,
		graphCloudPcDeviceImage.NewCloudPcDeviceImageDataSource,

		// Add microsoft 365 provider datasources here
	}
}
