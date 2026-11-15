package graphBetaIdentityAndAccessAuthenticationMethodsPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the authentication methods policy and its method configurations to the Terraform state.
//
// Every managed block is populated from the tenant, including blocks that are not configured, so import and refresh
// produce a complete state. A method configuration that the tenant does not return is set to null.
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data AuthenticationMethodsPolicyResourceModel,
	remoteResource graphmodels.AuthenticationMethodsPolicyable,
) AuthenticationMethodsPolicyResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = types.StringValue(singletonID)
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.PolicyMigrationState = convert.GraphToFrameworkEnum(remoteResource.GetPolicyMigrationState())
	data.RegistrationEnforcement = mapRegistrationEnforcement(remoteResource.GetRegistrationEnforcement())
	data.ReportSuspiciousActivitySettings = mapReportSuspiciousActivitySettings(remoteResource.GetReportSuspiciousActivitySettings())
	data.SystemCredentialPreferences = mapSystemCredentialPreferences(remoteResource.GetSystemCredentialPreferences())

	data.Email = types.ObjectNull(emailMethodAttrTypes)
	data.TemporaryAccessPass = types.ObjectNull(temporaryAccessPassMethodAttrTypes)
	data.Sms = types.ObjectNull(smsMethodAttrTypes)
	data.SoftwareOath = types.ObjectNull(oathMethodAttrTypes)
	data.HardwareOath = types.ObjectNull(oathMethodAttrTypes)

	// Match on the fixed configuration ID: several method interfaces share the same method set, so a type
	// switch on the interfaces alone cannot tell them apart.
	for _, configuration := range remoteResource.GetAuthenticationMethodConfigurations() {
		if configuration == nil || configuration.GetId() == nil {
			continue
		}
		switch *configuration.GetId() {
		case methodIDEmail:
			if method, ok := configuration.(*graphmodels.EmailAuthenticationMethodConfiguration); ok {
				data.Email = types.ObjectValueMust(emailMethodAttrTypes, map[string]attr.Value{
					"state":                              convert.GraphToFrameworkEnum(method.GetState()),
					"allow_external_id_to_use_email_otp": convert.GraphToFrameworkEnum(method.GetAllowExternalIdToUseEmailOtp()),
					"include_targets":                    mapMethodTargets(method.GetIncludeTargets()),
					"exclude_targets":                    mapExcludeTargets(method.GetExcludeTargets()),
				})
			}
		case methodIDTemporaryAccessPass:
			if method, ok := configuration.(*graphmodels.TemporaryAccessPassAuthenticationMethodConfiguration); ok {
				data.TemporaryAccessPass = types.ObjectValueMust(temporaryAccessPassMethodAttrTypes, map[string]attr.Value{
					"state":                       convert.GraphToFrameworkEnum(method.GetState()),
					"default_lifetime_in_minutes": convert.GraphToFrameworkInt32(method.GetDefaultLifetimeInMinutes()),
					"default_length":              convert.GraphToFrameworkInt32(method.GetDefaultLength()),
					"minimum_lifetime_in_minutes": convert.GraphToFrameworkInt32(method.GetMinimumLifetimeInMinutes()),
					"maximum_lifetime_in_minutes": convert.GraphToFrameworkInt32(method.GetMaximumLifetimeInMinutes()),
					"is_usable_once":              convert.GraphToFrameworkBool(method.GetIsUsableOnce()),
					"include_targets":             mapMethodTargets(method.GetIncludeTargets()),
					"exclude_targets":             mapExcludeTargets(method.GetExcludeTargets()),
				})
			}
		case methodIDSms:
			if method, ok := configuration.(*graphmodels.SmsAuthenticationMethodConfiguration); ok {
				data.Sms = types.ObjectValueMust(smsMethodAttrTypes, map[string]attr.Value{
					"state":           convert.GraphToFrameworkEnum(method.GetState()),
					"include_targets": mapSmsMethodTargets(method.GetIncludeTargets()),
					"exclude_targets": mapExcludeTargets(method.GetExcludeTargets()),
				})
			}
		case methodIDSoftwareOath:
			if method, ok := configuration.(*graphmodels.SoftwareOathAuthenticationMethodConfiguration); ok {
				data.SoftwareOath = mapOathMethod(method.GetState(), method.GetIncludeTargets(), method.GetExcludeTargets())
			}
		case methodIDHardwareOath:
			if method, ok := configuration.(*graphmodels.HardwareOathAuthenticationMethodConfiguration); ok {
				data.HardwareOath = mapOathMethod(method.GetState(), method.GetIncludeTargets(), method.GetExcludeTargets())
			}
		}
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"resourceId": data.ID.ValueString(),
	})

	return data
}

