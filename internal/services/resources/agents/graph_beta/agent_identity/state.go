package graphBetaAgentIdentity

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *AgentIdentityResourceModel, remoteResource graphmodels.ServicePrincipalable) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.AccountEnabled = convert.GraphToFrameworkBool(remoteResource.GetAccountEnabled())
	data.ServicePrincipalType = convert.GraphToFrameworkString(remoteResource.GetServicePrincipalType())
	data.Tags = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetTags())
	data.DisabledByMicrosoftStatus = convert.GraphToFrameworkString(remoteResource.GetDisabledByMicrosoftStatus())

	// Get additional data for agent identity specific fields
	additionalData := remoteResource.GetAdditionalData()
	if additionalData != nil {
		if blueprintId, ok := additionalData["agentIdentityBlueprintId"].(string); ok {
			data.AgentIdentityBlueprintId = types.StringValue(blueprintId)
		}
		if createdByAppId, ok := additionalData["createdByAppId"].(string); ok {
			data.CreatedByAppId = types.StringValue(createdByAppId)
		} else {
			data.CreatedByAppId = types.StringNull()
		}
		if createdDateTime, ok := additionalData["createdDateTime"].(string); ok {
			data.CreatedDateTime = types.StringValue(createdDateTime)
		} else {
			data.CreatedDateTime = types.StringNull()
		}
	} else {
		// Ensure computed fields are set to null if additionalData is nil
		data.CreatedByAppId = types.StringNull()
		data.CreatedDateTime = types.StringNull()
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
