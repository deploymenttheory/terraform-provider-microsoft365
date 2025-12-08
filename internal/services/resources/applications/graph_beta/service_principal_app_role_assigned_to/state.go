package graphBetaServicePrincipalAppRoleAssignedTo

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the properties of an AppRoleAssignment to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *ServicePrincipalAppRoleAssignedToResourceModel, remoteResource graphmodels.AppRoleAssignmentable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"assignmentID": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())

	if resourceId := remoteResource.GetResourceId(); resourceId != nil {
		data.ResourceObjectID = convert.GraphToFrameworkUUID(resourceId)
	}

	if appRoleId := remoteResource.GetAppRoleId(); appRoleId != nil {
		data.AppRoleID = convert.GraphToFrameworkUUID(appRoleId)
	}

	if principalId := remoteResource.GetPrincipalId(); principalId != nil {
		data.TargetServicePrincipalObjectID = convert.GraphToFrameworkUUID(principalId)
	}

	data.PrincipalType = convert.GraphToFrameworkString(remoteResource.GetPrincipalType())
	data.PrincipalDisplayName = convert.GraphToFrameworkString(remoteResource.GetPrincipalDisplayName())
	data.ResourceDisplayName = convert.GraphToFrameworkString(remoteResource.GetResourceDisplayName())

	if creationTimestamp := remoteResource.GetCreationTimestamp(); creationTimestamp != nil {
		data.CreatedDateTime = convert.GraphToFrameworkTime(creationTimestamp)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state %s with id %s", ResourceName, data.ID.ValueString()))
}
