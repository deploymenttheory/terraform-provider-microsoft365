package graphBetaApplicationsAgentIdentityBlueprintServicePrincipal

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *AgentIdentityBlueprintServicePrincipalResourceModel) (graphmodels.ServicePrincipalable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewServicePrincipal()

	odataType := "#microsoft.graph.agentIdentityBlueprintPrincipal"
	requestBody.SetOdataType(&odataType)

	appId := data.AppId.ValueString()
	requestBody.SetAppId(&appId)

	if err := convert.FrameworkToGraphStringSet(ctx, data.Tags, requestBody.SetTags); err != nil {
		return nil, fmt.Errorf("failed to set tags: %w", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
