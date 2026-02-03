package graphBetaServicePrincipal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a service principal from the Graph API to the data source model
func MapRemoteStateToDataSource(ctx context.Context, servicePrincipal graphmodels.ServicePrincipalable, inputState ServicePrincipalDataSourceModel) ServicePrincipalDataSourceModel {
	model := ServicePrincipalDataSourceModel{
		Timeouts: inputState.Timeouts,
	}

	// Preserve input lookup attributes
	model.ObjectId = inputState.ObjectId
	model.AppId = inputState.AppId
	model.DisplayName = inputState.DisplayName
	model.ODataQuery = inputState.ODataQuery

	// Map ID
	if servicePrincipal.GetId() != nil {
		model.ID = types.StringValue(*servicePrincipal.GetId())
		// Also set object_id if it wasn't already set
		if model.ObjectId.IsNull() {
			model.ObjectId = types.StringValue(*servicePrincipal.GetId())
		}
	} else {
		model.ID = types.StringNull()
	}

	// Map AppID
	if servicePrincipal.GetAppId() != nil {
		if model.AppId.IsNull() {
			model.AppId = types.StringValue(*servicePrincipal.GetAppId())
		}
	}

	// Map DisplayName
	if servicePrincipal.GetDisplayName() != nil {
		if model.DisplayName.IsNull() {
			model.DisplayName = types.StringValue(*servicePrincipal.GetDisplayName())
		}
	} else {
		if model.DisplayName.IsNull() {
			model.DisplayName = types.StringNull()
		}
	}

	// Map AppDisplayName
	if servicePrincipal.GetAppDisplayName() != nil {
		model.AppDisplayName = types.StringValue(*servicePrincipal.GetAppDisplayName())
	} else {
		model.AppDisplayName = types.StringNull()
	}

	// Map DeletedDateTime
	if servicePrincipal.GetDeletedDateTime() != nil {
		model.DeletedDateTime = types.StringValue(servicePrincipal.GetDeletedDateTime().Format("2006-01-02T15:04:05Z"))
	} else {
		model.DeletedDateTime = types.StringNull()
	}

	// Map ApplicationTemplateID
	if servicePrincipal.GetApplicationTemplateId() != nil {
		model.ApplicationTemplateID = types.StringValue(*servicePrincipal.GetApplicationTemplateId())
	} else {
		model.ApplicationTemplateID = types.StringNull()
	}

	// Map AccountEnabled
	if servicePrincipal.GetAccountEnabled() != nil {
		model.AccountEnabled = types.BoolValue(*servicePrincipal.GetAccountEnabled())
	} else {
		model.AccountEnabled = types.BoolNull()
	}

	// Map AppRoleAssignmentRequired
	if servicePrincipal.GetAppRoleAssignmentRequired() != nil {
		model.AppRoleAssignmentRequired = types.BoolValue(*servicePrincipal.GetAppRoleAssignmentRequired())
	} else {
		model.AppRoleAssignmentRequired = types.BoolNull()
	}

	// Map ServicePrincipalType
	if servicePrincipal.GetServicePrincipalType() != nil {
		model.ServicePrincipalType = types.StringValue(*servicePrincipal.GetServicePrincipalType())
	} else {
		model.ServicePrincipalType = types.StringNull()
	}

	// Map SignInAudience
	if servicePrincipal.GetSignInAudience() != nil {
		model.SignInAudience = types.StringValue(*servicePrincipal.GetSignInAudience())
	} else {
		model.SignInAudience = types.StringNull()
	}

	// Map PreferredSingleSignOnMode
	if servicePrincipal.GetPreferredSingleSignOnMode() != nil {
		model.PreferredSingleSignOnMode = types.StringValue(*servicePrincipal.GetPreferredSingleSignOnMode())
	} else {
		model.PreferredSingleSignOnMode = types.StringNull()
	}

	// Map Homepage
	if servicePrincipal.GetHomepage() != nil {
		model.Homepage = types.StringValue(*servicePrincipal.GetHomepage())
	} else {
		model.Homepage = types.StringNull()
	}

	// Map ErrorUrl
	if servicePrincipal.GetErrorUrl() != nil {
		model.ErrorUrl = types.StringValue(*servicePrincipal.GetErrorUrl())
	} else {
		model.ErrorUrl = types.StringNull()
	}

	// Map PublisherName
	if servicePrincipal.GetPublisherName() != nil {
		model.PublisherName = types.StringValue(*servicePrincipal.GetPublisherName())
	} else {
		model.PublisherName = types.StringNull()
	}

	// Map reply URLs
	if replyUrls := servicePrincipal.GetReplyUrls(); len(replyUrls) > 0 {
		elements := make([]attr.Value, len(replyUrls))
		for i, url := range replyUrls {
			elements[i] = types.StringValue(url)
		}
		model.ReplyUrls, _ = types.SetValue(types.StringType, elements)
	} else {
		model.ReplyUrls = types.SetNull(types.StringType)
	}

	// Map service principal names
	if spNames := servicePrincipal.GetServicePrincipalNames(); len(spNames) > 0 {
		elements := make([]attr.Value, len(spNames))
		for i, name := range spNames {
			elements[i] = types.StringValue(name)
		}
		model.ServicePrincipalNames, _ = types.SetValue(types.StringType, elements)
	} else {
		model.ServicePrincipalNames = types.SetNull(types.StringType)
	}

	// Map tags
	if tags := servicePrincipal.GetTags(); len(tags) > 0 {
		elements := make([]attr.Value, len(tags))
		for i, tag := range tags {
			elements[i] = types.StringValue(tag)
		}
		model.Tags, _ = types.SetValue(types.StringType, elements)
	} else {
		model.Tags = types.SetNull(types.StringType)
	}

	// Map DisabledByMicrosoftStatus
	if servicePrincipal.GetDisabledByMicrosoftStatus() != nil {
		model.DisabledByMicrosoftStatus = types.StringValue(*servicePrincipal.GetDisabledByMicrosoftStatus())
	} else {
		model.DisabledByMicrosoftStatus = types.StringNull()
	}

	// Map AppOwnerOrganizationID
	if servicePrincipal.GetAppOwnerOrganizationId() != nil {
		model.AppOwnerOrganizationID = types.StringValue(servicePrincipal.GetAppOwnerOrganizationId().String())
	} else {
		model.AppOwnerOrganizationID = types.StringNull()
	}

	// Map LoginUrl
	if servicePrincipal.GetLoginUrl() != nil {
		model.LoginUrl = types.StringValue(*servicePrincipal.GetLoginUrl())
	} else {
		model.LoginUrl = types.StringNull()
	}

	// Map LogoutUrl
	if servicePrincipal.GetLogoutUrl() != nil {
		model.LogoutUrl = types.StringValue(*servicePrincipal.GetLogoutUrl())
	} else {
		model.LogoutUrl = types.StringNull()
	}

	// Map Notes
	if servicePrincipal.GetNotes() != nil {
		model.Notes = types.StringValue(*servicePrincipal.GetNotes())
	} else {
		model.Notes = types.StringNull()
	}

	// Map notification email addresses
	if emailAddresses := servicePrincipal.GetNotificationEmailAddresses(); len(emailAddresses) > 0 {
		elements := make([]attr.Value, len(emailAddresses))
		for i, email := range emailAddresses {
			elements[i] = types.StringValue(email)
		}
		model.NotificationEmailAddresses, _ = types.SetValue(types.StringType, elements)
	} else {
		model.NotificationEmailAddresses = types.SetNull(types.StringType)
	}

	// Map SamlMetadataUrl
	if servicePrincipal.GetSamlMetadataUrl() != nil {
		model.SamlMetadataUrl = types.StringValue(*servicePrincipal.GetSamlMetadataUrl())
	} else {
		model.SamlMetadataUrl = types.StringNull()
	}

	// Map PreferredTokenSigningKeyEndDateTime
	if servicePrincipal.GetPreferredTokenSigningKeyEndDateTime() != nil {
		model.PreferredTokenSigningKeyEndDateTime = types.StringValue(servicePrincipal.GetPreferredTokenSigningKeyEndDateTime().Format("2006-01-02T15:04:05Z"))
	} else {
		model.PreferredTokenSigningKeyEndDateTime = types.StringNull()
	}

	// Map PreferredTokenSigningKeyThumbprint
	if servicePrincipal.GetPreferredTokenSigningKeyThumbprint() != nil {
		model.PreferredTokenSigningKeyThumbprint = types.StringValue(*servicePrincipal.GetPreferredTokenSigningKeyThumbprint())
	} else {
		model.PreferredTokenSigningKeyThumbprint = types.StringNull()
	}

	// Map SAML single sign-on settings
	if samlSettings := servicePrincipal.GetSamlSingleSignOnSettings(); samlSettings != nil {
		samlAttrs := map[string]attr.Value{
			"relay_state": types.StringNull(),
		}
		if relayState := samlSettings.GetRelayState(); relayState != nil {
			samlAttrs["relay_state"] = types.StringValue(*relayState)
		}
		model.SamlSingleSignOnSettings, _ = types.ObjectValue(samlSingleSignOnSettingsAttrTypes, samlAttrs)
	} else {
		model.SamlSingleSignOnSettings = types.ObjectNull(samlSingleSignOnSettingsAttrTypes)
	}

	// Map verified publisher
	if verifiedPublisher := servicePrincipal.GetVerifiedPublisher(); verifiedPublisher != nil {
		vpAttrs := map[string]attr.Value{
			"display_name":          types.StringNull(),
			"verified_publisher_id": types.StringNull(),
			"added_date_time":       types.StringNull(),
		}
		if displayName := verifiedPublisher.GetDisplayName(); displayName != nil {
			vpAttrs["display_name"] = types.StringValue(*displayName)
		}
		if verifiedPublisherId := verifiedPublisher.GetVerifiedPublisherId(); verifiedPublisherId != nil {
			vpAttrs["verified_publisher_id"] = types.StringValue(*verifiedPublisherId)
		}
		if addedDateTime := verifiedPublisher.GetAddedDateTime(); addedDateTime != nil {
			vpAttrs["added_date_time"] = types.StringValue(addedDateTime.Format("2006-01-02T15:04:05Z"))
		}
		model.VerifiedPublisher, _ = types.ObjectValue(verifiedPublisherAttrTypes, vpAttrs)
	} else {
		model.VerifiedPublisher = types.ObjectNull(verifiedPublisherAttrTypes)
	}

	// Map informational URLs
	if info := servicePrincipal.GetInfo(); info != nil {
		infoAttrs := map[string]attr.Value{
			"terms_of_service_url":  types.StringNull(),
			"support_url":           types.StringNull(),
			"privacy_statement_url": types.StringNull(),
			"marketing_url":         types.StringNull(),
			"logo_url":              types.StringNull(),
		}
		if termsOfServiceUrl := info.GetTermsOfServiceUrl(); termsOfServiceUrl != nil {
			infoAttrs["terms_of_service_url"] = types.StringValue(*termsOfServiceUrl)
		}
		if supportUrl := info.GetSupportUrl(); supportUrl != nil {
			infoAttrs["support_url"] = types.StringValue(*supportUrl)
		}
		if privacyStatementUrl := info.GetPrivacyStatementUrl(); privacyStatementUrl != nil {
			infoAttrs["privacy_statement_url"] = types.StringValue(*privacyStatementUrl)
		}
		if marketingUrl := info.GetMarketingUrl(); marketingUrl != nil {
			infoAttrs["marketing_url"] = types.StringValue(*marketingUrl)
		}
		if logoUrl := info.GetLogoUrl(); logoUrl != nil {
			infoAttrs["logo_url"] = types.StringValue(*logoUrl)
		}
		model.Info, _ = types.ObjectValue(infoAttrTypes, infoAttrs)
	} else {
		model.Info = types.ObjectNull(infoAttrTypes)
	}

	tflog.Debug(ctx, "Successfully mapped service principal to data source model", map[string]any{
		"id":          model.ID.ValueString(),
		"displayName": model.DisplayName.ValueString(),
		"appId":       model.AppId.ValueString(),
	})

	return model
}
