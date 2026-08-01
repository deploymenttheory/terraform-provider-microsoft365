package graphBetaServicePrincipal

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// samlSingleSignOnSettingsAttrTypes defines the attribute types for the saml_single_sign_on_settings nested object
var samlSingleSignOnSettingsAttrTypes = map[string]attr.Type{
	"relay_state": types.StringType,
}

// MapRemoteStateToTerraform maps the remote state from Microsoft Graph API to the Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data ServicePrincipalResourceModel, remoteResource graphmodels.ServicePrincipalable) ServicePrincipalResourceModel {
	tflog.Debug(ctx, fmt.Sprintf("Mapping %s remote state to Terraform state", ResourceName))

	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	// Map basic fields using helpers
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.AppID = convert.GraphToFrameworkString(remoteResource.GetAppId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())

	// Map boolean fields
	data.AccountEnabled = convert.GraphToFrameworkBool(remoteResource.GetAccountEnabled())
	data.AppRoleAssignmentRequired = convert.GraphToFrameworkBool(remoteResource.GetAppRoleAssignmentRequired())

	// Map optional string fields
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Homepage = convert.GraphToFrameworkString(remoteResource.GetHomepage())
	data.LoginURL = convert.GraphToFrameworkString(remoteResource.GetLoginUrl())
	data.LogoutURL = convert.GraphToFrameworkString(remoteResource.GetLogoutUrl())
	data.Notes = convert.GraphToFrameworkString(remoteResource.GetNotes())
	data.PreferredSingleSignOnMode = convert.GraphToFrameworkString(remoteResource.GetPreferredSingleSignOnMode())

	// Map computed string fields
	data.ServicePrincipalType = convert.GraphToFrameworkString(remoteResource.GetServicePrincipalType())
	data.SignInAudience = convert.GraphToFrameworkString(remoteResource.GetSignInAudience())

	// Map collection fields
	data.ServicePrincipalNames = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetServicePrincipalNames())

	// Filter tags to only include configured values (excludes system-generated tags)
	// This prevents drift when Microsoft adds system tags like "WindowsAzureActiveDirectoryIntegratedApp"
	data.Tags = convert.GraphToFrameworkStringSetFiltered(ctx, remoteResource.GetTags(), data.Tags)

	data.NotificationEmailAddresses = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetNotificationEmailAddresses())

	// GraphToFrameworkStringSet maps an empty remote list to null; keep a configured
	// empty set ([] clears the collection) so the planned value round-trips.
	remoteAlternativeNames := convert.GraphToFrameworkStringSet(ctx, remoteResource.GetAlternativeNames())
	configuredEmptySet := !data.AlternativeNames.IsNull() && !data.AlternativeNames.IsUnknown() && len(data.AlternativeNames.Elements()) == 0
	if !(remoteAlternativeNames.IsNull() && configuredEmptySet) {
		data.AlternativeNames = remoteAlternativeNames
	}

	// Map SAML single sign-on settings
	if samlSettings := remoteResource.GetSamlSingleSignOnSettings(); samlSettings != nil {
		samlAttrs := map[string]attr.Value{
			"relay_state": convert.GraphToFrameworkString(samlSettings.GetRelayState()),
		}
		data.SamlSingleSignOnSettings, _ = types.ObjectValue(samlSingleSignOnSettingsAttrTypes, samlAttrs)
	} else {
		data.SamlSingleSignOnSettings = types.ObjectNull(samlSingleSignOnSettingsAttrTypes)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s remote state to Terraform state", ResourceName))

	return data
}
