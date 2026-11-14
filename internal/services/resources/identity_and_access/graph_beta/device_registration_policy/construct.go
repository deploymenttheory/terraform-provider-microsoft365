package graphBetaIdentityAndAccessDeviceRegistrationPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource converts the Terraform resource model to a Kiota SDK model for the full-replace PUT.
//
// PUT /policies/deviceRegistrationPolicy replaces the whole policy and requires every updatable property,
// so each value that is unknown or null in the plan (not configured) is taken from remoteResource, the
// policy read immediately before the write. Fresh SDK objects are built for every value so the Kiota
// backing store serializes all of them, including values copied unchanged from the remote policy.
func constructResource(
	ctx context.Context,
	data *DeviceRegistrationPolicyResourceModel,
	remoteResource graphmodels.DeviceRegistrationPolicyable,
) (graphmodels.DeviceRegistrationPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if remoteResource == nil {
		return nil, errRemotePolicyRequired
	}

	if err := validateRequest(ctx, data, remoteResource); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := graphmodels.NewDeviceRegistrationPolicy()

	if isKnown(data.UserDeviceQuota) {
		convert.FrameworkToGraphInt32(data.UserDeviceQuota, requestBody.SetUserDeviceQuota)
	} else {
		requestBody.SetUserDeviceQuota(remoteResource.GetUserDeviceQuota())
	}

	if isKnown(data.MultiFactorAuthConfiguration) {
		err := convert.FrameworkToGraphEnum(data.MultiFactorAuthConfiguration, graphmodels.ParseMultiFactorAuthConfiguration, requestBody.SetMultiFactorAuthConfiguration)
		if err != nil {
			return nil, fmt.Errorf("failed to set multi factor auth configuration: %w", err)
		}
	} else {
		requestBody.SetMultiFactorAuthConfiguration(remoteResource.GetMultiFactorAuthConfiguration())
	}

	localAdminPassword, err := constructLocalAdminPassword(ctx, data.LocalAdminPassword, remoteResource.GetLocalAdminPassword())
	if err != nil {
		return nil, fmt.Errorf("failed to construct local admin password: %w", err)
	}
	requestBody.SetLocalAdminPassword(localAdminPassword)

	azureADRegistration, err := constructAzureADRegistration(ctx, data.AzureADRegistration, remoteResource.GetAzureADRegistration())
	if err != nil {
		return nil, fmt.Errorf("failed to construct azure ad registration: %w", err)
	}
	requestBody.SetAzureADRegistration(azureADRegistration)

	azureADJoin, err := constructAzureADJoin(ctx, data.AzureADJoin, remoteResource.GetAzureADJoin())
	if err != nil {
		return nil, fmt.Errorf("failed to construct azure ad join: %w", err)
	}
	requestBody.SetAzureADJoin(azureADJoin)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructLocalAdminPassword builds the Windows LAPS settings.
func constructLocalAdminPassword(ctx context.Context, data types.Object, remote graphmodels.LocalAdminPasswordSettingsable) (graphmodels.LocalAdminPasswordSettingsable, error) {
	settings := graphmodels.NewLocalAdminPasswordSettings()

	if isKnown(data) {
		var model LocalAdminPasswordModel
		if diags := data.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, fmt.Errorf("%w local_admin_password: %v", errReadPlanValue, diags.Errors())
		}
		if isKnown(model.IsEnabled) {
			convert.FrameworkToGraphBool(model.IsEnabled, settings.SetIsEnabled)
			return settings, nil
		}
	}

	if remote != nil {
		settings.SetIsEnabled(remote.GetIsEnabled())
	}
	return settings, nil
}

// constructAzureADRegistration builds the Microsoft Entra registered device settings.
// isAdminConfigurable is read-only and is always sent back with its current value.
func constructAzureADRegistration(ctx context.Context, data types.Object, remote graphmodels.AzureADRegistrationPolicyable) (graphmodels.AzureADRegistrationPolicyable, error) {
	policy := graphmodels.NewAzureADRegistrationPolicy()

	var remoteMembership graphmodels.DeviceRegistrationMembershipable
	if remote != nil {
		policy.SetIsAdminConfigurable(remote.GetIsAdminConfigurable())
		remoteMembership = remote.GetAllowedToRegister()
	}

	var model AzureADRegistrationModel
	if isKnown(data) {
		if diags := data.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, fmt.Errorf("%w azure_ad_registration: %v", errReadPlanValue, diags.Errors())
		}
	}

	membership, err := constructMembership(ctx, model.AllowedToRegister, remoteMembership)
	if err != nil {
		return nil, fmt.Errorf("failed to construct allowed_to_register: %w", err)
	}
	policy.SetAllowedToRegister(membership)

	return policy, nil
}

