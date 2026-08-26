package schema

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
)

// HardwareConfigurationAssignmentsSchema is a schema for hardware configuration assignments.
// The `/deviceManagement/hardwareConfigurations/{id}/assign` endpoint only accepts group
// inclusion and exclusion targets, matching the Intune 'BIOS configuration and other settings'
// blade, which offers included and excluded groups but no 'all devices' / 'all users' targets.
// Included group targets support assignment filters.
func HardwareConfigurationAssignmentsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Assignments for the hardware configuration. Each assignment targets an Entra ID group for inclusion or exclusion. Supports assignment filters.",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Type of assignment target. Must be one of: 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.OneOf(
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
					MarkdownDescription: "ID of the filter to apply to the assignment. Required when filter_type is 'include' or 'exclude'. Should be omitted when filter_type is 'none'.",
					// Have to set a default value here to satify the set hash calculation behaviour.
					Default: stringdefault.StaticString("00000000-0000-0000-0000-000000000000"),
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
						customValidator.RequiredWhenEquals(
							"filter_type",
							types.StringValue("include"),
						),
						customValidator.RequiredWhenEquals(
							"filter_type",
							types.StringValue("exclude"),
						),
						stringvalidator.NoneOf("00000000-0000-0000-0000-000000000000"),
					},
				},
				"filter_type": schema.StringAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.",
					Default:             stringdefault.StaticString("none"),
					Validators: []validator.String{
						stringvalidator.OneOf("include", "exclude", "none"),
					},
				},
			},
		},
	}
}
