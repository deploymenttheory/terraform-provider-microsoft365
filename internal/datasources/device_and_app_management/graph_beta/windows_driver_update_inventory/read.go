package graphBetaWindowsDriverUpdateInventory

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for Windows Driver Update Inventory data sources.
//
// The function supports two methods of looking up a Windows Driver Update Inventory:
// 1. By ID - Uses a direct API call to fetch the specific inventory
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
//   - Returns error if no matching inventory is found
func (d *WindowsDriverUpdateInventoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object WindowsDriverUpdateInventoryDataSourceModel

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

	// Validate that Windows Driver Update Profile ID is provided
	if object.WindowsDriverUpdateProfileID.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to search for a driver inventory.",
		)
		return
	}

	profileID := object.WindowsDriverUpdateProfileID.ValueString()

	// Validate that exactly one of ID or Name is provided
	if object.ID.IsNull() && object.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Either id or name must be provided",
		)
		return
	}
	if !object.ID.IsNull() && !object.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Only one of id or name should be provided, not both",
		)
		return
	}

	// Look up by ID if provided
	if !object.ID.IsNull() {
		inventoryID := object.ID.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Fetching Windows Driver Update Inventory by ID: %s", inventoryID))

		inventory, err := d.client.
			DeviceManagement().
			WindowsDriverUpdateProfiles().
			ByWindowsDriverUpdateProfileId(profileID).
			DriverInventories().
			ByWindowsDriverUpdateInventoryId(inventoryID).
			Get(ctx, nil)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Windows Driver Update Inventory",
				fmt.Sprintf("Could not read Windows Driver Update Inventory ID %s: %s", inventoryID, err),
			)
			return
		}

		// Map the response to the data source model using the dedicated function
		MapRemoteStateToDataSource(ctx, &object, inventory)
	} else {
		// Look up by Name
		name := object.Name.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Fetching Windows Driver Update Inventory by name: %s", name))

		// Get all inventories and filter by name
		inventoriesResult, err := d.client.
			DeviceManagement().
			WindowsDriverUpdateProfiles().
			ByWindowsDriverUpdateProfileId(profileID).
			DriverInventories().
			Get(ctx, nil)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Windows Driver Update Inventories",
				fmt.Sprintf("Could not read Windows Driver Update Inventories: %s", err),
			)
			return
		}

		// Find inventory by name
		var foundInventory graphmodels.WindowsDriverUpdateInventoryable
		for _, inventory := range inventoriesResult.GetValue() {
			if *inventory.GetName() == name {
				foundInventory = inventory
				break
			}
		}

		if foundInventory == nil {
			resp.Diagnostics.AddError(
				"Error Reading Windows Driver Update Inventory",
				fmt.Sprintf("No Windows Driver Update Inventory found with name: %s in profile: %s", name, profileID),
			)
			return
		}

		// Map the response to the data source model using the dedicated function
		MapRemoteStateToDataSource(ctx, &object, foundInventory)
	}

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s", d.ProviderTypeName, d.TypeName))
}
