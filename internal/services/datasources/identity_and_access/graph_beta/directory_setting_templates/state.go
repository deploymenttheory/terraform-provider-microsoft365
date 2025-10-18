package graphBetaDirectorySettingTemplates

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a DirectorySettingTemplate to a model
func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.DirectorySettingTemplateable) DirectorySettingTemplateModel {
	model := DirectorySettingTemplateModel{
		ID:          convert.GraphToFrameworkString(data.GetId()),
		Description: convert.GraphToFrameworkString(data.GetDescription()),
		DisplayName: convert.GraphToFrameworkString(data.GetDisplayName()),
		Values:      []SettingTemplateValueModel{},
	}

	// Map the setting template values
	if values := data.GetValues(); values != nil {
		for _, value := range values {
			if value == nil {
				continue
			}

			valueModel := SettingTemplateValueModel{
				Name: convert.GraphToFrameworkString(value.GetName()),
				// The SDK uses GetTypeEscaped() for the type property
				Type:         convert.GraphToFrameworkString(value.GetTypeEscaped()),
				DefaultValue: convert.GraphToFrameworkString(value.GetDefaultValue()),
				Description:  convert.GraphToFrameworkString(value.GetDescription()),
			}

			model.Values = append(model.Values, valueModel)
		}
	}

	return model
}
