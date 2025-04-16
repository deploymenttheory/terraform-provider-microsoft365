package graphBetaLinuxPlatformScript

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	resource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/linux_platform_script"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for Linux Platform Script data sources.
//
// The function supports two methods of looking up a Linux Platform Script:
// 1. By ID - Uses a direct API call to fetch the specific script
// 2. By Name - Lists all configuration policies and finds the matching Linux Platform Script
//
// The function ensures that:
// - Either ID or Name is provided (but not both)
// - The lookup method is optimized based on the provided identifier
// - The remote state is properly mapped to the Terraform state
func (d *LinuxPlatformScriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object resource.LinuxPlatformScriptResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

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

	if !object.ID.IsNull() {
		// Lookup by ID
		tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", d.ProviderTypeName, d.TypeName, object.ID.ValueString()))

		baseResource, err := d.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		resource.MapRemoteResourceStateToTerraform(ctx, &object, baseResource)

	} else {
		// Lookup by Name
		tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with Name: %s", d.ProviderTypeName, d.TypeName, object.Name.ValueString()))

		configPolicies, err := d.client.
			DeviceManagement().
			ConfigurationPolicies().
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		var matchingPolicy graphmodels.DeviceManagementConfigurationPolicyable
		for _, policy := range configPolicies.GetValue() {
			// We need to check if this is a Linux platform script by checking platforms
			platforms := policy.GetPlatforms()
			if platforms != nil && *platforms == graphmodels.LINUX_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS {
				if policy.GetName() != nil && *policy.GetName() == object.Name.ValueString() {
					matchingPolicy = policy
					break
				}
			}
		}

		if matchingPolicy == nil {
			resp.Diagnostics.AddError(
				"Error Reading Linux Platform Script",
				fmt.Sprintf("No Linux platform script found with name: %s", object.Name.ValueString()),
			)
			return
		}

		// Now fetch the full details using the ID
		policyId := *matchingPolicy.GetId()
		object.ID = types.StringValue(policyId)

		baseResource, err := d.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(policyId).
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		resource.MapRemoteResourceStateToTerraform(ctx, &object, baseResource)

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Linux Platform Script Datasource Read Method: %s_%s", d.ProviderTypeName, d.TypeName))
}
