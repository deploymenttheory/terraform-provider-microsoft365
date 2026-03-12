package graphBetaCrossTenantAccessDefaultSettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote crossTenantAccessPolicyConfigurationDefault API response to Terraform state.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *CrossTenantAccessDefaultSettingsResourceModel, remoteResource graphmodels.CrossTenantAccessPolicyConfigurationDefaultable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state")

	// This is a singleton — the ID is always the static identifier used throughout the provider.
	data.ID = types.StringValue(singletonID)

	// Map the is_service_default property
	if isServiceDefault := remoteResource.GetIsServiceDefault(); isServiceDefault != nil {
		data.IsServiceDefault = types.BoolValue(*isServiceDefault)
	} else {
		data.IsServiceDefault = types.BoolNull()
	}

	// Map B2bCollaborationInbound
	if b2bCollabInbound := remoteResource.GetB2bCollaborationInbound(); b2bCollabInbound != nil {
		data.B2bCollaborationInbound = mapB2BSettingToTerraform(ctx, b2bCollabInbound)
	} else {
		data.B2bCollaborationInbound = types.ObjectNull(getB2BSettingAttrTypes())
	}

	// Map B2bCollaborationOutbound
	if b2bCollabOutbound := remoteResource.GetB2bCollaborationOutbound(); b2bCollabOutbound != nil {
		data.B2bCollaborationOutbound = mapB2BSettingToTerraform(ctx, b2bCollabOutbound)
	} else {
		data.B2bCollaborationOutbound = types.ObjectNull(getB2BSettingAttrTypes())
	}

	// Map B2bDirectConnectInbound
	if b2bDirectInbound := remoteResource.GetB2bDirectConnectInbound(); b2bDirectInbound != nil {
		data.B2bDirectConnectInbound = mapB2BSettingToTerraform(ctx, b2bDirectInbound)
	} else {
		data.B2bDirectConnectInbound = types.ObjectNull(getB2BSettingAttrTypes())
	}

	// Map B2bDirectConnectOutbound
	if b2bDirectOutbound := remoteResource.GetB2bDirectConnectOutbound(); b2bDirectOutbound != nil {
		data.B2bDirectConnectOutbound = mapB2BSettingToTerraform(ctx, b2bDirectOutbound)
	} else {
		data.B2bDirectConnectOutbound = types.ObjectNull(getB2BSettingAttrTypes())
	}

	// Map InboundTrust
	if inboundTrust := remoteResource.GetInboundTrust(); inboundTrust != nil {
		data.InboundTrust = mapInboundTrustToTerraform(ctx, inboundTrust)
	} else {
		data.InboundTrust = types.ObjectNull(getInboundTrustAttrTypes())
	}

	// Map InvitationRedemptionIdentityProviderConfiguration
	if invitationConfig := remoteResource.GetInvitationRedemptionIdentityProviderConfiguration(); invitationConfig != nil {
		data.InvitationRedemptionIdentityProviderConfiguration = mapInvitationRedemptionConfigToTerraform(ctx, invitationConfig)
	} else {
		data.InvitationRedemptionIdentityProviderConfiguration = types.ObjectNull(getInvitationRedemptionConfigAttrTypes())
	}

	// Map TenantRestrictions
	if tenantRestrictions := remoteResource.GetTenantRestrictions(); tenantRestrictions != nil {
		data.TenantRestrictions = mapTenantRestrictionsToTerraform(ctx, tenantRestrictions)
	} else {
		data.TenantRestrictions = types.ObjectNull(getTenantRestrictionsAttrTypes())
	}

	// Map AutomaticUserConsentSettings
	if autoConsent := remoteResource.GetAutomaticUserConsentSettings(); autoConsent != nil {
		data.AutomaticUserConsentSettings = mapAutomaticUserConsentToTerraform(ctx, autoConsent)
	} else {
		data.AutomaticUserConsentSettings = types.ObjectNull(getAutomaticUserConsentAttrTypes())
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// mapB2BSettingToTerraform maps a CrossTenantAccessPolicyB2BSetting to types.Object
func mapB2BSettingToTerraform(ctx context.Context, b2bSetting graphmodels.CrossTenantAccessPolicyB2BSettingable) types.Object {
	if b2bSetting == nil {
		return types.ObjectNull(getB2BSettingAttrTypes())
	}

	attrs := map[string]attr.Value{
		"users_and_groups": types.ObjectNull(getTargetConfigurationAttrTypes()),
		"applications":     types.ObjectNull(getTargetConfigurationAttrTypes()),
	}

	if usersAndGroups := b2bSetting.GetUsersAndGroups(); usersAndGroups != nil {
		attrs["users_and_groups"] = mapTargetConfigurationToTerraform(ctx, usersAndGroups)
	}

	if applications := b2bSetting.GetApplications(); applications != nil {
		attrs["applications"] = mapTargetConfigurationToTerraform(ctx, applications)
	}

	obj, diags := types.ObjectValue(getB2BSettingAttrTypes(), attrs)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create B2B setting object", map[string]any{"diagnostics": diags})
		return types.ObjectNull(getB2BSettingAttrTypes())
	}

	return obj
}

