package graphBetaWindowsDriverUpdateProfile

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for Windows Driver Update profile data sources.
//
// The function supports two methods of looking up a Windows Driver Update profile:
// 1. By ID - Uses a direct API call to fetch the specific profile
// 2. By Name - Lists all inventories and finds the matching one
//
// The function ensures that:
// - Either ID or Name is provided (but not both)
// - The Windows Driver Update Profile ID is always required
// - The lookup method is optimized based on the provided identifier
// - The remote state is properly mapped to the Terraform state
//
// The function will:
//  1. Extract and validate the configuration
//  2. Verify that exactly one identifier (ID or Name) is provided
//  3. Perform the appropriate API call based on the provided identifier
//  4. Map the remote state to the Terraform state
//  5. Handle any errors and return appropriate diagnostics
//
// If using ID:
//   - Makes a direct GET request to the specific resource endpoint
//   - Returns error if the ID is not found
//
// If using Name:
//   - Retrieves all inventories and searches for matching name
//   - Returns error if no matching profile is found
//
// Read handles the Read operation for Windows Driver Update Profile data sources.

// Read handles the Read operation for Windows Driver Update Profile data sources.
func (d *WindowsDriverUpdateProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object WindowsDriverUpdateProfileDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s", d.ProviderTypeName, d.TypeName))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Validate that exactly one of ID or DisplayName is provided
	if object.ID.IsNull() && object.DisplayName.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Either id or display_name must be provided",
		)
		return
	}
	if !object.ID.IsNull() && !object.DisplayName.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Only one of id or display_name should be provided, not both",
		)
		return
	}

	// Look up by ID if provided
	if !object.ID.IsNull() {
		profileID := object.ID.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Fetching Windows Driver Update Profile by ID: %s", profileID))

		profile, err := d.client.
			DeviceManagement().
			WindowsDriverUpdateProfiles().
			ByWindowsDriverUpdateProfileId(profileID).
			Get(ctx, nil)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Windows Driver Update Profile",
				fmt.Sprintf("Could not read Windows Driver Update Profile ID %s: %s", profileID, err),
			)
			return
		}

		// Map the response to the data source model using the dedicated function
		MapRemoteStateToDataSource(ctx, &object, profile)
	} else {
		// Look up by DisplayName
		displayName := object.DisplayName.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Fetching Windows Driver Update Profile by display name: %s", displayName))

		// Get all profiles and filter by display name
		profilesResult, err := d.client.
			DeviceManagement().
			WindowsDriverUpdateProfiles().
			Get(ctx, nil)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Windows Driver Update Profiles",
				fmt.Sprintf("Could not read Windows Driver Update Profiles: %s", err),
			)
			return
		}

		// Find profile by display name
		var foundProfile graphmodels.WindowsDriverUpdateProfileable
		for _, profile := range profilesResult.GetValue() {
			if *profile.GetDisplayName() == displayName {
				foundProfile = profile
				break
			}
		}

		if foundProfile == nil {
			resp.Diagnostics.AddError(
				"Error Reading Windows Driver Update Profile",
				fmt.Sprintf("No Windows Driver Update Profile found with display name: %s", displayName),
			)
			return
		}

		// Map the response to the data source model using the dedicated function
		MapRemoteStateToDataSource(ctx, &object, foundProfile)
	}

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s", d.ProviderTypeName, d.TypeName))
}
