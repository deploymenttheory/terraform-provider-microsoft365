package graphBetaApplicationsAgentIdentityBlueprint

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *AgentIdentityBlueprintResourceModel, isCreate bool) (graphmodels.Applicationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewApplication()

	// Set @odata.type to specify this is an agentIdentityBlueprint
	odataType := "#microsoft.graph.agentIdentityBlueprint"
	requestBody.SetOdataType(&odataType)

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.SignInAudience, requestBody.SetSignInAudience)

	if err := convert.FrameworkToGraphStringSet(ctx, data.Tags, requestBody.SetTags); err != nil {
		return nil, fmt.Errorf("failed to set tags: %w", err)
	}

	// Set sponsors and owners using OData bind properties - only during creation
	// seperate constructor for update scenarios that requires a separate set of api endpoints.
	if isCreate {
		additionalData := requestBody.GetAdditionalData()
		if additionalData == nil {
			additionalData = make(map[string]any)
		}

		if !data.SponsorUserIds.IsNull() && !data.SponsorUserIds.IsUnknown() {
			var sponsors []string
			diags := data.SponsorUserIds.ElementsAs(ctx, &sponsors, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract sponsor_user_ids: %v", diags.Errors())
			}
			if len(sponsors) > 0 {
				sponsorUrls := make([]string, len(sponsors))
				for i, sponsorId := range sponsors {
					sponsorUrls[i] = fmt.Sprintf("https://graph.microsoft.com/beta/users/%s", sponsorId)
				}
				additionalData["sponsors@odata.bind"] = sponsorUrls
				tflog.Debug(ctx, fmt.Sprintf("Adding %d sponsors to agent identity blueprint", len(sponsors)))
			}
		}

		if !data.OwnerUserIds.IsNull() && !data.OwnerUserIds.IsUnknown() {
			var owners []string
			diags := data.OwnerUserIds.ElementsAs(ctx, &owners, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract owner_user_ids: %v", diags.Errors())
			}
			if len(owners) > 0 {
				ownerUrls := make([]string, len(owners))
				for i, ownerId := range owners {
					ownerUrls[i] = fmt.Sprintf("https://graph.microsoft.com/beta/users/%s", ownerId)
				}
				additionalData["owners@odata.bind"] = ownerUrls
				tflog.Debug(ctx, fmt.Sprintf("Adding %d owners to agent identity blueprint", len(owners)))
			}
		}

		requestBody.SetAdditionalData(additionalData)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
