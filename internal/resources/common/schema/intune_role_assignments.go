package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func IntuneRoleAssignmentsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: "The Role Assignment configuration for managing role assignments in Microsoft 365.",
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
			"scope_members": schema.SetAttribute{
				MarkdownDescription: "List of ids of role scope member security groups. These are IDs from Azure Active Directory.",
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
				MarkdownDescription: "Specifies the type of scope for a Role Assignment. Default type 'ResourceScope' allows assignment of ResourceScopes. " +
					"For 'AllDevices', 'AllLicensedUsers', and 'AllDevicesAndLicensedUsers', the ResourceScopes property should be left empty. " +
					"Possible values are: `resourceScope`, `allDevices`, `allLicensedUsers`, `allDevicesAndLicensedUsers`.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"resourceScope",
						"allDevices",
						"allLicensedUsers",
						"allDevicesAndLicensedUsers",
					),
				},
			},
			"resource_scopes": schema.SetAttribute{
				MarkdownDescription: "List of ids of role scope member security groups. These are IDs from Azure Active Directory.",
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
	}
}
