package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for Windows Autopilot Device Preparation Policy data sources.
//
// The function supports two methods of looking up a policy:
// 1. By ID - Uses a direct API call to fetch the specific configuration policy
// 2. By Name - Lists configuration policies and finds the matching one
//
// Device preparation policies are settings catalog policies, so both lookups verify the
// template reference of the returned policy to ensure an unrelated settings catalog policy
// is never returned.
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object WindowsAutopilotDevicePreparationPolicyDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", DataSourceName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Validate that either ID or name is provided, but not both
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
		// Direct lookup by ID
		respResource, err := d.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
			Get(ctx, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		if !isDevicePreparationPolicy(respResource) {
			resp.Diagnostics.AddError(
				"Error Reading Windows Autopilot Device Preparation Policy",
				fmt.Sprintf("The configuration policy with ID %s is not a Windows Autopilot device preparation policy", object.ID.ValueString()),
			)
			return
		}

		MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	} else {
		// Lookup by name
		name := object.Name.ValueString()

		policies, err := d.listConfigurationPoliciesByName(ctx, name)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		var foundPolicy graphmodels.DeviceManagementConfigurationPolicyable
		for _, policy := range policies {
			if !isDevicePreparationPolicy(policy) {
				continue
			}
			if policy.GetName() != nil && *policy.GetName() == name {
				foundPolicy = policy
				break
			}
		}

		if foundPolicy == nil {
			resp.Diagnostics.AddError(
				"Error Reading Windows Autopilot Device Preparation Policy",
				fmt.Sprintf("No Windows Autopilot device preparation policy found with name: %s", name),
			)
			return
		}

		MapRemoteResourceStateToTerraform(ctx, &object, foundPolicy)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s", DataSourceName))
}

// listConfigurationPoliciesByName retrieves configuration policies, preferring a server side
// OData filter on name and falling back to an unfiltered list when the endpoint rejects the filter.
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) listConfigurationPoliciesByName(ctx context.Context, name string) ([]graphmodels.DeviceManagementConfigurationPolicyable, error) {
	filter := fmt.Sprintf("name eq '%s'", strings.ReplaceAll(name, "'", "''"))

	requestOptions := &devicemanagement.ConfigurationPoliciesRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.ConfigurationPoliciesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	respList, err := d.client.
		DeviceManagement().
		ConfigurationPolicies().
		Get(ctx, requestOptions)

	if err == nil {
		return respList.GetValue(), nil
	}

	tflog.Debug(ctx, "Server side name filter failed, falling back to listing all configuration policies", map[string]any{
		"error": err.Error(),
	})

	respList, err = d.client.
		DeviceManagement().
		ConfigurationPolicies().
		Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	return respList.GetValue(), nil
}

// isDevicePreparationPolicy checks that a configuration policy was created from one of the
// Autopilot device preparation templates.
func isDevicePreparationPolicy(policy graphmodels.DeviceManagementConfigurationPolicyable) bool {
	if policy == nil {
		return false
	}

	templateRef := policy.GetTemplateReference()
	if templateRef == nil {
		return false
	}

	family := templateRef.GetTemplateFamily()
	if family == nil || family.String() != templateFamily {
		return false
	}

	templateID := templateRef.GetTemplateId()
	if templateID == nil {
		return false
	}

	// Template ids are versioned, e.g. "<guid>_1", so match on the base template guid.
	base := strings.SplitN(*templateID, "_", 2)[0]

	return base == baseTemplateIDAutomatic || base == baseTemplateIDUserDriven
}
