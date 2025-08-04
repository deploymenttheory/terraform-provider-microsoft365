package graphBetaRoleScopeTag

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a role scope tag resource for the Microsoft Graph API
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *RoleScopeTagResourceModel, isUpdate bool) (graphmodels.RoleScopeTagable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		var excludeResourceID *string
		// For updates, exclude the current resource from validation
		if isUpdate && !data.ID.IsNull() && !data.ID.IsUnknown() {
			id := data.ID.ValueString()
			excludeResourceID = &id
		}
		if err := validateRequest(ctx, client, data.DisplayName.ValueString(), excludeResourceID); err != nil {
			return nil, err
		}
	}

	requestBody := graphmodels.NewRoleScopeTag()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
