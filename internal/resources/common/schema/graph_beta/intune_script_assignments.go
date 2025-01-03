package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func IntuneScriptAssignmentsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: "The assignment configuration for this Windows Settings Catalog profile.",
		Attributes: map[string]schema.Attribute{
			"all_devices": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Specifies whether this assignment applies to all devices. " +
					"When set to `true`, the assignment targets all devices in the organization." +
					"Can be used in conjuction with `all_users`." +
					"Can be used as an alternative to `include_groups`." +
					"Can be used in conjuction with `all_users` and `exclude_group_ids`.",
			},
			"all_users": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Specifies whether this assignment applies to all users. " +
					"When set to `true`, the assignment targets all licensed users within the organization." +
					"Can be used in conjuction with `all_devices`." +
					"Can be used as an alternative to `include_groups`." +
					"Can be used in conjuction with `all_devices` and `exclude_group_ids`.",
			},
			"include_group_ids": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "A set of entra id group Id's to include in the assignment.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							guidRegex,
							"assignment include group id must be a valid GUID",
						),
					),
				},
			},
			"exclude_group_ids": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				MarkdownDescription: "A set of group IDs to exclude from the assignment. " +
					"These groups will not receive the assignment, even if they match other inclusion criteria.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
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
