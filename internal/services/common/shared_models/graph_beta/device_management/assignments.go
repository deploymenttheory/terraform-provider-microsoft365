package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

// DeviceCompliancePolicyAssignmentResourceModel struct to hold device compliance policy assignment configuration
type DeviceCompliancePolicyAssignmentResourceModel struct {
	// Target assignment fields - only one should be used at a time
	Type    types.String `tfsdk:"type"`     // "allDevicesAssignmentTarget", "allLicensedUsersAssignmentTarget", "groupAssignmentTarget", "exclusionGroupAssignmentTarget"
	GroupId types.String `tfsdk:"group_id"` // For group targets (both include and exclude)
	// Assignment filter fields
	FilterId   types.String `tfsdk:"filter_id"`
	FilterType types.String `tfsdk:"filter_type"` // "include", "exclude", or "none"
}

// DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel defines the schema for a Windows Remediation Script assignment.
type DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel struct {
	Type       types.String `tfsdk:"type"`     // "allDevicesAssignmentTarget", "allLicensedUsersAssignmentTarget", "groupAssignmentTarget", "exclusionGroupAssignmentTarget"
	GroupId    types.String `tfsdk:"group_id"` // For group targets (both include and exclude)
	FilterId   types.String `tfsdk:"filter_id"`
	FilterType types.String `tfsdk:"filter_type"` // "include", "exclude", or "none"
}

// DeviceManagementDeviceConfigurationAssignmentWithoutGroupFilterModel defines the schema for a Windows Remediation Script assignment.
type DeviceManagementDeviceConfigurationAssignmentWithoutGroupFilterModel struct {
	Type    types.String `tfsdk:"type"`     // "allDevicesAssignmentTarget", "allLicensedUsersAssignmentTarget", "groupAssignmentTarget", "exclusionGroupAssignmentTarget"
	GroupId types.String `tfsdk:"group_id"` // For group targets (both include and exclude)
}

// DeviceManagementDeviceConfigurationAssignmentWithAllLicensedUsersInclusionGroupConfigurationManagerCollectionAssignmentModel defines the schema for a Windows Remediation Script assignment.
type DeviceManagementDeviceConfigurationAssignmentWithAllLicensedUsersInclusionGroupConfigurationManagerCollectionAssignmentModel struct {
	Type         types.String `tfsdk:"type"`          // "allDevicesAssignmentTarget", "allLicensedUsersAssignmentTarget", "groupAssignmentTarget", "exclusionGroupAssignmentTarget"
	GroupId      types.String `tfsdk:"group_id"`      // For group targets (both include and exclude)
	CollectionId types.String `tfsdk:"collection_id"` // For configuration manager collection targets
}

// WindowsSoftwareUpdateAssignmentModel defines the schema for a Windows Software Update assignment.
type WindowsSoftwareUpdateAssignmentModel struct {
	Type    types.String `tfsdk:"type"`     // "groupAssignmentTarget", "exclusionGroupAssignmentTarget"
	GroupId types.String `tfsdk:"group_id"` // For group targets (both include and exclude)
}

// PlatformScriptAssignmentModel defines the schema for a MacOS Platform Script assignment.
type PlatformScriptAssignmentModel struct {
	// Target assignment fields - only one should be used at a time
	Type    types.String `tfsdk:"type"`     // "allDevicesAssignmentTarget", "allLicensedUsersAssignmentTarget", "groupAssignmentTarget", "exclusionGroupAssignmentTarget"
	GroupId types.String `tfsdk:"group_id"` // For group targets (both include and exclude)
}
