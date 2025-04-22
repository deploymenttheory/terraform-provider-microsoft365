package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func RoleAssignmentsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Optional:            true,
		MarkdownDescription: "The Role Assignment configurations for managing role assignments in Microsoft 365.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					MarkdownDescription: "Key of the entity. This is read-only and automatically generated.",
					Computed:            true,
				},
				"display_name": schema.StringAttribute{
					MarkdownDescription: "The display or friendly name of the role Assignment.",
					Required:            true,
				},
				"description": schema.StringAttribute{
					MarkdownDescription: "Description of the Role Assignment.",
					Optional:            true,
				},
				"admin_group_users_group_ids": schema.SetAttribute{
					MarkdownDescription: "Group ids that are assigned as members of this role scope.",
					Optional:            true,
					ElementType:         types.StringType,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(
							stringvalidator.RegexMatches(
								guidRegex,
								"scope member id must be a valid GUID",
							),
						),
					},
				},
				"scope_type": schema.StringAttribute{
					MarkdownDescription: "Administrators in this role assignment can target policies, applications and remote tasks to a scope type of:" +
						"'AllDevices', 'AllLicensedUsers', and 'AllDevicesAndLicensedUsers'. If the scope intent is for a entra id group then leave this empty. " +
						"Possible values are: `AllDevices`, `AllLicensedUsers`, `AllDevicesAndLicensedUsers`.",
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							"AllDevices",
							"AllLicensedUsers",
							"AllDevicesAndLicensedUsers",
						),
					},
				},
				"resource_scopes": schema.SetAttribute{
					MarkdownDescription: "Administrators in this role assignment can target policies, applications and remote tasks. List of ids of role scope member security groups from Entra ID.",
					Optional:            true,
					ElementType:         types.StringType,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(
							stringvalidator.RegexMatches(
								guidRegex,
								"resource scope id must be a valid GUID",
							),
						),
					},
				},
			},
		},
	}
}
