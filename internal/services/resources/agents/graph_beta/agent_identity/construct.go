package graphBetaAgentIdentity

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource builds the request body for creating or updating an agent identity.
// For create (isCreate=true): returns a ServicePrincipal with @odata.type set to agentIdentity
// For update (isCreate=false): returns a ServicePrincipal without the create-only fields
func constructResource(ctx context.Context, data *AgentIdentityResourceModel, isCreate bool) (graphmodels.ServicePrincipalable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model (isCreate: %t)", ResourceName, isCreate))

	requestBody := graphmodels.NewAgentIdentity()

	if isCreate {
		// Set @odata.type to specify this is an agentIdentity
		// odataType := "#microsoft.graph.agentIdentity"
		// requestBody.SetOdataType(&odataType)

		// Set the blueprint ID - required for creation
		additionalData := make(map[string]any)
		additionalData["agentIdentityBlueprintId"] = data.AgentIdentityBlueprintId.ValueString()

		// Add sponsors using OData bind - only during creation
		// Separate API calls handle add/delete for updates
		if !data.SponsorIds.IsNull() && !data.SponsorIds.IsUnknown() {
			var sponsors []string
			diags := data.SponsorIds.ElementsAs(ctx, &sponsors, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract sponsor_ids: %v", diags.Errors())
			}
			if len(sponsors) > 0 {
				sponsorUrls := make([]string, len(sponsors))
				for i, sponsorId := range sponsors {
					sponsorUrls[i] = fmt.Sprintf("https://graph.microsoft.com/beta/users/%s", sponsorId)
				}
				additionalData["sponsors@odata.bind"] = sponsorUrls
				tflog.Debug(ctx, fmt.Sprintf("Adding %d sponsors to agent identity", len(sponsors)))
			}
		}

		// Add owners using OData bind - only during creation
		// Separate API calls handle add/delete for updates
		if !data.OwnerIds.IsNull() && !data.OwnerIds.IsUnknown() {
			var owners []string
			diags := data.OwnerIds.ElementsAs(ctx, &owners, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract owner_ids: %v", diags.Errors())
			}
			if len(owners) > 0 {
				ownerUrls := make([]string, len(owners))
				for i, ownerId := range owners {
					ownerUrls[i] = fmt.Sprintf("https://graph.microsoft.com/beta/users/%s", ownerId)
				}
				additionalData["owners@odata.bind"] = ownerUrls
				tflog.Debug(ctx, fmt.Sprintf("Adding %d owners to agent identity", len(owners)))
			}
		}

		requestBody.SetAdditionalData(additionalData)
	}

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphBool(data.AccountEnabled, requestBody.SetAccountEnabled)

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
