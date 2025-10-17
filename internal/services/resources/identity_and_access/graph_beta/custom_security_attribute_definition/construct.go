package graphBetaCustomSecurityAttributeDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs the CustomSecurityAttributeDefinition resource for Microsoft Graph API calls
func constructResource(ctx context.Context, data *CustomSecurityAttributeDefinitionResourceModel) (graphmodels.CustomSecurityAttributeDefinitionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewCustomSecurityAttributeDefinition()

	convert.FrameworkToGraphString(data.AttributeSet, requestBody.SetAttributeSet)
	convert.FrameworkToGraphString(data.Name, requestBody.SetName)
	convert.FrameworkToGraphBool(data.IsCollection, requestBody.SetIsCollection)
	convert.FrameworkToGraphBool(data.IsSearchable, requestBody.SetIsSearchable)
	convert.FrameworkToGraphString(data.Status, requestBody.SetStatus)
	convert.FrameworkToGraphString(data.Type, requestBody.SetTypeEscaped)
	convert.FrameworkToGraphBool(data.UsePreDefinedValuesOnly, requestBody.SetUsePreDefinedValuesOnly)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