// constructAzureADJoin builds the Microsoft Entra joined device settings.
// isAdminConfigurable is read-only and is always sent back with its current value.
func constructAzureADJoin(ctx context.Context, data types.Object, remote graphmodels.AzureADJoinPolicyable) (graphmodels.AzureADJoinPolicyable, error) {
	policy := graphmodels.NewAzureADJoinPolicy()

	var remoteAllowedToJoin graphmodels.DeviceRegistrationMembershipable
	var remoteLocalAdmins graphmodels.LocalAdminSettingsable
	if remote != nil {
		policy.SetIsAdminConfigurable(remote.GetIsAdminConfigurable())
		remoteAllowedToJoin = remote.GetAllowedToJoin()
		remoteLocalAdmins = remote.GetLocalAdmins()
	}

	var model AzureADJoinModel
	if isKnown(data) {
		if diags := data.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, fmt.Errorf("%w azure_ad_join: %v", errReadPlanValue, diags.Errors())
		}
	}

	allowedToJoin, err := constructMembership(ctx, model.AllowedToJoin, remoteAllowedToJoin)
	if err != nil {
		return nil, fmt.Errorf("failed to construct allowed_to_join: %w", err)
	}
	policy.SetAllowedToJoin(allowedToJoin)

	localAdmins, err := constructLocalAdmins(ctx, model.LocalAdmins, remoteLocalAdmins)
	if err != nil {
		return nil, fmt.Errorf("failed to construct local_admins: %w", err)
	}
	policy.SetLocalAdmins(localAdmins)

	return policy, nil
}

// constructLocalAdmins builds the local administrator settings for Microsoft Entra joined devices.
func constructLocalAdmins(ctx context.Context, data types.Object, remote graphmodels.LocalAdminSettingsable) (graphmodels.LocalAdminSettingsable, error) {
	settings := graphmodels.NewLocalAdminSettings()

	var remoteRegisteringUsers graphmodels.DeviceRegistrationMembershipable
	if remote != nil {
		settings.SetEnableGlobalAdmins(remote.GetEnableGlobalAdmins())
		remoteRegisteringUsers = remote.GetRegisteringUsers()
	}

	var model LocalAdminsModel
	if isKnown(data) {
		if diags := data.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, fmt.Errorf("%w local_admins: %v", errReadPlanValue, diags.Errors())
		}
	}

	if isKnown(model.EnableGlobalAdmins) {
		convert.FrameworkToGraphBool(model.EnableGlobalAdmins, settings.SetEnableGlobalAdmins)
	}

	registeringUsers, err := constructMembership(ctx, model.RegisteringUsers, remoteRegisteringUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to construct registering_users: %w", err)
	}
	settings.SetRegisteringUsers(registeringUsers)

	return settings, nil
}

// constructMembership builds a polymorphic deviceRegistrationMembership from the plan, or copies the
// remote membership when the plan value is not configured. users and groups are only sent for
// enumerated memberships, and unset lists are omitted rather than sent as null.
func constructMembership(ctx context.Context, data types.Object, remote graphmodels.DeviceRegistrationMembershipable) (graphmodels.DeviceRegistrationMembershipable, error) {
	if !isKnown(data) {
		return copyMembership(remote), nil
	}

	var model MembershipModel
	if diags := data.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, fmt.Errorf("%w membership: %v", errReadPlanValue, diags.Errors())
	}

	switch model.Type.ValueString() {
	case membershipTypeAll:
		return graphmodels.NewAllDeviceRegistrationMembership(), nil
	case membershipTypeNone:
		return graphmodels.NewNoDeviceRegistrationMembership(), nil
	case membershipTypeEnumerated:
		membership := graphmodels.NewEnumeratedDeviceRegistrationMembership()
		if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, model.Users, membership.SetUsers); err != nil {
			return nil, fmt.Errorf("failed to convert users: %w", err)
		}
		if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, model.Groups, membership.SetGroups); err != nil {
			return nil, fmt.Errorf("failed to convert groups: %w", err)
		}
		return membership, nil
	default:
		return nil, fmt.Errorf("%w %q", errUnsupportedMembershipType, model.Type.ValueString())
	}
}

// copyMembership rebuilds a remote membership as a fresh SDK object so every value is serialized.
func copyMembership(remote graphmodels.DeviceRegistrationMembershipable) graphmodels.DeviceRegistrationMembershipable {
	if remote == nil {
		return nil
	}

	if enumerated, ok := remote.(graphmodels.EnumeratedDeviceRegistrationMembershipable); ok {
		membership := graphmodels.NewEnumeratedDeviceRegistrationMembership()
		if users := enumerated.GetUsers(); users != nil {
			membership.SetUsers(append([]string{}, users...))
		}
		if groups := enumerated.GetGroups(); groups != nil {
			membership.SetGroups(append([]string{}, groups...))
		}
		return membership
	}

	switch membershipTypeFromODataType(remote.GetOdataType()) {
	case membershipTypeAll:
		return graphmodels.NewAllDeviceRegistrationMembership()
	case membershipTypeNone:
		return graphmodels.NewNoDeviceRegistrationMembership()
	default:
		membership := graphmodels.NewDeviceRegistrationMembership()
		membership.SetOdataType(remote.GetOdataType())
		return membership
	}
}
