package schema

import (
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// AndroidNotificationAssignmentsSchema returns the schema for the assignments block. This is used by the resources
// Android Enterprise Notifications and Android device administrator Notifications
func AndroidNotificationAssignmentsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Assignments for the compliance policy. Each assignment specifies the target group and schedule for script execution.",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				// Target assignment fields - only one should be used at a time
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Type of assignment target. Must be one of: 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.OneOf(
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
