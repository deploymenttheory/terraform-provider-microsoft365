package graphBetaServicePrincipal

import (
	"context"
	"encoding/base64"
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

	// Map boolean fields
	data.AccountEnabled = convert.GraphToFrameworkBool(remoteResource.GetAccountEnabled())
	data.AppRoleAssignmentRequired = convert.GraphToFrameworkBool(remoteResource.GetAppRoleAssignmentRequired())
	data.IsDisabled = convert.GraphToFrameworkBool(remoteResource.GetIsDisabled())

	// Map optional string fields
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.LoginURL = convert.GraphToFrameworkString(remoteResource.GetLoginUrl())
	data.Notes = convert.GraphToFrameworkString(remoteResource.GetNotes())
	data.PreferredSingleSignOnMode = convert.GraphToFrameworkString(remoteResource.GetPreferredSingleSignOnMode())
	data.TokenEncryptionKeyID = convert.GraphToFrameworkUUID(remoteResource.GetTokenEncryptionKeyId())

	// Map computed string fields
	data.ServicePrincipalType = convert.GraphToFrameworkString(remoteResource.GetServicePrincipalType())
	data.AppOwnerOrganizationID = convert.GraphToFrameworkUUID(remoteResource.GetAppOwnerOrganizationId())
	data.ApplicationTemplateID = convert.GraphToFrameworkString(remoteResource.GetApplicationTemplateId())
	data.CreatedByAppID = convert.GraphToFrameworkString(remoteResource.GetCreatedByAppId())
	data.PreferredTokenSigningKeyEndDateTime = convert.GraphToFrameworkTime(remoteResource.GetPreferredTokenSigningKeyEndDateTime())
	data.PreferredTokenSigningKeyThumbprint = convert.GraphToFrameworkString(remoteResource.GetPreferredTokenSigningKeyThumbprint())

	// Map computed collections
	data.KeyCredentials = mapKeyCredentials(remoteResource.GetKeyCredentials())
	data.PasswordCredentials = mapPasswordCredentials(remoteResource.GetPasswordCredentials())

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

// bytesToBase64String encodes a binary Graph value (e.g. customKeyIdentifier) as base64,
// matching how the Graph API represents the value on the wire.
func bytesToBase64String(value []byte) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(base64.StdEncoding.EncodeToString(value))
}

func mapKeyCredentials(keyCredentials []graphmodels.KeyCredentialable) types.Set {
	elemType := types.ObjectType{AttrTypes: keyCredentialAttrTypes}
	if len(keyCredentials) == 0 {
		emptySet, _ := types.SetValue(elemType, []attr.Value{})
		return emptySet
	}

	elements := make([]attr.Value, 0, len(keyCredentials))
	for _, keyCred := range keyCredentials {
		if keyCred == nil {
			continue
		}

		attrs := map[string]attr.Value{
			"custom_key_identifier": bytesToBase64String(keyCred.GetCustomKeyIdentifier()),
			"display_name":          convert.GraphToFrameworkString(keyCred.GetDisplayName()),
			"end_date_time":         convert.GraphToFrameworkTime(keyCred.GetEndDateTime()),
			"key_id":                convert.GraphToFrameworkUUID(keyCred.GetKeyId()),
			"start_date_time":       convert.GraphToFrameworkTime(keyCred.GetStartDateTime()),
			"type":                  convert.GraphToFrameworkString(keyCred.GetTypeEscaped()),
			"usage":                 convert.GraphToFrameworkString(keyCred.GetUsage()),
		}

		objVal, diags := types.ObjectValue(keyCredentialAttrTypes, attrs)
		if diags.HasError() {
			continue
		}
		elements = append(elements, objVal)
	}

	setVal, diags := types.SetValue(elemType, elements)
	if diags.HasError() {
		return types.SetNull(elemType)
	}
	return setVal
}

func mapPasswordCredentials(passwordCredentials []graphmodels.PasswordCredentialable) types.Set {
	elemType := types.ObjectType{AttrTypes: passwordCredentialAttrTypes}
	if len(passwordCredentials) == 0 {
		emptySet, _ := types.SetValue(elemType, []attr.Value{})
		return emptySet
	}

	elements := make([]attr.Value, 0, len(passwordCredentials))
	for _, passCred := range passwordCredentials {
		if passCred == nil {
			continue
		}

		attrs := map[string]attr.Value{
			"custom_key_identifier": bytesToBase64String(passCred.GetCustomKeyIdentifier()),
			"display_name":          convert.GraphToFrameworkString(passCred.GetDisplayName()),
			"end_date_time":         convert.GraphToFrameworkTime(passCred.GetEndDateTime()),
			"hint":                  convert.GraphToFrameworkString(passCred.GetHint()),
			"key_id":                convert.GraphToFrameworkUUID(passCred.GetKeyId()),
			"start_date_time":       convert.GraphToFrameworkTime(passCred.GetStartDateTime()),
		}

		objVal, diags := types.ObjectValue(passwordCredentialAttrTypes, attrs)
		if diags.HasError() {
			continue
		}
		elements = append(elements, objVal)
	}

	setVal, diags := types.SetValue(elemType, elements)
	if diags.HasError() {
		return types.SetNull(elemType)
	}
	return setVal
}
