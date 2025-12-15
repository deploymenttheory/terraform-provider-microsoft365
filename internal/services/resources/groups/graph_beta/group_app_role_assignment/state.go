package graphBetaGroupAppRoleAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of a GroupAppRoleAssignment resource to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *GroupAppRoleAssignmentResourceModel, remoteResource graphmodels.AppRoleAssignmentable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.TargetGroupID = convert.GraphToFrameworkUUID(remoteResource.GetPrincipalId())
	data.ResourceObjectID = convert.GraphToFrameworkUUID(remoteResource.GetResourceId())
	data.AppRoleID = convert.GraphToFrameworkUUID(remoteResource.GetAppRoleId())
	data.PrincipalDisplayName = convert.GraphToFrameworkString(remoteResource.GetPrincipalDisplayName())
	data.ResourceDisplayName = convert.GraphToFrameworkString(remoteResource.GetResourceDisplayName())
	data.PrincipalType = convert.GraphToFrameworkString(remoteResource.GetPrincipalType())

	if creationTime := remoteResource.GetCreationTimestamp(); creationTime != nil {
		timestampStr := creationTime.String()
		data.CreationTimestamp = convert.GraphToFrameworkString(&timestampStr)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state %s with id %s", ResourceName, data.ID.ValueString()))
}
