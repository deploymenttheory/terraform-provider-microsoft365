package graphBetaIdentityAndAccessDeviceRegistrationPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote device registration policy from Kiota SDK to Terraform state.
// The prior values in data are used only to keep an unset users or groups list null when Graph returns an empty list.
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data DeviceRegistrationPolicyResourceModel,
	remote graphmodels.DeviceRegistrationPolicyable,
) DeviceRegistrationPolicyResourceModel {
	if remote == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": remote.GetId(),
	})

	data.ID = convert.GraphToFrameworkStringWithDefault(remote.GetId(), singletonID)
	data.DisplayName = convert.GraphToFrameworkString(remote.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.UserDeviceQuota = convert.GraphToFrameworkInt32(remote.GetUserDeviceQuota())
	data.MultiFactorAuthConfiguration = convert.GraphToFrameworkEnum(remote.GetMultiFactorAuthConfiguration())
	data.LocalAdminPassword = mapLocalAdminPassword(remote.GetLocalAdminPassword())
	data.AzureADRegistration = mapAzureADRegistration(ctx, data.AzureADRegistration, remote.GetAzureADRegistration())
	data.AzureADJoin = mapAzureADJoin(ctx, data.AzureADJoin, remote.GetAzureADJoin())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
	return data
}

// mapLocalAdminPassword maps the Windows LAPS settings.
func mapLocalAdminPassword(remote graphmodels.LocalAdminPasswordSettingsable) types.Object {
	if remote == nil {
		return types.ObjectNull(localAdminPasswordAttrTypes)
	}
	return types.ObjectValueMust(localAdminPasswordAttrTypes, map[string]attr.Value{
		"is_enabled": convert.GraphToFrameworkBool(remote.GetIsEnabled()),
	})
}

// mapAzureADRegistration maps the Microsoft Entra registered device settings.
func mapAzureADRegistration(ctx context.Context, prior attr.Value, remote graphmodels.AzureADRegistrationPolicyable) types.Object {
	if remote == nil {
		return types.ObjectNull(azureADRegistrationAttrTypes)
	}
	priorMembership := priorAttribute(prior, "allowed_to_register")
	return types.ObjectValueMust(azureADRegistrationAttrTypes, map[string]attr.Value{
		"is_admin_configurable": convert.GraphToFrameworkBool(remote.GetIsAdminConfigurable()),
		"allowed_to_register":   mapMembership(ctx, priorMembership, remote.GetAllowedToRegister()),
	})
}

// mapAzureADJoin maps the Microsoft Entra joined device settings.
func mapAzureADJoin(ctx context.Context, prior attr.Value, remote graphmodels.AzureADJoinPolicyable) types.Object {
	if remote == nil {
		return types.ObjectNull(azureADJoinAttrTypes)
	}

	localAdmins := types.ObjectNull(localAdminsAttrTypes)
	if remoteLocalAdmins := remote.GetLocalAdmins(); remoteLocalAdmins != nil {
		priorRegisteringUsers := priorAttribute(priorAttribute(prior, "local_admins"), "registering_users")
		localAdmins = types.ObjectValueMust(localAdminsAttrTypes, map[string]attr.Value{
			"enable_global_admins": convert.GraphToFrameworkBool(remoteLocalAdmins.GetEnableGlobalAdmins()),
			"registering_users":    mapMembership(ctx, priorRegisteringUsers, remoteLocalAdmins.GetRegisteringUsers()),
		})
	}

	return types.ObjectValueMust(azureADJoinAttrTypes, map[string]attr.Value{
		"is_admin_configurable": convert.GraphToFrameworkBool(remote.GetIsAdminConfigurable()),
		"allowed_to_join":       mapMembership(ctx, priorAttribute(prior, "allowed_to_join"), remote.GetAllowedToJoin()),
		"local_admins":          localAdmins,
	})
}

// mapMembership maps a polymorphic deviceRegistrationMembership to the type/users/groups object.
// users and groups are only set for enumerated memberships. Graph returns an empty list for an
// unset users or groups list, so an empty list stays null unless the prior value was an explicit empty set.
func mapMembership(ctx context.Context, prior attr.Value, remote graphmodels.DeviceRegistrationMembershipable) types.Object {
	if remote == nil {
		return types.ObjectNull(membershipAttrTypes)
	}

	users := types.SetNull(types.StringType)
	groups := types.SetNull(types.StringType)
	membershipType := membershipTypeFromODataType(remote.GetOdataType())

	if enumerated, ok := remote.(graphmodels.EnumeratedDeviceRegistrationMembershipable); ok {
		membershipType = membershipTypeEnumerated
		users = mapMembershipSet(ctx, priorAttribute(prior, "users"), enumerated.GetUsers())
		groups = mapMembershipSet(ctx, priorAttribute(prior, "groups"), enumerated.GetGroups())
	}

	return types.ObjectValueMust(membershipAttrTypes, map[string]attr.Value{
		"type":   types.StringValue(membershipType),
		"users":  users,
		"groups": groups,
	})
}

// mapMembershipSet keeps an unset list null when Graph returns it empty.
func mapMembershipSet(ctx context.Context, prior attr.Value, values []string) types.Set {
	if len(values) > 0 {
		return convert.GraphToFrameworkStringSet(ctx, values)
	}
	if priorSet, ok := prior.(types.Set); ok && !priorSet.IsNull() && !priorSet.IsUnknown() {
		return types.SetValueMust(types.StringType, []attr.Value{})
	}
	return types.SetNull(types.StringType)
}
