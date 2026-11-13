package schema

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// ForwardingProfileSettingsAttributes returns the writable traffic forwarding
// profile v2 settings shared by custom and default profiles.
func ForwardingProfileSettingsAttributes() map[string]schema.Attribute {
	idSet := func(kind string) schema.SetAttribute {
		return schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			MarkdownDescription: fmt.Sprintf(
				"Microsoft Entra %s object IDs. Allowed only when `membership` is `enumerated`.",
				kind,
			),
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID",
					),
				),
			},
		}
	}
	return map[string]schema.Attribute{
		"state": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Whether the profile is enabled. Values: `enabled`, `disabled`.",
			Validators:          []validator.String{stringvalidator.OneOf("enabled", "disabled")},
		},
		"client_fallback_action": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "What the client does with acquired traffic when it cannot reach Global Secure Access. Values: `bypass`, `block`. Graph defaults to `bypass`. Omitting it leaves the current value unchanged.",
			Validators:          []validator.String{stringvalidator.OneOf("bypass", "block")},
			PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"forwarding_source_matching": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Source matching conditions (preview). Omitting it leaves the current value unchanged.",
			PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
			Attributes: map[string]schema.Attribute{
				"session_type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Traffic sources the profile matches. `any` matches all sources; `agent` matches traffic from local agents only.",
					Validators:          []validator.String{stringvalidator.OneOf("any", "agent")},
				},
			},
		},
		"local_network_settings": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Home network support (preview). Supported only for `private` profiles. Omitting it leaves the current value unchanged.",
			PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
			Attributes: map[string]schema.Attribute{
				"allow_local_subnet_bypass": schema.BoolAttribute{
					Required:            true,
					MarkdownDescription: "Lets users turn on **Prefer Local Network** in the client so traffic to the directly connected IPv4 RFC 1918 subnet bypasses Private Access.",
				},
			},
		},
		"assignment_rules": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "User, group, device and platform assignments. Graph replaces assignments as a whole. Omitting it leaves the current assignments unchanged.",
			PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
			Attributes: map[string]schema.Attribute{
				"membership": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "`all` assigns all users and devices, `enumerated` assigns the listed IDs, `none` assigns nobody.",
					Validators: []validator.String{stringvalidator.OneOf(
						sharedmodels.MembershipAll,
						sharedmodels.MembershipEnumerated,
						sharedmodels.MembershipNone,
					)},
				},
				"user_ids":   idSet("user"),
				"group_ids":  idSet("security group"),
				"device_ids": idSet("device"),
				"platforms": schema.SetAttribute{
					ElementType:         types.StringType,
					Required:            true,
					MarkdownDescription: "Device platforms the profile applies to. Values: `windows`, `macOS`, `android`, `iOS`. At least one is required.",
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf(sharedmodels.ForwardingProfilePlatforms...),
						),
					},
				},
			},
		},
	}
}

// ValidateForwardingProfileSettings checks cross-attribute rules that Graph
// would otherwise reject or silently rewrite. trafficType may be unknown.
func ValidateForwardingProfileSettings(
	ctx context.Context,
	trafficType types.String,
	localNetwork types.Object,
	rules types.Object,
) diag.Diagnostics {
	var diags diag.Diagnostics
	if !localNetwork.IsNull() && !localNetwork.IsUnknown() &&
		!trafficType.IsNull() && !trafficType.IsUnknown() && trafficType.ValueString() != "private" {
		diags.AddAttributeError(
			path.Root("local_network_settings"),
			"Unsupported local network settings",
			"local_network_settings applies only to private traffic forwarding profiles.",
		)
	}
	if rules.IsNull() || rules.IsUnknown() {
		return diags
	}
	var m sharedmodels.AssignmentRulesModel
	if d := rules.As(
		ctx,
		&m,
		basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true},
	); d.HasError() {
		return append(diags, d...)
	}
	if m.Membership.IsUnknown() {
		return diags
	}
	sets := map[string]types.Set{
		"user_ids":   m.UserIDs,
		"group_ids":  m.GroupIDs,
		"device_ids": m.DeviceIDs,
	}
	if m.Membership.ValueString() != sharedmodels.MembershipEnumerated {
		for name, set := range sets {
			if !set.IsNull() {
				diags.AddAttributeError(
					path.Root("assignment_rules").AtName(name),
					"Unexpected assignment IDs",
					name+" can be set only when membership is \"enumerated\".",
				)
			}
		}
		return diags
	}
	total := 0
	for _, set := range sets {
		if set.IsUnknown() {
			return diags
		}
		if !set.IsNull() {
			total += len(set.Elements())
		}
	}
	if total == 0 {
		// Graph turns an empty enumerated membership into "none".
		diags.AddAttributeError(
			path.Root("assignment_rules"),
			"Missing assignment IDs",
			"membership \"enumerated\" requires at least one ID in user_ids, group_ids or device_ids. Use membership \"none\" to assign nobody.",
		)
	}
	return diags
}
