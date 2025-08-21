package schema

import (
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema is a schema for device configuration
// assignments that supports all group assignment types, and group filters.
func DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Assignments for the device configuration. Each assignment specifies the target group and schedule for script execution. Supports group filters.",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"allDevicesAssignmentTarget",
							"allLicensedUsersAssignmentTarget",
							"groupAssignmentTarget",
							"exclusionGroupAssignmentTarget",
						),
					},
				},
				"group_id": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					},
				},
				// Assignment filter fields
				"filter_id": schema.StringAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "ID of the filter to apply to the assignment.",
					Default:             stringdefault.StaticString("00000000-0000-0000-0000-000000000000"),
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					},
				},
				"filter_type": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.",
					Computed:            true,
					Default:             stringdefault.StaticString("none"),
					Validators: []validator.String{
						stringvalidator.OneOf("include", "exclude", "none"),
					},
				},
			},
		},
	}
}

// DeviceConfigurationWithAllGroupAssignmentsSchema is a schema for device configuration
// assignments that supports all group types, but no group filters.
func DeviceConfigurationWithAllGroupAssignmentsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Assignments for the device configuration. Each assignment specifies the target group and schedule for script execution.",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"allDevicesAssignmentTarget",
							"allLicensedUsersAssignmentTarget",
							"groupAssignmentTarget",
							"exclusionGroupAssignmentTarget",
						),
					},
				},
				"group_id": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					},
				},
			},
		},
	}
}

// DeviceConfigurationWithAllLicensedUsersAllDevicesInclusionGroupAssignmentsSchema is a schema for device configuration
// assignments that supports inclusion group types, but no exclusion groups types and nogroup filters.
func DeviceConfigurationWithAllLicensedUsersAllDevicesInclusionGroupAssignmentsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Assignments for the device configuration. Each assignment specifies the target group and schedule for script execution.",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"allDevicesAssignmentTarget",
							"allLicensedUsersAssignmentTarget",
							"groupAssignmentTarget",
						),
					},
				},
				"group_id": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					},
				},
			},
		},
	}
}

// DeviceConfigurationWithAllLicensedUsersInclusionGroupConfigurationManagerCollectionAssignmentsSchema is a schema for device configuration
// assignments that only support inclusion group assignments.
func DeviceConfigurationWithAllLicensedUsersInclusionGroupConfigurationManagerCollectionAssignmentsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Assignments for the device configuration. Each assignment specifies the target group and schedule for script execution.",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required: true,
					MarkdownDescription: "The target group type for the assignment. Possible values are:\n\n" +
						"- **allLicensedUsersAssignmentTarget**: Target all users with Intune licenses\n" +
						"- **groupAssignmentTarget**: Target a specific Entra ID group\n" +
						"- **configurationManagerCollection**: Target System Center Configuration Manager collection",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"allLicensedUsersAssignmentTarget",
							"groupAssignmentTarget",
							"configurationManagerCollection",
						),
					},
				},
				"group_id": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					},
				},
				"collection_id": schema.StringAttribute{
					MarkdownDescription: "The SCCM group collection ID for the assignment target. Default collections start with 'SMS', while custom collections start with your site code (e.g., 'MEM').",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[A-Za-z]{2,8}[0-9A-Za-z]{8}$`),
							"Must be a valid SCCM collection ID format. Default collections start with 'SMS' followed by an alphanumeric ID. Custom collections start with your site code (e.g., 'MEM') followed by an alphanumeric ID.",
						),
					},
				},
			},
		},
	}
}

// DeviceConfigurationWithInclusionGroupAssignmentsSchema is a schema for device configuration
// assignments that only support inclusion group assignments.
func DeviceConfigurationWithInclusionGroupAssignmentsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Assignments for the device configuration. Each assignment specifies the target group and schedule for script execution.",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Type of assignment target. Must be one of: 'groupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"groupAssignmentTarget",
						),
					},
				},
				"group_id": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					},
				},
			},
		},
	}
}
