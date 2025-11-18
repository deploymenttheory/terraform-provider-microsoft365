package graphBetaNetworkFilteringPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *NetworkFilteringPolicyResourceModel) (models.FilteringPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := models.NewFilteringPolicy()

	convert.FrameworkToGraphString(data.Name, requestBody.SetName)

	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphEnum(data.Action, models.ParseFilteringPolicyAction, requestBody.SetAction); err != nil {
		return nil, fmt.Errorf("invalid filtering policy action: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
