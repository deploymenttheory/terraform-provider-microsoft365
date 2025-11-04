package graphBetaServicePrincipal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a service principal from the Graph API to the data source model
func MapRemoteStateToDataSource(ctx context.Context, servicePrincipal graphmodels.ServicePrincipalable) ServicePrincipalModel {
	model := ServicePrincipalModel{}

	if servicePrincipal.GetId() != nil {
		model.ID = types.StringValue(*servicePrincipal.GetId())
	} else {
		model.ID = types.StringNull()
	}

	if servicePrincipal.GetAppId() != nil {
		model.AppID = types.StringValue(*servicePrincipal.GetAppId())
	} else {
		model.AppID = types.StringNull()
	}

	if servicePrincipal.GetAppDisplayName() != nil {
		model.AppDisplayName = types.StringValue(*servicePrincipal.GetAppDisplayName())
	} else {
		model.AppDisplayName = types.StringNull()
	}

	if servicePrincipal.GetDisplayName() != nil {
		model.DisplayName = types.StringValue(*servicePrincipal.GetDisplayName())
	} else {
		model.DisplayName = types.StringNull()
	}

	if servicePrincipal.GetDeletedDateTime() != nil {
		model.DeletedDateTime = types.StringValue(servicePrincipal.GetDeletedDateTime().Format("2006-01-02T15:04:05Z"))
	} else {
		model.DeletedDateTime = types.StringNull()
	}

	if servicePrincipal.GetApplicationTemplateId() != nil {
		model.ApplicationTemplateID = types.StringValue(*servicePrincipal.GetApplicationTemplateId())
	} else {
		model.ApplicationTemplateID = types.StringNull()
	}

	if servicePrincipal.GetAccountEnabled() != nil {
		model.AccountEnabled = types.BoolValue(*servicePrincipal.GetAccountEnabled())
	} else {
		model.AccountEnabled = types.BoolNull()
	}

	if servicePrincipal.GetAppRoleAssignmentRequired() != nil {
		model.AppRoleAssignmentRequired = types.BoolValue(*servicePrincipal.GetAppRoleAssignmentRequired())
	} else {
		model.AppRoleAssignmentRequired = types.BoolNull()
	}

	if servicePrincipal.GetServicePrincipalType() != nil {
		model.ServicePrincipalType = types.StringValue(*servicePrincipal.GetServicePrincipalType())
	} else {
		model.ServicePrincipalType = types.StringNull()
	}

	if servicePrincipal.GetSignInAudience() != nil {
		model.SignInAudience = types.StringValue(*servicePrincipal.GetSignInAudience())
	} else {
		model.SignInAudience = types.StringNull()
	}

	if servicePrincipal.GetPreferredSingleSignOnMode() != nil {
		model.PreferredSingleSignOnMode = types.StringValue(*servicePrincipal.GetPreferredSingleSignOnMode())
	} else {
		model.PreferredSingleSignOnMode = types.StringNull()
	}

	if servicePrincipal.GetHomepage() != nil {
		model.Homepage = types.StringValue(*servicePrincipal.GetHomepage())
	} else {
		model.Homepage = types.StringNull()
	}

	if servicePrincipal.GetErrorUrl() != nil {
		model.ErrorUrl = types.StringValue(*servicePrincipal.GetErrorUrl())
	} else {
		model.ErrorUrl = types.StringNull()
	}

	if servicePrincipal.GetPublisherName() != nil {
		model.PublisherName = types.StringValue(*servicePrincipal.GetPublisherName())
	} else {
		model.PublisherName = types.StringNull()
	}

	// Map reply URLs
	if replyUrls := servicePrincipal.GetReplyUrls(); replyUrls != nil {
		model.ReplyUrls = make([]types.String, len(replyUrls))
		for i, url := range replyUrls {
			model.ReplyUrls[i] = types.StringValue(url)
		}
	} else {
		model.ReplyUrls = []types.String{}
	}

	// Map service principal names
	if spNames := servicePrincipal.GetServicePrincipalNames(); spNames != nil {
		model.ServicePrincipalNames = make([]types.String, len(spNames))
		for i, name := range spNames {
			model.ServicePrincipalNames[i] = types.StringValue(name)
		}
	} else {
		model.ServicePrincipalNames = []types.String{}
	}

	// Map tags
	if tags := servicePrincipal.GetTags(); tags != nil {
		model.Tags = make([]types.String, len(tags))
		for i, tag := range tags {
			model.Tags[i] = types.StringValue(tag)
		}
	} else {
		model.Tags = []types.String{}
	}

	if servicePrincipal.GetDisabledByMicrosoftStatus() != nil {
		model.DisabledByMicrosoftStatus = types.StringValue(*servicePrincipal.GetDisabledByMicrosoftStatus())
	} else {
		model.DisabledByMicrosoftStatus = types.StringNull()
	}

	if servicePrincipal.GetAppOwnerOrganizationId() != nil {
		model.AppOwnerOrganizationID = types.StringValue(servicePrincipal.GetAppOwnerOrganizationId().String())
	} else {
		model.AppOwnerOrganizationID = types.StringNull()
	}

	if servicePrincipal.GetLoginUrl() != nil {
		model.LoginUrl = types.StringValue(*servicePrincipal.GetLoginUrl())
	} else {
		model.LoginUrl = types.StringNull()
	}

	if servicePrincipal.GetLogoutUrl() != nil {
		model.LogoutUrl = types.StringValue(*servicePrincipal.GetLogoutUrl())
	} else {
		model.LogoutUrl = types.StringNull()
	}

	if servicePrincipal.GetNotes() != nil {
		model.Notes = types.StringValue(*servicePrincipal.GetNotes())
	} else {
		model.Notes = types.StringNull()
	}

	// Map notification email addresses
	if emailAddresses := servicePrincipal.GetNotificationEmailAddresses(); emailAddresses != nil {
		model.NotificationEmailAddresses = make([]types.String, len(emailAddresses))
		for i, email := range emailAddresses {
			model.NotificationEmailAddresses[i] = types.StringValue(email)
		}
	} else {
		model.NotificationEmailAddresses = []types.String{}
	}

	if servicePrincipal.GetSamlMetadataUrl() != nil {
		model.SamlMetadataUrl = types.StringValue(*servicePrincipal.GetSamlMetadataUrl())
	} else {
		model.SamlMetadataUrl = types.StringNull()
	}

	if servicePrincipal.GetPreferredTokenSigningKeyEndDateTime() != nil {
		model.PreferredTokenSigningKeyEndDateTime = types.StringValue(servicePrincipal.GetPreferredTokenSigningKeyEndDateTime().Format("2006-01-02T15:04:05Z"))
	} else {
		model.PreferredTokenSigningKeyEndDateTime = types.StringNull()
	}

	if servicePrincipal.GetPreferredTokenSigningKeyThumbprint() != nil {
		model.PreferredTokenSigningKeyThumbprint = types.StringValue(*servicePrincipal.GetPreferredTokenSigningKeyThumbprint())
	} else {
		model.PreferredTokenSigningKeyThumbprint = types.StringNull()
	}

	// Map SAML single sign-on settings
	if samlSettings := servicePrincipal.GetSamlSingleSignOnSettings(); samlSettings != nil {
		model.SamlSingleSignOnSettings = &SamlSingleSignOnSettingsModel{}
		if relayState := samlSettings.GetRelayState(); relayState != nil {
			model.SamlSingleSignOnSettings.RelayState = types.StringValue(*relayState)
		} else {
			model.SamlSingleSignOnSettings.RelayState = types.StringNull()
		}
	} else {
		model.SamlSingleSignOnSettings = nil
	}

	// Map verified publisher
	if verifiedPublisher := servicePrincipal.GetVerifiedPublisher(); verifiedPublisher != nil {
		model.VerifiedPublisher = &VerifiedPublisherModel{}
		if displayName := verifiedPublisher.GetDisplayName(); displayName != nil {
			model.VerifiedPublisher.DisplayName = types.StringValue(*displayName)
		} else {
			model.VerifiedPublisher.DisplayName = types.StringNull()
		}
		if verifiedPublisherId := verifiedPublisher.GetVerifiedPublisherId(); verifiedPublisherId != nil {
			model.VerifiedPublisher.VerifiedPublisherID = types.StringValue(*verifiedPublisherId)
		} else {
			model.VerifiedPublisher.VerifiedPublisherID = types.StringNull()
		}
		if addedDateTime := verifiedPublisher.GetAddedDateTime(); addedDateTime != nil {
			model.VerifiedPublisher.AddedDateTime = types.StringValue(addedDateTime.Format("2006-01-02T15:04:05Z"))
		} else {
			model.VerifiedPublisher.AddedDateTime = types.StringNull()
		}
	} else {
		model.VerifiedPublisher = nil
	}

	// Map informational URLs
	if info := servicePrincipal.GetInfo(); info != nil {
		model.Info = &InformationalUrlModel{}
		if termsOfServiceUrl := info.GetTermsOfServiceUrl(); termsOfServiceUrl != nil {
			model.Info.TermsOfServiceUrl = types.StringValue(*termsOfServiceUrl)
		} else {
			model.Info.TermsOfServiceUrl = types.StringNull()
		}
		if supportUrl := info.GetSupportUrl(); supportUrl != nil {
			model.Info.SupportUrl = types.StringValue(*supportUrl)
		} else {
			model.Info.SupportUrl = types.StringNull()
		}
		if privacyStatementUrl := info.GetPrivacyStatementUrl(); privacyStatementUrl != nil {
			model.Info.PrivacyStatementUrl = types.StringValue(*privacyStatementUrl)
		} else {
			model.Info.PrivacyStatementUrl = types.StringNull()
		}
		if marketingUrl := info.GetMarketingUrl(); marketingUrl != nil {
			model.Info.MarketingUrl = types.StringValue(*marketingUrl)
		} else {
			model.Info.MarketingUrl = types.StringNull()
		}
		if logoUrl := info.GetLogoUrl(); logoUrl != nil {
			model.Info.LogoUrl = types.StringValue(*logoUrl)
		} else {
			model.Info.LogoUrl = types.StringNull()
		}
	} else {
		model.Info = nil
	}

	tflog.Debug(ctx, "Successfully mapped service principal to data source model", map[string]any{
		"id":          model.ID.ValueString(),
		"displayName": model.DisplayName.ValueString(),
		"appId":       model.AppID.ValueString(),
	})

	return model
}
