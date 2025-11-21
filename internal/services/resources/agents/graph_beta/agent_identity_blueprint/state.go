package graphBetaAgentsAgentIdentityBlueprint

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state from Microsoft Graph to the Terraform state model.
func MapRemoteStateToTerraform(ctx context.Context, data *AgentIdentityBlueprintResourceModel, remoteResource graphmodels.Applicationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	// Core identity fields
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.AppId = convert.GraphToFrameworkString(remoteResource.GetAppId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.UniqueName = convert.GraphToFrameworkString(remoteResource.GetUniqueName())

	// Configuration fields
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.SignInAudience = convert.GraphToFrameworkString(remoteResource.GetSignInAudience())
	data.GroupMembershipClaims = convert.GraphToFrameworkString(remoteResource.GetGroupMembershipClaims())
	data.PublisherDomain = convert.GraphToFrameworkString(remoteResource.GetPublisherDomain())
	data.ServiceManagementReference = convert.GraphToFrameworkString(remoteResource.GetServiceManagementReference())

	// TokenEncryptionKeyId returns a UUID, not a string, so we need to handle it differently
	if tokenKeyId := remoteResource.GetTokenEncryptionKeyId(); tokenKeyId != nil {
		data.TokenEncryptionKeyId = types.StringValue(tokenKeyId.String())
	} else {
		data.TokenEncryptionKeyId = types.StringNull()
	}

	// Read-only fields
	// Note: createdByAppId might not be available in the SDK yet
	// data.CreatedByAppId = convert.GraphToFrameworkString(remoteResource.GetCreatedByAppId())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.DisabledByMicrosoftStatus = convert.GraphToFrameworkString(remoteResource.GetDisabledByMicrosoftStatus())

	// Collection fields
	if identifierUris := remoteResource.GetIdentifierUris(); identifierUris != nil {
		data.IdentifierUris, _ = types.SetValueFrom(ctx, types.StringType, identifierUris)
	} else {
		data.IdentifierUris = types.SetNull(types.StringType)
	}

	if tags := remoteResource.GetTags(); tags != nil {
		data.Tags, _ = types.SetValueFrom(ctx, types.StringType, tags)
	} else {
		data.Tags = types.SetNull(types.StringType)
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"resourceId": data.ID.ValueString(),
	})
}

// MapRemoteStateToTerraformFromJSON maps the remote state from Microsoft Graph JSON response to the Terraform state model.
// This is used when working with HTTP client responses instead of SDK models.
func MapRemoteStateToTerraformFromJSON(ctx context.Context, data *AgentIdentityBlueprintResourceModel, remoteResource map[string]any) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote JSON state to Terraform state")

	// Helper function to safely get string values from JSON
	getString := func(key string) types.String {
		if val, ok := remoteResource[key].(string); ok && val != "" {
			return types.StringValue(val)
		}
		return types.StringNull()
	}

	// Helper function to safely get string slice values from JSON
	getStringSet := func(key string) types.Set {
		if val, ok := remoteResource[key].([]any); ok && len(val) > 0 {
			stringVals := make([]string, 0, len(val))
			for _, v := range val {
				if str, ok := v.(string); ok {
					stringVals = append(stringVals, str)
				}
			}
			if len(stringVals) > 0 {
				set, _ := types.SetValueFrom(ctx, types.StringType, stringVals)
				return set
			}
		}
		return types.SetNull(types.StringType)
	}

	// Core identity fields
	data.ID = getString("id")
	data.AppId = getString("appId")
	data.DisplayName = getString("displayName")
	data.UniqueName = getString("uniqueName")

	// Configuration fields
	data.Description = getString("description")
	data.SignInAudience = getString("signInAudience")
	data.GroupMembershipClaims = getString("groupMembershipClaims")
	data.PublisherDomain = getString("publisherDomain")
	data.ServiceManagementReference = getString("serviceManagementReference")

	// TokenEncryptionKeyId - may be returned as string or UUID
	if val, ok := remoteResource["tokenEncryptionKeyId"].(string); ok && val != "" {
		data.TokenEncryptionKeyId = types.StringValue(val)
	} else {
		data.TokenEncryptionKeyId = types.StringNull()
	}

	// Read-only fields
	// Note: createdByAppId might not be in the response
	data.CreatedByAppId = getString("createdByAppId")
	data.CreatedDateTime = getString("createdDateTime")
	data.DisabledByMicrosoftStatus = getString("disabledByMicrosoftStatus")

	// Collection fields
	data.IdentifierUris = getStringSet("identifierUris")
	data.Tags = getStringSet("tags")

	tflog.Debug(ctx, "Finished mapping remote JSON state to Terraform state", map[string]any{
		"resourceId": data.ID.ValueString(),
	})
}