// mapTargetConfigurationToTerraform maps a CrossTenantAccessPolicyTargetConfiguration to types.Object
func mapTargetConfigurationToTerraform(ctx context.Context, targetConfig graphmodels.CrossTenantAccessPolicyTargetConfigurationable) types.Object {
	if targetConfig == nil {
		return types.ObjectNull(getTargetConfigurationAttrTypes())
	}

	attrs := map[string]attr.Value{
		"access_type": types.StringNull(),
		"targets":     types.SetNull(types.ObjectType{AttrTypes: getTargetAttrTypes()}),
	}

	// Map AccessType
	if accessType := targetConfig.GetAccessType(); accessType != nil {
		switch *accessType {
		case graphmodels.ALLOWED_CROSSTENANTACCESSPOLICYTARGETCONFIGURATIONACCESSTYPE:
			attrs["access_type"] = types.StringValue("allowed")
		case graphmodels.BLOCKED_CROSSTENANTACCESSPOLICYTARGETCONFIGURATIONACCESSTYPE:
			attrs["access_type"] = types.StringValue("blocked")
		}
	}

	// Map Targets
	if targets := targetConfig.GetTargets(); len(targets) > 0 {
		targetElements := make([]attr.Value, 0, len(targets))
		for _, target := range targets {
			targetAttrs := map[string]attr.Value{
				"target":      types.StringNull(),
				"target_type": types.StringNull(),
			}

			if targetVal := target.GetTarget(); targetVal != nil {
				targetAttrs["target"] = types.StringValue(*targetVal)
			}

			if targetType := target.GetTargetType(); targetType != nil {
				switch *targetType {
				case graphmodels.USER_CROSSTENANTACCESSPOLICYTARGETTYPE:
					targetAttrs["target_type"] = types.StringValue("user")
				case graphmodels.GROUP_CROSSTENANTACCESSPOLICYTARGETTYPE:
					targetAttrs["target_type"] = types.StringValue("group")
				case graphmodels.APPLICATION_CROSSTENANTACCESSPOLICYTARGETTYPE:
					targetAttrs["target_type"] = types.StringValue("application")
				}
			}

			targetObj, diags := types.ObjectValue(getTargetAttrTypes(), targetAttrs)
			if !diags.HasError() {
				targetElements = append(targetElements, targetObj)
			}
		}

		if len(targetElements) > 0 {
			targetSet, diags := types.SetValue(types.ObjectType{AttrTypes: getTargetAttrTypes()}, targetElements)
			if !diags.HasError() {
				attrs["targets"] = targetSet
			}
		}
	}

	obj, diags := types.ObjectValue(getTargetConfigurationAttrTypes(), attrs)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create target configuration object", map[string]any{"diagnostics": diags})
		return types.ObjectNull(getTargetConfigurationAttrTypes())
	}

	return obj
}

