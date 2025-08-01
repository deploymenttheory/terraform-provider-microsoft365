package schema

import (
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ComplianceScriptAssignmentsSchema returns the schema for the assignments block
func ComplianceScriptAssignmentsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Assignments for the compliance policy. Each assignment specifies the target group and schedule for script execution.",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				// Target assignment fields - only one should be used at a time
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
