// REF: https://learn.microsoft.com/en-us/graph/api/intune-deviceconfigv2-devicemanagementconfigurationpolicy-list?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/intune-deviceconfigv2-devicemanagementconfigurationpolicy-get?view=graph-rest-beta
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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

type lookupMethod int

const (
	lookupByPolicyId lookupMethod = iota
	lookupByName
	lookupByODataQuery
	lookupListAll
)

// Read handles the Read operation.
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object WindowsAutopilotDevicePreparationPolicyDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	method := d.determineLookupMethod(object)

	var policies []graphmodels.DeviceManagementConfigurationPolicyable
	var err error

	switch method {
	case lookupByPolicyId:
		policies, err = d.getPolicyByPolicyId(ctx, object)
	case lookupByName:
		policies, err = d.getPoliciesByName(ctx, object)
	case lookupByODataQuery:
		policies, err = d.getPoliciesByODataQuery(ctx, object)
	case lookupListAll:
		policies, err = d.listAllPolicies(ctx)
	}

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	// A settings catalog policy is only a device preparation policy when it was created from one of
	// the Autopilot device preparation templates, so unrelated policies are discarded here.
	policies = filterDevicePreparationPolicies(policies)

	if method == lookupByPolicyId && len(policies) == 0 {
		resp.Diagnostics.AddError(
			"Error Reading Windows Autopilot Device Preparation Policy",
			fmt.Sprintf("The configuration policy with ID %s is not a Windows Autopilot device preparation policy", object.PolicyId.ValueString()),
		)
		return
	}

	if method == lookupByName && len(policies) == 0 {
		resp.Diagnostics.AddError(
			"Error Reading Windows Autopilot Device Preparation Policy",
			fmt.Sprintf("No Windows Autopilot device preparation policy found with name: %s", object.Name.ValueString()),
		)
		return
	}

	object.Items = ConstructPolicyItems(policies)

	// Assignments are a navigation property, so they need a separate request against a single policy.
	if !object.ListAssignments.IsNull() && object.ListAssignments.ValueBool() && len(policies) > 0 {
		policyId := ""
		if policies[0].GetId() != nil {
			policyId = *policies[0].GetId()
		}

		if policyId != "" {
			assignments, err := d.getPolicyAssignments(ctx, policyId)
			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
				return
			}

			object.Assignments = ConstructPolicyAssignments(assignments)
		}
	}

	object.ID = types.StringValue(fmt.Sprintf("windows-autopilot-device-preparation-policy-datasource-%d", time.Now().Unix()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(object.Items)))
}

// determineLookupMethod determines which lookup method to use based on the provided attributes
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) determineLookupMethod(object WindowsAutopilotDevicePreparationPolicyDataSourceModel) lookupMethod {
	switch {
	case !object.PolicyId.IsNull() && object.PolicyId.ValueString() != "":
		return lookupByPolicyId
	case !object.Name.IsNull() && object.Name.ValueString() != "":
		return lookupByName
	case !object.ODataQuery.IsNull() && object.ODataQuery.ValueString() != "":
		return lookupByODataQuery
	default:
		return lookupListAll
	}
}

// getPolicyByPolicyId retrieves a single policy by its id
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) getPolicyByPolicyId(ctx context.Context, object WindowsAutopilotDevicePreparationPolicyDataSourceModel) ([]graphmodels.DeviceManagementConfigurationPolicyable, error) {
	policyId := object.PolicyId.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up device preparation policy by policy ID: %s", policyId))

	policy, err := d.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(policyId).
		Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	if policy == nil {
		return []graphmodels.DeviceManagementConfigurationPolicyable{}, nil
	}

	return []graphmodels.DeviceManagementConfigurationPolicyable{policy}, nil
}