// mapInboundTrustToTerraform maps a CrossTenantAccessPolicyInboundTrust to types.Object
func mapInboundTrustToTerraform(ctx context.Context, inboundTrust graphmodels.CrossTenantAccessPolicyInboundTrustable) types.Object {
	if inboundTrust == nil {
		return types.ObjectNull(getInboundTrustAttrTypes())
	}

	attrs := map[string]attr.Value{
		"is_mfa_accepted":                         types.BoolNull(),
		"is_compliant_device_accepted":            types.BoolNull(),
		"is_hybrid_azure_ad_joined_device_accepted": types.BoolNull(),
	}

	if isMfaAccepted := inboundTrust.GetIsMfaAccepted(); isMfaAccepted != nil {
		attrs["is_mfa_accepted"] = types.BoolValue(*isMfaAccepted)
	}

	if isCompliantDeviceAccepted := inboundTrust.GetIsCompliantDeviceAccepted(); isCompliantDeviceAccepted != nil {
		attrs["is_compliant_device_accepted"] = types.BoolValue(*isCompliantDeviceAccepted)
	}

	if isHybridAzureADJoinedDeviceAccepted := inboundTrust.GetIsHybridAzureADJoinedDeviceAccepted(); isHybridAzureADJoinedDeviceAccepted != nil {
		attrs["is_hybrid_azure_ad_joined_device_accepted"] = types.BoolValue(*isHybridAzureADJoinedDeviceAccepted)
	}

	obj, diags := types.ObjectValue(getInboundTrustAttrTypes(), attrs)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create inbound trust object", map[string]any{"diagnostics": diags})
		return types.ObjectNull(getInboundTrustAttrTypes())
	}

	return obj
}

// mapInvitationRedemptionConfigToTerraform maps a DefaultInvitationRedemptionIdentityProviderConfiguration to types.Object
func mapInvitationRedemptionConfigToTerraform(ctx context.Context, config graphmodels.DefaultInvitationRedemptionIdentityProviderConfigurationable) types.Object {
	if config == nil {
		return types.ObjectNull(getInvitationRedemptionConfigAttrTypes())
	}

	attrs := map[string]attr.Value{
		"primary_identity_provider_precedence_order": types.ListNull(types.StringType),
		"fallback_identity_provider":                 types.StringNull(),
	}

	// Map PrimaryIdentityProviderPrecedenceOrder
	if precedenceOrder := config.GetPrimaryIdentityProviderPrecedenceOrder(); len(precedenceOrder) > 0 {
		precedenceElements := make([]attr.Value, 0, len(precedenceOrder))
		for _, provider := range precedenceOrder {
			precedenceElements = append(precedenceElements, types.StringValue(provider.String()))
		}
		if len(precedenceElements) > 0 {
			precedenceList, diags := types.ListValue(types.StringType, precedenceElements)
			if !diags.HasError() {
				attrs["primary_identity_provider_precedence_order"] = precedenceList
			}
		}
	}

	// Map FallbackIdentityProvider
	if fallback := config.GetFallbackIdentityProvider(); fallback != nil {
		attrs["fallback_identity_provider"] = types.StringValue(fallback.String())
	}

	obj, diags := types.ObjectValue(getInvitationRedemptionConfigAttrTypes(), attrs)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create invitation redemption config object", map[string]any{"diagnostics": diags})
		return types.ObjectNull(getInvitationRedemptionConfigAttrTypes())
	}

	return obj
}

