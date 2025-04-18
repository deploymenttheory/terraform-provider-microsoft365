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
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
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
		tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", d.ProviderTypeName, d.TypeName, object.ID.ValueString()))
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with Name: %s", d.ProviderTypeName, d.TypeName, object.Name.ValueString()))

		policyId, err := d.getResourceIdByName(ctx, object.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Linux Platform Script",
				err.Error(),
			)
			return
		}

		object.ID = types.StringValue(policyId)
	}

	if err := d.getDataSource(ctx, &object); err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Linux Platform Script Datasource Read Method: %s_%s", d.ProviderTypeName, d.TypeName))
}

// findPolicyIdByName looks up a Linux platform script by name and returns its ID
// findPolicyIdByName looks up a Linux platform script by name and returns its ID
func (d *LinuxPlatformScriptDataSource) getResourceIdByName(ctx context.Context, name string) (string, error) {
	tflog.Debug(ctx, fmt.Sprintf("Looking for Linux platform script with name: '%s'", name))

	// Build the filter string using the technologies field
	filterValue := fmt.Sprintf("technologies/any(t:t eq 'linuxMdm') and name eq '%s'", name)

	// Set up expand to get more complete data
	expand := []string{"settings"}

	// Create request options with proper filtering and expand
	requestOptions := &devicemanagement.ConfigurationPoliciesRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.ConfigurationPoliciesRequestBuilderGetQueryParameters{
			Filter: &filterValue,
			Expand: expand,
		},
	}

	configPolicies, err := d.client.
		DeviceManagement().
		ConfigurationPolicies().
		Get(ctx, requestOptions)

	if err != nil {
		return "", err
	}

	// Check if we found anything with the combined filter
	if configPolicies.GetValue() != nil && len(configPolicies.GetValue()) > 0 {
		policy := configPolicies.GetValue()[0]
		tflog.Debug(ctx, fmt.Sprintf("Found Linux script with name: '%s' and ID: %s", name, *policy.GetId()))
		return *policy.GetId(), nil
	}

	// If we didn't find with combined filter, try just filtering by technology
	technologyFilter := "technologies/any(t:t eq 'linuxMdm')"
	requestOptions = &devicemanagement.ConfigurationPoliciesRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.ConfigurationPoliciesRequestBuilderGetQueryParameters{
			Filter: &technologyFilter,
			Expand: expand,
		},
	}

	configPolicies, err = d.client.
		DeviceManagement().
		ConfigurationPolicies().
		Get(ctx, requestOptions)

	if err != nil {
		return "", err
	}

	// Now search by name in the returned Linux policies
	if configPolicies.GetValue() != nil {
		for _, policy := range configPolicies.GetValue() {
			if policy.GetName() != nil && *policy.GetName() == name {
				tflog.Debug(ctx, fmt.Sprintf("Found Linux script with name: '%s' and ID: %s", name, *policy.GetId()))
				return *policy.GetId(), nil
			}
		}

		// Log available Linux scripts for debugging
		tflog.Debug(ctx, fmt.Sprintf("Found %d Linux scripts:", len(configPolicies.GetValue())))

		linuxScriptNames := []string{}
		for _, policy := range configPolicies.GetValue() {
			if policy.GetName() != nil {
				linuxScriptNames = append(linuxScriptNames, *policy.GetName())
				tflog.Debug(ctx, fmt.Sprintf("Linux script: Name='%s', ID='%s'", *policy.GetName(), *policy.GetId()))
			}
		}

		if len(linuxScriptNames) > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Available Linux script names: %v", linuxScriptNames))
		}
	} else {
		tflog.Debug(ctx, "No Linux scripts found with technology 'linuxMdm'")
	}

	return "", fmt.Errorf("no Linux platform script found with name: %s", name)
}

// getDataSource fetches all details for a Linux script and maps them to the Terraform model
func (d *LinuxPlatformScriptDataSource) getDataSource(ctx context.Context, object *resource.LinuxPlatformScriptResourceModel) error {
	// 1. Fetch base resource
	baseResource, err := d.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		return err
	}

	// 2. Map base resource
	resource.MapRemoteResourceStateToTerraform(ctx, object, baseResource)

	// 3. Fetch settings
	settingsResponse, err := d.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Settings().
		Get(ctx, nil)

	if err != nil {
		return err
	}

	// 4. Map settings
	resource.MapRemoteSettingsStateToTerraform(ctx, object, settingsResponse.GetValue())

	// 5. Fetch assignments
	assignmentsResponse, err := d.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		return err
	}

	// 6. Map assignments (using a type assertion since assignments API might be different)
	scriptAssignments, ok := assignmentsResponse.(models.DeviceManagementScriptAssignmentCollectionResponseable)
	if ok {
		resource.MapRemoteAssignmentStateToTerraform(ctx, object, scriptAssignments)
	} else {
		tflog.Warn(ctx, "Couldn't cast assignments to DeviceManagementScriptAssignmentCollectionResponseable")
	}

	return nil
}
