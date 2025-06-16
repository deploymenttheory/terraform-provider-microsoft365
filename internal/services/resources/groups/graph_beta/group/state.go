package graphBetaGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of a Group resource to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *GroupResourceModel, remoteResource graphmodels.Groupable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": types.StringPointerValue(remoteResource.GetId()),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.MailNickname = types.StringPointerValue(remoteResource.GetMailNickname())

	if mailEnabled := remoteResource.GetMailEnabled(); mailEnabled != nil {
		data.MailEnabled = types.BoolValue(*mailEnabled)
	}

	if securityEnabled := remoteResource.GetSecurityEnabled(); securityEnabled != nil {
		data.SecurityEnabled = types.BoolValue(*securityEnabled)
	}

	data.GroupTypes = state.StringSliceToSet(ctx, remoteResource.GetGroupTypes())
	data.Visibility = types.StringPointerValue(remoteResource.GetVisibility())

	if isAssignableToRole := remoteResource.GetIsAssignableToRole(); isAssignableToRole != nil {
		data.IsAssignableToRole = types.BoolValue(*isAssignableToRole)
	}

	data.MembershipRule = types.StringPointerValue(remoteResource.GetMembershipRule())
	data.MembershipRuleProcessingState = types.StringPointerValue(remoteResource.GetMembershipRuleProcessingState())

	if createdDateTime := remoteResource.GetCreatedDateTime(); createdDateTime != nil {
		data.CreatedDateTime = types.StringValue(createdDateTime.Format("2006-01-02T15:04:05Z"))
	}

	data.Mail = types.StringPointerValue(remoteResource.GetMail())
	data.ProxyAddresses = state.StringSliceToSet(ctx, remoteResource.GetProxyAddresses())

	if onPremisesSyncEnabled := remoteResource.GetOnPremisesSyncEnabled(); onPremisesSyncEnabled != nil {
		data.OnPremisesSyncEnabled = types.BoolValue(*onPremisesSyncEnabled)
	}

	data.PreferredDataLocation = types.StringPointerValue(remoteResource.GetPreferredDataLocation())
	data.PreferredLanguage = types.StringPointerValue(remoteResource.GetPreferredLanguage())
	data.Theme = types.StringPointerValue(remoteResource.GetTheme())
	data.Classification = types.StringPointerValue(remoteResource.GetClassification())

	if expirationDateTime := remoteResource.GetExpirationDateTime(); expirationDateTime != nil {
		data.ExpirationDateTime = types.StringValue(expirationDateTime.Format("2006-01-02T15:04:05Z"))
	}

	if renewedDateTime := remoteResource.GetRenewedDateTime(); renewedDateTime != nil {
		data.RenewedDateTime = types.StringValue(renewedDateTime.Format("2006-01-02T15:04:05Z"))
	}

	data.SecurityIdentifier = types.StringPointerValue(remoteResource.GetSecurityIdentifier())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state %s with id %s", ResourceName, data.ID.ValueString()))
}