// mapTenantRestrictionsToTerraform maps a CrossTenantAccessPolicyTenantRestrictions to types.Object
// Note: The devices property is not supported on the server side yet, so we don't map it.
func mapTenantRestrictionsToTerraform(ctx context.Context, tenantRestrictions graphmodels.CrossTenantAccessPolicyTenantRestrictionsable) types.Object {
	if tenantRestrictions == nil {
		return types.ObjectNull(getTenantRestrictionsAttrTypes())
	}

	attrs := map[string]attr.Value{
		"users_and_groups": types.ObjectNull(getTargetConfigurationAttrTypes()),
		"applications":     types.ObjectNull(getTargetConfigurationAttrTypes()),
	}

	if usersAndGroups := tenantRestrictions.GetUsersAndGroups(); usersAndGroups != nil {
		attrs["users_and_groups"] = mapTargetConfigurationToTerraform(ctx, usersAndGroups)
	}

	if applications := tenantRestrictions.GetApplications(); applications != nil {
		attrs["applications"] = mapTargetConfigurationToTerraform(ctx, applications)
	}

	// Devices property is not supported on the server side yet according to SDK documentation
	// if devices := tenantRestrictions.GetDevices(); devices != nil {
	// 	attrs["devices"] = mapDevicesFilterToTerraform(ctx, devices)
	// }

	obj, diags := types.ObjectValue(getTenantRestrictionsAttrTypes(), attrs)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create tenant restrictions object", map[string]any{"diagnostics": diags})
		return types.ObjectNull(getTenantRestrictionsAttrTypes())
	}

	return obj
}

// mapAutomaticUserConsentToTerraform maps an InboundOutboundPolicyConfiguration to types.Object
func mapAutomaticUserConsentToTerraform(ctx context.Context, autoConsent graphmodels.InboundOutboundPolicyConfigurationable) types.Object {
	if autoConsent == nil {
		return types.ObjectNull(getAutomaticUserConsentAttrTypes())
	}

	attrs := map[string]attr.Value{
		"inbound_allowed":  types.BoolNull(),
		"outbound_allowed": types.BoolNull(),
	}

	if inboundAllowed := autoConsent.GetInboundAllowed(); inboundAllowed != nil {
		attrs["inbound_allowed"] = types.BoolValue(*inboundAllowed)
	}

	if outboundAllowed := autoConsent.GetOutboundAllowed(); outboundAllowed != nil {
		attrs["outbound_allowed"] = types.BoolValue(*outboundAllowed)
	}

	obj, diags := types.ObjectValue(getAutomaticUserConsentAttrTypes(), attrs)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create automatic user consent object", map[string]any{"diagnostics": diags})
		return types.ObjectNull(getAutomaticUserConsentAttrTypes())
	}

	return obj
}

// Helper functions to get attribute types for type definitions

func getB2BSettingAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"users_and_groups": types.ObjectType{AttrTypes: getTargetConfigurationAttrTypes()},
		"applications":     types.ObjectType{AttrTypes: getTargetConfigurationAttrTypes()},
	}
}

func getTargetConfigurationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"access_type": types.StringType,
		"targets":     types.SetType{ElemType: types.ObjectType{AttrTypes: getTargetAttrTypes()}},
	}
}

func getTargetAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target":      types.StringType,
		"target_type": types.StringType,
	}
}

func getInboundTrustAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"is_mfa_accepted":                         types.BoolType,
		"is_compliant_device_accepted":            types.BoolType,
		"is_hybrid_azure_ad_joined_device_accepted": types.BoolType,
	}
}

func getInvitationRedemptionConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"primary_identity_provider_precedence_order": types.ListType{ElemType: types.StringType},
		"fallback_identity_provider":                 types.StringType,
	}
}

func getTenantRestrictionsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"users_and_groups": types.ObjectType{AttrTypes: getTargetConfigurationAttrTypes()},
		"applications":     types.ObjectType{AttrTypes: getTargetConfigurationAttrTypes()},
	}
}

func getAutomaticUserConsentAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"inbound_allowed":  types.BoolType,
		"outbound_allowed": types.BoolType,
	}
}
