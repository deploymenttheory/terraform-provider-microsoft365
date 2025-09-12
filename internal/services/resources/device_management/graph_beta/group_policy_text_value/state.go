package graphBetaGroupPolicyTextValue

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote GroupPolicyPresentationValueText and GroupPolicyDefinitionValue to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data *GroupPolicyTextValueResourceModel, remoteResource graphmodels.GroupPolicyPresentationValueTextable, definitionValue graphmodels.GroupPolicyDefinitionValueable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId":    remoteResource.GetId(),
		"resourceValue": remoteResource.GetValue(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Value = convert.GraphToFrameworkString(remoteResource.GetValue())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	// Map the enabled state from the definition value
	if definitionValue != nil {
		data.Enabled = convert.GraphToFrameworkBool(definitionValue.GetEnabled())
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
