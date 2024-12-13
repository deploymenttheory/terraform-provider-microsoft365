package graphBetaRoleScopeTag

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	resource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/role_scope_tag"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for the RoleScopeTagDataSource.
//
// The function:
// - Gets the role scope tag ID or display name from the configuration
// - Retrieves the role scope tag details from Intune
// - Maps the properties to the data source schema
// - Updates the Terraform state with the current configuration
func (d *RoleScopeTagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object resource.RoleScopeTagResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Datasource Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, resource.ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var respResource graphmodels.RoleScopeTagable
	var err error

	// If ID is provided, get by ID
	if !object.ID.IsNull() && !object.ID.IsUnknown() {
		tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", d.ProviderTypeName, d.TypeName, object.ID.ValueString()))

		respResource, err = d.client.
			DeviceManagement().
			RoleScopeTags().
			ByRoleScopeTagId(object.ID.ValueString()).
			Get(ctx, nil)
	} else if !object.DisplayName.IsNull() && !object.DisplayName.IsUnknown() {
		// If display name is provided, filter by display name
		tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with DisplayName: %s", d.ProviderTypeName, d.TypeName, object.DisplayName.ValueString()))

		filter := fmt.Sprintf("displayName eq '%s'", object.DisplayName.ValueString())
		respCollection, err := d.client.
			DeviceManagement().
			RoleScopeTags().
			Get(ctx, &msgraphsdk.RoleScopeTagsRequestBuilderGetRequestConfiguration{
				QueryParameters: &msgraphsdk.RoleScopeTagsRequestBuilderGetQueryParameters{
					Filter: &filter,
				},
			})

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		if roles := respCollection.GetValue(); len(roles) > 0 {
			respResource = roles[0]
		} else {
			resp.Diagnostics.AddError(
				"Role Scope Tag Not Found",
				fmt.Sprintf("No role scope tag found with display name: %s", object.DisplayName.ValueString()),
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Missing Required Field",
			"Either id or display_name must be provided",
		)
		return
	}

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return
	}

	resource.MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s", d.ProviderTypeName, d.TypeName))
}
