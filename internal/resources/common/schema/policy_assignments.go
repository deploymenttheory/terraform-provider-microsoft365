package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AssignmentsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: "The assignment configuration for this Windows Settings Catalog profile.",
		Attributes: map[string]schema.Attribute{
			"all_devices": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Specifies whether this assignment applies to all devices. " +
					"When set to `true`, the assignment targets all devices in the organization." +
					"Can be used in conjuction with `all_devices_filter_type` or `all_devices_filter_id`." +
					"Can be used as an alternative to `include_groups`." +
					"Can be used in conjuction with `all_users` and `all_users_filter_type` or `all_users_filter_id`.",
			},
			"all_devices_filter_type": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The filter type for all devices assignment. " +
					"Valid values are:\n" +
					"- `include`: Apply the assignment to devices that match the filter.\n" +
					"- `exclude`: Do not apply the assignment to devices that match the filter.\n" +
					"- `none`: No filter applied.",
				Validators: []validator.String{
					stringvalidator.OneOf("include", "exclude", "none"),
				},
			},
			"all_devices_filter_id": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The ID of the device group filter to apply when `all_devices` is set to `true`. " +
					"This should be a valid GUID of an existing device group filter.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						guidRegex,
						"must be a valid GUID",
					),
				},
			},
			"all_users": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Specifies whether this assignment applies to all users. " +
					"When set to `true`, the assignment targets all licensed users within the organization." +
					"Can be used in conjuction with `all_users_filter_type` or `all_users_filter_id`." +
					"Can be used as an alternative to `include_groups`." +
					"Can be used in conjuction with `all_devices` and `all_devices_filter_type` or `all_devices_filter_id`.",
			},
			"all_users_filter_type": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The filter type for all users assignment. " +
					"Valid values are:\n" +
					"- `include`: Apply the assignment to users that match the filter.\n" +
					"- `exclude`: Do not apply the assignment to users that match the filter.\n" +
					"- `none`: No filter applied.",
				Validators: []validator.String{
					stringvalidator.OneOf("include", "exclude", "none"),
				},
			},
			"all_users_filter_id": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The ID of the filter to apply when `all_users` is set to `true`. " +
					"This should be a valid GUID of an existing filter.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						guidRegex,
						"assignment filer for all_users must be a valid GUID",
					),
				},
			},
			"include_groups": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "A list of entra id group Id's to include in the assignment. " +
					"Each group can have its own filter type and filter ID.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The entra ID group ID of the group to include in the assignment. " +
								"This should be a valid GUID of an existing group.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									guidRegex,
									"assignment include group(s) must be a valid GUID",
								),
							},
						},
						"include_groups_filter_type": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The device group filter type for the included group. " +
								"Valid values are:\n" +
								"- `include`: Apply the assignment to group members that match the filter.\n" +
								"- `exclude`: Do not apply the assignment to group members that match the filter.\n" +
								"- `none`: No filter applied.",
							Validators: []validator.String{
								stringvalidator.OneOf("include", "exclude", "none"),
							},
						},
						"include_groups_filter_id": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The Entra ID Group ID of the filter to apply to the included group. " +
								"This should be a valid GUID of an existing filter.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									guidRegex,
									"assignment group filter id must be a valid GUID",
								),
							},
						},
					},
				},
			},
			"exclude_group_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				MarkdownDescription: "A list of group IDs to exclude from the assignment. " +
					"These groups will not receive the assignment, even if they match other inclusion criteria.",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							guidRegex,
							"assignment exclude group id must be a valid GUID",
						),
					),
				},
			},
		},
	}
}
