package graphBetaCrossTenantAccessPartnerSettings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *CrossTenantAccessPartnerSettingsResourceModel) (graphmodels.CrossTenantAccessPolicyConfigurationPartnerable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewCrossTenantAccessPolicyConfigurationPartner()

	tenantID := data.TenantID.ValueString()
	requestBody.SetTenantId(&tenantID)

	if !data.IsServiceProvider.IsNull() && !data.IsServiceProvider.IsUnknown() {
		val := data.IsServiceProvider.ValueBool()
		requestBody.SetIsServiceProvider(&val)
	}

	if !data.B2bCollaborationInbound.IsNull() && !data.B2bCollaborationInbound.IsUnknown() {
		b2bSetting, err := constructB2BSetting(ctx, data.B2bCollaborationInbound)
		if err != nil {
			return nil, fmt.Errorf("failed to construct b2b_collaboration_inbound: %w", err)
		}
		requestBody.SetB2bCollaborationInbound(b2bSetting)
	}

	if !data.B2bCollaborationOutbound.IsNull() && !data.B2bCollaborationOutbound.IsUnknown() {
		b2bSetting, err := constructB2BSetting(ctx, data.B2bCollaborationOutbound)
		if err != nil {
			return nil, fmt.Errorf("failed to construct b2b_collaboration_outbound: %w", err)
		}
		requestBody.SetB2bCollaborationOutbound(b2bSetting)
	}

	if !data.B2bDirectConnectInbound.IsNull() && !data.B2bDirectConnectInbound.IsUnknown() {
		b2bSetting, err := constructB2BSetting(ctx, data.B2bDirectConnectInbound)
		if err != nil {
			return nil, fmt.Errorf("failed to construct b2b_direct_connect_inbound: %w", err)
		}
		requestBody.SetB2bDirectConnectInbound(b2bSetting)
	}

	if !data.B2bDirectConnectOutbound.IsNull() && !data.B2bDirectConnectOutbound.IsUnknown() {
		b2bSetting, err := constructB2BSetting(ctx, data.B2bDirectConnectOutbound)
		if err != nil {
			return nil, fmt.Errorf("failed to construct b2b_direct_connect_outbound: %w", err)
		}
		requestBody.SetB2bDirectConnectOutbound(b2bSetting)
	}

	if !data.InboundTrust.IsNull() && !data.InboundTrust.IsUnknown() {
		inboundTrust, err := constructInboundTrust(ctx, data.InboundTrust)
		if err != nil {
			return nil, fmt.Errorf("failed to construct inbound_trust: %w", err)
		}
		requestBody.SetInboundTrust(inboundTrust)
	}

	if !data.AutomaticUserConsentSettings.IsNull() && !data.AutomaticUserConsentSettings.IsUnknown() {
		autoConsent, err := constructAutomaticUserConsentSettings(ctx, data.AutomaticUserConsentSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to construct automatic_user_consent_settings: %w", err)
		}
		requestBody.SetAutomaticUserConsentSettings(autoConsent)
	}

	if !data.TenantRestrictions.IsNull() && !data.TenantRestrictions.IsUnknown() {
		tenantRestr, err := constructTenantRestrictions(ctx, data.TenantRestrictions)
		if err != nil {
			return nil, fmt.Errorf("failed to construct tenant_restrictions: %w", err)
		}
		requestBody.SetTenantRestrictions(tenantRestr)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

func constructB2BSetting(ctx context.Context, obj types.Object) (graphmodels.CrossTenantAccessPolicyB2BSettingable, error) {
	if obj.IsNull() || obj.IsUnknown() {
		return nil, nil
	}

	b2bSetting := graphmodels.NewCrossTenantAccessPolicyB2BSetting()

	attrs := obj.Attributes()

	if usersAndGroupsObj, ok := attrs["users_and_groups"].(types.Object); ok && !usersAndGroupsObj.IsNull() {
		targetConfig, err := constructTargetConfiguration(ctx, usersAndGroupsObj)
		if err != nil {
			return nil, fmt.Errorf("failed to construct users_and_groups: %w", err)
		}
		b2bSetting.SetUsersAndGroups(targetConfig)
	}

	if applicationsObj, ok := attrs["applications"].(types.Object); ok && !applicationsObj.IsNull() {
		targetConfig, err := constructTargetConfiguration(ctx, applicationsObj)
		if err != nil {
			return nil, fmt.Errorf("failed to construct applications: %w", err)
		}
		b2bSetting.SetApplications(targetConfig)
	}

	return b2bSetting, nil
}

func constructTargetConfiguration(_ context.Context, obj types.Object) (graphmodels.CrossTenantAccessPolicyTargetConfigurationable, error) {
	if obj.IsNull() || obj.IsUnknown() {
		return nil, nil
	}

	targetConfig := graphmodels.NewCrossTenantAccessPolicyTargetConfiguration()

	attrs := obj.Attributes()

	if accessType, ok := attrs["access_type"].(types.String); ok && !accessType.IsNull() {
		accessTypeStr := accessType.ValueString()
		var accessTypeEnum *graphmodels.CrossTenantAccessPolicyTargetConfigurationAccessType
		switch accessTypeStr {
		case "allowed":
			val := graphmodels.ALLOWED_CROSSTENANTACCESSPOLICYTARGETCONFIGURATIONACCESSTYPE
			accessTypeEnum = &val
		case "blocked":
			val := graphmodels.BLOCKED_CROSSTENANTACCESSPOLICYTARGETCONFIGURATIONACCESSTYPE
			accessTypeEnum = &val
		}
		if accessTypeEnum != nil {
			targetConfig.SetAccessType(accessTypeEnum)
		}
	}

	if targetsSet, ok := attrs["targets"].(types.Set); ok && !targetsSet.IsNull() {
		var targets []graphmodels.CrossTenantAccessPolicyTargetable

		for _, elem := range targetsSet.Elements() {
			if targetObj, ok := elem.(types.Object); ok && !targetObj.IsNull() {
				target := graphmodels.NewCrossTenantAccessPolicyTarget()
				targetAttrs := targetObj.Attributes()

				if targetVal, ok := targetAttrs["target"].(types.String); ok && !targetVal.IsNull() {
					val := targetVal.ValueString()
					target.SetTarget(&val)
				}

				if targetType, ok := targetAttrs["target_type"].(types.String); ok && !targetType.IsNull() {
					targetTypeStr := targetType.ValueString()
					var targetTypeEnum *graphmodels.CrossTenantAccessPolicyTargetType
					switch targetTypeStr {
					case "user":
						val := graphmodels.USER_CROSSTENANTACCESSPOLICYTARGETTYPE
						targetTypeEnum = &val
					case "group":
						val := graphmodels.GROUP_CROSSTENANTACCESSPOLICYTARGETTYPE
						targetTypeEnum = &val
					case "application":
						val := graphmodels.APPLICATION_CROSSTENANTACCESSPOLICYTARGETTYPE
						targetTypeEnum = &val
					}
					if targetTypeEnum != nil {
						target.SetTargetType(targetTypeEnum)
					}
				}

				targets = append(targets, target)
			}
		}

		if len(targets) > 0 {
			targetConfig.SetTargets(targets)
		}
	}

	return targetConfig, nil
}

func constructInboundTrust(_ context.Context, obj types.Object) (graphmodels.CrossTenantAccessPolicyInboundTrustable, error) {
	if obj.IsNull() || obj.IsUnknown() {
		return nil, nil
	}

	inboundTrust := graphmodels.NewCrossTenantAccessPolicyInboundTrust()

	attrs := obj.Attributes()

	if isMfaAccepted, ok := attrs["is_mfa_accepted"].(types.Bool); ok && !isMfaAccepted.IsNull() {
		val := isMfaAccepted.ValueBool()
		inboundTrust.SetIsMfaAccepted(&val)
	}

	if isCompliantDeviceAccepted, ok := attrs["is_compliant_device_accepted"].(types.Bool); ok && !isCompliantDeviceAccepted.IsNull() {
		val := isCompliantDeviceAccepted.ValueBool()
		inboundTrust.SetIsCompliantDeviceAccepted(&val)
	}

	if isHybridAzureADJoinedDeviceAccepted, ok := attrs["is_hybrid_azure_ad_joined_device_accepted"].(types.Bool); ok && !isHybridAzureADJoinedDeviceAccepted.IsNull() {
		val := isHybridAzureADJoinedDeviceAccepted.ValueBool()
		inboundTrust.SetIsHybridAzureADJoinedDeviceAccepted(&val)
	}

	return inboundTrust, nil
}

func constructAutomaticUserConsentSettings(_ context.Context, obj types.Object) (graphmodels.InboundOutboundPolicyConfigurationable, error) {
	if obj.IsNull() || obj.IsUnknown() {
		return nil, nil
	}

	autoConsent := graphmodels.NewInboundOutboundPolicyConfiguration()

	attrs := obj.Attributes()

	if inboundAllowed, ok := attrs["inbound_allowed"].(types.Bool); ok && !inboundAllowed.IsNull() {
		val := inboundAllowed.ValueBool()
		autoConsent.SetInboundAllowed(&val)
	}

	if outboundAllowed, ok := attrs["outbound_allowed"].(types.Bool); ok && !outboundAllowed.IsNull() {
		val := outboundAllowed.ValueBool()
		autoConsent.SetOutboundAllowed(&val)
	}

	return autoConsent, nil
}

func constructTenantRestrictions(ctx context.Context, obj types.Object) (graphmodels.CrossTenantAccessPolicyTenantRestrictionsable, error) {
	if obj.IsNull() || obj.IsUnknown() {
		return nil, nil
	}

	// Use the base B2BSetting constructor to avoid setting @odata.type
	// The partner settings endpoint rejects the specialized tenant restrictions type
	tenantRestrictions := &tenantRestrictionsB2BSetting{
		CrossTenantAccessPolicyB2BSetting: *graphmodels.NewCrossTenantAccessPolicyB2BSetting(),
	}

	attrs := obj.Attributes()

	if usersAndGroupsObj, ok := attrs["users_and_groups"].(types.Object); ok && !usersAndGroupsObj.IsNull() {
		targetConfig, err := constructTargetConfiguration(ctx, usersAndGroupsObj)
		if err != nil {
			return nil, fmt.Errorf("failed to construct users_and_groups: %w", err)
		}
		tenantRestrictions.SetUsersAndGroups(targetConfig)
	}

	if applicationsObj, ok := attrs["applications"].(types.Object); ok && !applicationsObj.IsNull() {
		targetConfig, err := constructTargetConfiguration(ctx, applicationsObj)
		if err != nil {
			return nil, fmt.Errorf("failed to construct applications: %w", err)
		}
		tenantRestrictions.SetApplications(targetConfig)
	}

	return tenantRestrictions, nil
}
