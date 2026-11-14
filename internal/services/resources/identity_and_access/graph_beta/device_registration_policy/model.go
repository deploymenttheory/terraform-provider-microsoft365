// REF: https://learn.microsoft.com/en-us/graph/api/resources/deviceregistrationpolicy?view=graph-rest-beta
package graphBetaIdentityAndAccessDeviceRegistrationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DeviceRegistrationPolicyResourceModel represents the schema for the DeviceRegistrationPolicy resource.
// The nested blocks are types.Object because they are Optional+Computed: when they are not configured
// they are unknown in the plan and are filled from the current tenant policy.
type DeviceRegistrationPolicyResourceModel struct {
	ID                           types.String   `tfsdk:"id"`
	DisplayName                  types.String   `tfsdk:"display_name"`
	Description                  types.String   `tfsdk:"description"`
	UserDeviceQuota              types.Int32    `tfsdk:"user_device_quota"`
	MultiFactorAuthConfiguration types.String   `tfsdk:"multi_factor_auth_configuration"`
	LocalAdminPassword           types.Object   `tfsdk:"local_admin_password"`
	AzureADRegistration          types.Object   `tfsdk:"azure_ad_registration"`
	AzureADJoin                  types.Object   `tfsdk:"azure_ad_join"`
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`
}

// LocalAdminPasswordModel represents the localAdminPassword (Windows LAPS) settings.
type LocalAdminPasswordModel struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

// AzureADRegistrationModel represents the azureADRegistration (Microsoft Entra registered) settings.
type AzureADRegistrationModel struct {
	IsAdminConfigurable types.Bool   `tfsdk:"is_admin_configurable"`
	AllowedToRegister   types.Object `tfsdk:"allowed_to_register"`
}

// AzureADJoinModel represents the azureADJoin (Microsoft Entra joined) settings.
type AzureADJoinModel struct {
	IsAdminConfigurable types.Bool   `tfsdk:"is_admin_configurable"`
	AllowedToJoin       types.Object `tfsdk:"allowed_to_join"`
	LocalAdmins         types.Object `tfsdk:"local_admins"`
}

// LocalAdminsModel represents the local administrator settings for Microsoft Entra joined devices.
type LocalAdminsModel struct {
	EnableGlobalAdmins types.Bool   `tfsdk:"enable_global_admins"`
	RegisteringUsers   types.Object `tfsdk:"registering_users"`
}

// MembershipModel represents a polymorphic deviceRegistrationMembership.
// type maps to the @odata.type: all, none or enumerated (users and groups).
type MembershipModel struct {
	Type   types.String `tfsdk:"type"`
	Users  types.Set    `tfsdk:"users"`
	Groups types.Set    `tfsdk:"groups"`
}

const (
	membershipTypeAll        = "all"
	membershipTypeNone       = "none"
	membershipTypeEnumerated = "enumerated"

	odataTypeAllMembership        = "#microsoft.graph.allDeviceRegistrationMembership"
	odataTypeNoMembership         = "#microsoft.graph.noDeviceRegistrationMembership"
	odataTypeEnumeratedMembership = "#microsoft.graph.enumeratedDeviceRegistrationMembership"
)

var membershipAttrTypes = map[string]attr.Type{
	"type":   types.StringType,
	"users":  types.SetType{ElemType: types.StringType},
	"groups": types.SetType{ElemType: types.StringType},
}

var localAdminPasswordAttrTypes = map[string]attr.Type{
	"is_enabled": types.BoolType,
}

var azureADRegistrationAttrTypes = map[string]attr.Type{
	"is_admin_configurable": types.BoolType,
	"allowed_to_register":   types.ObjectType{AttrTypes: membershipAttrTypes},
}

var localAdminsAttrTypes = map[string]attr.Type{
	"enable_global_admins": types.BoolType,
	"registering_users":    types.ObjectType{AttrTypes: membershipAttrTypes},
}

var azureADJoinAttrTypes = map[string]attr.Type{
	"is_admin_configurable": types.BoolType,
	"allowed_to_join":       types.ObjectType{AttrTypes: membershipAttrTypes},
	"local_admins":          types.ObjectType{AttrTypes: localAdminsAttrTypes},
}