// mapOathMethod maps a softwareOath or hardwareOath method configuration.
func mapOathMethod(
	state *graphmodels.AuthenticationMethodState,
	includeTargets []graphmodels.AuthenticationMethodTargetable,
	excludeTargets []graphmodels.ExcludeTargetable,
) types.Object {
	return types.ObjectValueMust(oathMethodAttrTypes, map[string]attr.Value{
		"state":           convert.GraphToFrameworkEnum(state),
		"include_targets": mapMethodTargets(includeTargets),
		"exclude_targets": mapExcludeTargets(excludeTargets),
	})
}

// mapRegistrationEnforcement maps the registrationEnforcement complex type.
func mapRegistrationEnforcement(enforcement graphmodels.RegistrationEnforcementable) types.Object {
	if enforcement == nil {
		return types.ObjectNull(registrationEnforcementAttrTypes)
	}

	campaignValue := types.ObjectNull(registrationCampaignAttrTypes)
	if campaign := enforcement.GetAuthenticationMethodsRegistrationCampaign(); campaign != nil {
		includeElements := make([]attr.Value, 0, len(campaign.GetIncludeTargets()))
		for _, target := range campaign.GetIncludeTargets() {
			includeElements = append(includeElements, types.ObjectValueMust(registrationCampaignIncludeTargetAttrTypes, map[string]attr.Value{
				"id":                             convert.GraphToFrameworkString(target.GetId()),
				"target_type":                    convert.GraphToFrameworkEnum(target.GetTargetType()),
				"targeted_authentication_method": convert.GraphToFrameworkString(target.GetTargetedAuthenticationMethod()),
			}))
		}
		campaignValue = types.ObjectValueMust(registrationCampaignAttrTypes, map[string]attr.Value{
			"snooze_duration_in_days":                    convert.GraphToFrameworkInt32(campaign.GetSnoozeDurationInDays()),
			"enforce_registration_after_allowed_snoozes": convert.GraphToFrameworkBool(campaign.GetEnforceRegistrationAfterAllowedSnoozes()),
			"state":           convert.GraphToFrameworkEnum(campaign.GetState()),
			"include_targets": types.SetValueMust(types.ObjectType{AttrTypes: registrationCampaignIncludeTargetAttrTypes}, includeElements),
			"exclude_targets": mapExcludeTargets(campaign.GetExcludeTargets()),
		})
	}

	return types.ObjectValueMust(registrationEnforcementAttrTypes, map[string]attr.Value{
		"authentication_methods_registration_campaign": campaignValue,
	})
}

// mapReportSuspiciousActivitySettings maps the reportSuspiciousActivitySettings complex type.
func mapReportSuspiciousActivitySettings(settings graphmodels.ReportSuspiciousActivitySettingsable) types.Object {
	if settings == nil {
		return types.ObjectNull(reportSuspiciousActivitySettingsAttrTypes)
	}

	includeTarget := types.ObjectNull(targetAttrTypes)
	if target := settings.GetIncludeTarget(); target != nil {
		includeTarget = types.ObjectValueMust(targetAttrTypes, map[string]attr.Value{
			"id":          convert.GraphToFrameworkString(target.GetId()),
			"target_type": convert.GraphToFrameworkEnum(target.GetTargetType()),
		})
	}

	return types.ObjectValueMust(reportSuspiciousActivitySettingsAttrTypes, map[string]attr.Value{
		"state":                convert.GraphToFrameworkEnum(settings.GetState()),
		"voice_reporting_code": convert.GraphToFrameworkInt32(settings.GetVoiceReportingCode()),
		"include_target":       includeTarget,
	})
}