// getPoliciesByName retrieves policies matching an exact name
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) getPoliciesByName(ctx context.Context, object WindowsAutopilotDevicePreparationPolicyDataSourceModel) ([]graphmodels.DeviceManagementConfigurationPolicyable, error) {
	name := object.Name.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up device preparation policies by name: %s", name))

	filter := fmt.Sprintf("name eq '%s'", strings.ReplaceAll(name, "'", "''"))

	requestConfig := &devicemanagement.ConfigurationPoliciesRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.ConfigurationPoliciesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	policies, err := d.listPoliciesWithPageIterator(ctx, requestConfig)
	if err != nil {
		return nil, err
	}

	// The server side filter is a convenience only, the exact match is enforced locally so that a
	// case insensitive or partial server side match can never widen the result set.
	matched := make([]graphmodels.DeviceManagementConfigurationPolicyable, 0, len(policies))
	for _, policy := range policies {
		if policy.GetName() != nil && *policy.GetName() == name {
			matched = append(matched, policy)
		}
	}

	return matched, nil
}

// getPoliciesByODataQuery retrieves policies using a custom OData filter
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) getPoliciesByODataQuery(ctx context.Context, object WindowsAutopilotDevicePreparationPolicyDataSourceModel) ([]graphmodels.DeviceManagementConfigurationPolicyable, error) {
	filter := object.ODataQuery.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up device preparation policies with OData query: %s", filter))

	requestConfig := &devicemanagement.ConfigurationPoliciesRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.ConfigurationPoliciesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.listPoliciesWithPageIterator(ctx, requestConfig)
}

// listAllPolicies retrieves all configuration policies in the tenant
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) listAllPolicies(ctx context.Context) ([]graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Listing all device preparation policies")

	return d.listPoliciesWithPageIterator(ctx, nil)
}

// listPoliciesWithPageIterator handles pagination for configuration policy list requests
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) listPoliciesWithPageIterator(ctx context.Context, requestConfig *devicemanagement.ConfigurationPoliciesRequestBuilderGetRequestConfiguration) ([]graphmodels.DeviceManagementConfigurationPolicyable, error) {
	var allPolicies []graphmodels.DeviceManagementConfigurationPolicyable

	result, err := d.client.
		DeviceManagement().
		ConfigurationPolicies().
		Get(ctx, requestConfig)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return allPolicies, nil
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.DeviceManagementConfigurationPolicyable](
		result,
		d.client.GetAdapter(),
		graphmodels.CreateDeviceManagementConfigurationPolicyCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(policy graphmodels.DeviceManagementConfigurationPolicyable) bool {
		if policy != nil {
			allPolicies = append(allPolicies, policy)
		}
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate configuration policy pages: %w", err)
	}

	return allPolicies, nil
}

// getPolicyAssignments retrieves the assignments of a single policy
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) getPolicyAssignments(ctx context.Context, policyId string) ([]graphmodels.DeviceManagementConfigurationPolicyAssignmentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Looking up assignments for device preparation policy: %s", policyId))

	result, err := d.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(policyId).
		Assignments().
		Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return []graphmodels.DeviceManagementConfigurationPolicyAssignmentable{}, nil
	}

	return result.GetValue(), nil
}

// filterDevicePreparationPolicies discards policies that were not created from an Autopilot
// device preparation template.
func filterDevicePreparationPolicies(policies []graphmodels.DeviceManagementConfigurationPolicyable) []graphmodels.DeviceManagementConfigurationPolicyable {
	filtered := make([]graphmodels.DeviceManagementConfigurationPolicyable, 0, len(policies))
	for _, policy := range policies {
		if isDevicePreparationPolicy(policy) {
			filtered = append(filtered, policy)
		}
	}

	return filtered
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

	return deploymentModeFromTemplateId(templateRef.GetTemplateId()) != ""
}

// deploymentModeFromTemplateId maps an Autopilot device preparation template id to its deployment
// mode, returning an empty string when the template is not a device preparation template.
func deploymentModeFromTemplateId(templateId *string) string {
	if templateId == nil {
		return ""
	}

	// Template ids are versioned, e.g. "<guid>_1", so match on the base template guid.
	switch strings.SplitN(*templateId, "_", 2)[0] {
	case baseTemplateIDAutomatic:
		return "automatic"
	case baseTemplateIDUserDriven:
		return "user_driven"
	default:
		return ""
	}
}
