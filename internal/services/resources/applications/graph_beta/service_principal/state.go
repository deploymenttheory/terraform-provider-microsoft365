package graphBetaServicePrincipal

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state from Microsoft Graph API to the Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data ServicePrincipalResourceModel, remoteResource graphmodels.ServicePrincipalable) ServicePrincipalResourceModel {
	tflog.Debug(ctx, fmt.Sprintf("Mapping %s remote state to Terraform state", ResourceName))

	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.AppID = convert.GraphToFrameworkString(remoteResource.GetAppId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.AccountEnabled = convert.GraphToFrameworkBool(remoteResource.GetAccountEnabled())
	data.AppRoleAssignmentRequired = convert.GraphToFrameworkBool(remoteResource.GetAppRoleAssignmentRequired())
	data.ServicePrincipalType = convert.GraphToFrameworkString(remoteResource.GetServicePrincipalType())
	data.ServicePrincipalNames = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetServicePrincipalNames())
	data.SignInAudience = convert.GraphToFrameworkString(remoteResource.GetSignInAudience())
	data.Tags = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetTags())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s remote state to Terraform state", ResourceName))

	return data
}
