package graphBetaAuthenticationStrength

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MapRemoteResourceStateToTerraform maps the remote authentication strength policy to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AuthenticationStrengthResourceModel, remoteResource map[string]interface{}) {
	if id, ok := remoteResource["id"].(string); ok {
		data.ID = types.StringValue(id)
	}

	data.DisplayName = convert.MapToFrameworkString(remoteResource, "displayName")
	data.Description = convert.MapToFrameworkString(remoteResource, "description")
	data.PolicyType = convert.MapToFrameworkString(remoteResource, "policyType")
	data.RequirementsSatisfied = convert.MapToFrameworkString(remoteResource, "requirementsSatisfied")
	data.CreatedDateTime = convert.MapToFrameworkString(remoteResource, "createdDateTime")
	data.ModifiedDateTime = convert.MapToFrameworkString(remoteResource, "modifiedDateTime")
	data.AllowedCombinations = convert.MapToFrameworkStringSet(ctx, remoteResource, "allowedCombinations")

}