// mapSystemCredentialPreferences maps the systemCredentialPreferences complex type.
func mapSystemCredentialPreferences(preferences graphmodels.SystemCredentialPreferencesable) types.Object {
	if preferences == nil {
		return types.ObjectNull(systemCredentialPreferencesAttrTypes)
	}

	includeElements := make([]attr.Value, 0, len(preferences.GetIncludeTargets()))
	for _, target := range preferences.GetIncludeTargets() {
		includeElements = append(includeElements, types.ObjectValueMust(targetAttrTypes, map[string]attr.Value{
			"id":          convert.GraphToFrameworkString(target.GetId()),
			"target_type": convert.GraphToFrameworkEnum(target.GetTargetType()),
		}))
	}

	return types.ObjectValueMust(systemCredentialPreferencesAttrTypes, map[string]attr.Value{
		"state":           convert.GraphToFrameworkEnum(preferences.GetState()),
		"include_targets": types.SetValueMust(types.ObjectType{AttrTypes: targetAttrTypes}, includeElements),
		"exclude_targets": mapExcludeTargets(preferences.GetExcludeTargets()),
	})
}

// mapMethodTargets maps authenticationMethodTarget values. An empty or missing collection maps to an empty set.
func mapMethodTargets(targets []graphmodels.AuthenticationMethodTargetable) types.Set {
	elements := make([]attr.Value, 0, len(targets))
	for _, target := range targets {
		elements = append(elements, types.ObjectValueMust(methodTargetAttrTypes, map[string]attr.Value{
			"id":                       convert.GraphToFrameworkString(target.GetId()),
			"target_type":              convert.GraphToFrameworkEnum(target.GetTargetType()),
			"is_registration_required": convert.GraphToFrameworkBoolWithDefault(target.GetIsRegistrationRequired(), false),
		}))
	}
	return types.SetValueMust(types.ObjectType{AttrTypes: methodTargetAttrTypes}, elements)
}

// mapSmsMethodTargets maps smsAuthenticationMethodTarget values. An empty or missing collection maps to an empty set.
func mapSmsMethodTargets(targets []graphmodels.SmsAuthenticationMethodTargetable) types.Set {
	elements := make([]attr.Value, 0, len(targets))
	for _, target := range targets {
		elements = append(elements, types.ObjectValueMust(smsMethodTargetAttrTypes, map[string]attr.Value{
			"id":                       convert.GraphToFrameworkString(target.GetId()),
			"target_type":              convert.GraphToFrameworkEnum(target.GetTargetType()),
			"is_registration_required": convert.GraphToFrameworkBoolWithDefault(target.GetIsRegistrationRequired(), false),
			"is_usable_for_sign_in":    convert.GraphToFrameworkBoolWithDefault(target.GetIsUsableForSignIn(), false),
		}))
	}
	return types.SetValueMust(types.ObjectType{AttrTypes: smsMethodTargetAttrTypes}, elements)
}

// mapExcludeTargets maps excludeTarget values. An empty or missing collection maps to an empty set.
func mapExcludeTargets(targets []graphmodels.ExcludeTargetable) types.Set {
	elements := make([]attr.Value, 0, len(targets))
	for _, target := range targets {
		elements = append(elements, types.ObjectValueMust(targetAttrTypes, map[string]attr.Value{
			"id":          convert.GraphToFrameworkString(target.GetId()),
			"target_type": convert.GraphToFrameworkEnum(target.GetTargetType()),
		}))
	}
	return types.SetValueMust(types.ObjectType{AttrTypes: targetAttrTypes}, elements)
}
