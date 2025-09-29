package graphBetaGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of a Group resource to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *GroupResourceModel, remoteResource graphmodels.Groupable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.MailNickname = convert.GraphToFrameworkString(remoteResource.GetMailNickname())
	data.MailEnabled = convert.GraphToFrameworkBool(remoteResource.GetMailEnabled())
	data.SecurityEnabled = convert.GraphToFrameworkBool(remoteResource.GetSecurityEnabled())
	data.GroupTypes = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetGroupTypes())
	data.Visibility = convert.GraphToFrameworkString(remoteResource.GetVisibility())
	data.IsAssignableToRole = convert.GraphToFrameworkBool(remoteResource.GetIsAssignableToRole())
	data.MembershipRule = convert.GraphToFrameworkString(remoteResource.GetMembershipRule())
	data.MembershipRuleProcessingState = convert.GraphToFrameworkString(remoteResource.GetMembershipRuleProcessingState())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.Mail = convert.GraphToFrameworkString(remoteResource.GetMail())
	data.ProxyAddresses = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetProxyAddresses())
	data.OnPremisesSyncEnabled = convert.GraphToFrameworkBool(remoteResource.GetOnPremisesSyncEnabled())
	data.PreferredDataLocation = convert.GraphToFrameworkString(remoteResource.GetPreferredDataLocation())
	data.PreferredLanguage = convert.GraphToFrameworkString(remoteResource.GetPreferredLanguage())
	data.Theme = convert.GraphToFrameworkString(remoteResource.GetTheme())
	data.Classification = convert.GraphToFrameworkString(remoteResource.GetClassification())
	data.ExpirationDateTime = convert.GraphToFrameworkTime(remoteResource.GetExpirationDateTime())
	data.RenewedDateTime = convert.GraphToFrameworkTime(remoteResource.GetRenewedDateTime())
	data.SecurityIdentifier = convert.GraphToFrameworkString(remoteResource.GetSecurityIdentifier())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state %s with id %s", ResourceName, data.ID.ValueString()))
}
