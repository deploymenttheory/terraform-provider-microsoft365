// read.go (updated with proper technology enum handling)
package graphBetaLinuxPlatformScript

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for Linux Platform Scripts data source.
func (d *LinuxPlatformScriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object LinuxPlatformScriptDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with filter_type: %s", d.ProviderTypeName, d.TypeName, filterType))

	if filterType != "all" && (object.FilterValue.IsNull() || object.FilterValue.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			fmt.Sprintf("filter_value must be provided when filter_type is '%s'", filterType),
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var filteredItems []LinuxPlatformScriptModel
	filterValue := object.FilterValue.ValueString()

	// For ID filter, we can make a direct API call
	if filterType == "id" {
		constants.GraphSDKMutex.Lock()
		respItem, err := d.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(filterValue).
			Get(ctx, nil)
		constants.GraphSDKMutex.Unlock()

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		// Verify this is a Linux platform script
		if hasLinuxMdmTechnology(respItem) {
			filteredItems = append(filteredItems, MapRemoteStateToDataSource(respItem))
		} else {
			resp.Diagnostics.AddError(
				"Error Reading Linux Platform Script",
				fmt.Sprintf("The configuration policy with ID %s is not a Linux platform script", filterValue),
			)
			return
		}
	} else {
		// For all other filters, we need to get all Linux scripts and filter locally
		// Set up technology filter for Linux platform scripts
		technologyFilter := "technologies/any(t:t eq 'linuxMdm')"

		requestOptions := &devicemanagement.ConfigurationPoliciesRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.ConfigurationPoliciesRequestBuilderGetQueryParameters{
				Filter: &technologyFilter,
			},
		}

		constants.GraphSDKMutex.Lock()
		respList, err := d.client.
			DeviceManagement().
			ConfigurationPolicies().
			Get(ctx, requestOptions)
		constants.GraphSDKMutex.Unlock()

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, item := range respList.GetValue() {
			// Ensure this is a Linux platform script
			if !hasLinuxMdmTechnology(item) {
				continue
			}

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))

			case "display_name":
				if item.GetName() != nil && strings.Contains(
					strings.ToLower(*item.GetName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))
				}
			}
		}
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s, found %d items", d.ProviderTypeName, d.TypeName, len(filteredItems)))
}

// hasLinuxMdmTechnology checks if a configuration policy has the Linux MDM technology
func hasLinuxMdmTechnology(policy graphmodels.DeviceManagementConfigurationPolicyable) bool {
	techEnum := policy.GetTechnologies()
	if techEnum == nil {
		return false
	}
	techString := (*techEnum).String()

	for _, tech := range strings.Split(techString, ",") {
		if tech == "linuxMdm" {
			return true
		}
	}

	return false
}
