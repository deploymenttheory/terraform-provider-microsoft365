package graphBetaRoleDefinitionAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps a remote role assignment to the Terraform resource model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *RoleDefinitionAssignmentResourceModel, assignment graphmodels.DeviceAndAppManagementRoleAssignmentable) {
	if assignment == nil {
		tflog.Debug(ctx, "Remote assignment is nil")
		return
	}

	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]interface{}{
		"assignmentId": convert.GraphToFrameworkString(assignment.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(assignment.GetId())
	data.DisplayName = convert.GraphToFrameworkString(assignment.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(assignment.GetDescription())
	data.ScopeType = convert.GraphToFrameworkEnum(assignment.GetScopeType())

	if members := assignment.GetScopeMembers(); len(members) > 0 {
		data.ScopeMembers = convert.GraphToFrameworkStringSet(ctx, members)
	} else {
		data.ScopeMembers = types.SetNull(types.StringType)
	}

	// Set resource scopes
	if scopes := assignment.GetResourceScopes(); len(scopes) > 0 {
		data.ResourceScopes = convert.GraphToFrameworkStringSet(ctx, scopes)
	} else {
		data.ResourceScopes = types.SetNull(types.StringType)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
