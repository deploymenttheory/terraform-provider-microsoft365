package graphBetaTenantWideGroupSettings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *TenantWideGroupSettingsResourceModel) (graphmodels.DirectorySettingable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDirectorySetting()

	convert.FrameworkToGraphString(data.TemplateID, requestBody.SetTemplateId)

	// Convert values set to setting values
	if !data.Values.IsNull() && !data.Values.IsUnknown() {
		var settingValues []types.Object
		diags := data.Values.ElementsAs(ctx, &settingValues, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert values set: %v", diags.Errors())
		}

		var values []graphmodels.SettingValueable
		for _, settingValueObj := range settingValues {
			var settingValue SettingValueModel
			diags := settingValueObj.As(ctx, &settingValue, basetypes.ObjectAsOptions{})
			if diags.HasError() {
				return nil, fmt.Errorf("failed to convert setting value object: %v", diags.Errors())
			}

			value := graphmodels.NewSettingValue()
			convert.FrameworkToGraphString(settingValue.Name, value.SetName)
			convert.FrameworkToGraphString(settingValue.Value, value.SetValue)
			values = append(values, value)
		}
		requestBody.SetValues(values)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
