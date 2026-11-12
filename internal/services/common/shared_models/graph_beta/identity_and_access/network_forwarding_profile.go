package sharedmodels

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Graph beta forwarding profile v2 properties. They are not in the public
// $metadata or the generated SDK; shapes follow Entra portal traffic.
const (
	odataAssignmentsAll        = "#microsoft.graph.networkaccess.forwardingProfileAssignmentsAllMembership"
	odataAssignmentsEnumerated = "#microsoft.graph.networkaccess.forwardingProfileAssignmentsEnumeratedMembership"
	odataAssignmentsNone       = "#microsoft.graph.networkaccess.forwardingProfileAssignmentsNoMembership"

	MembershipAll        = "all"
	MembershipEnumerated = "enumerated"
	MembershipNone       = "none"

	// platformsNone is what Graph returns when no platform is selected. Graph
	// rejects it on PATCH, so it is only ever read.
	platformsNone = "none"
)

// ForwardingProfilePlatforms lists the platform values in Graph's canonical casing.
var ForwardingProfilePlatforms = []string{"windows", "macOS", "android", "iOS"}

// ForwardingProfileWire is the JSON shape of networkaccess/forwardingProfiles.
type ForwardingProfileWire struct {
	ID                       *string                         `json:"id,omitempty"`
	Name                     *string                         `json:"name,omitempty"`
	Description              *string                         `json:"description,omitempty"`
	State                    *string                         `json:"state,omitempty"`
	Version                  *string                         `json:"version,omitempty"`
	LastModifiedDateTime     *string                         `json:"lastModifiedDateTime,omitempty"`
	TrafficForwardingType    *string                         `json:"trafficForwardingType,omitempty"`
	Priority                 *int32                          `json:"priority,omitempty"`
	IsCustomProfile          *bool                           `json:"isCustomProfile,omitempty"`
	ClientFallbackAction     *string                         `json:"clientFallbackAction,omitempty"`
	LocalNetworkSettings     *ForwardingProfileLocalNetwork  `json:"localNetworkSettings,omitempty"`
	ForwardingSourceMatching *ForwardingProfileSourceMatch   `json:"forwardingSourceMatching,omitempty"`
	AssignmentRules          *ForwardingProfileAssignmentsRW `json:"assignmentRules,omitempty"`
}

type ForwardingProfileLocalNetwork struct {
	AllowLocalSubnetBypass *bool `json:"allowLocalSubnetBypass"`
}

type ForwardingProfileSourceMatch struct {
	SessionType *string `json:"sessionType"`
}

type ForwardingProfileAssignmentsRW struct {
	Assignments           *ForwardingProfileAssignmentsWire `json:"assignments"`
	DevicePropertyFilters *ForwardingProfileDeviceFilters   `json:"devicePropertyFilters"`
}

type ForwardingProfileAssignmentsWire struct {
	ODataType string                         `json:"@odata.type"` //nolint:tagliatelle // Graph wire name
	Users     []ForwardingProfileDirectoryID `json:"users,omitempty"`
	Groups    []ForwardingProfileDirectoryID `json:"groups,omitempty"`
	Devices   []ForwardingProfileDirectoryID `json:"devices,omitempty"`
}

type ForwardingProfileDirectoryID struct {
	ID string `json:"id"`
}

type ForwardingProfileDeviceFilters struct {
	Platforms *string `json:"platforms"`
}

// Terraform object models for the shared nested attributes.

type ForwardingSourceMatchingModel struct {
	SessionType types.String `tfsdk:"session_type"`
}

type LocalNetworkSettingsModel struct {
	AllowLocalSubnetBypass types.Bool `tfsdk:"allow_local_subnet_bypass"`
}

type AssignmentRulesModel struct {
	Membership types.String `tfsdk:"membership"`
	UserIDs    types.Set    `tfsdk:"user_ids"`
	GroupIDs   types.Set    `tfsdk:"group_ids"`
	DeviceIDs  types.Set    `tfsdk:"device_ids"`
	Platforms  types.Set    `tfsdk:"platforms"`
}

var (
	ForwardingSourceMatchingAttrTypes = map[string]attr.Type{"session_type": types.StringType}
	LocalNetworkSettingsAttrTypes     = map[string]attr.Type{
		"allow_local_subnet_bypass": types.BoolType,
	}
	AssignmentRulesAttrTypes = map[string]attr.Type{
		"membership": types.StringType,
		"user_ids":   types.SetType{ElemType: types.StringType},
		"group_ids":  types.SetType{ElemType: types.StringType},
		"device_ids": types.SetType{ElemType: types.StringType},
		"platforms":  types.SetType{ElemType: types.StringType},
	}
)

// ForwardingProfileSettings holds the writable v2 settings shared by custom and
// default profiles.
type ForwardingProfileSettings struct {
	State                    types.String
	ClientFallbackAction     types.String
	ForwardingSourceMatching types.Object
	LocalNetworkSettings     types.Object
	AssignmentRules          types.Object
}

func known(v attr.Value) bool { return !v.IsNull() && !v.IsUnknown() }

// SettingsPatch adds the settings that are configured and differ from state to
// body. A nil state treats every configured value as changed. Graph replaces
// assignmentRules as a whole, so it is always sent complete.
func SettingsPatch(
	ctx context.Context,
	body map[string]any,
	plan ForwardingProfileSettings,
	state *ForwardingProfileSettings,
) diag.Diagnostics {
	var diags diag.Diagnostics
	changed := func(p, s attr.Value) bool { return known(p) && (state == nil || !p.Equal(s)) }
	var s ForwardingProfileSettings
	if state != nil {
		s = *state
	}
	if changed(plan.State, s.State) {
		body["state"] = plan.State.ValueString()
	}
	if changed(plan.ClientFallbackAction, s.ClientFallbackAction) {
		body["clientFallbackAction"] = plan.ClientFallbackAction.ValueString()
	}
	if changed(plan.ForwardingSourceMatching, s.ForwardingSourceMatching) {
		var m ForwardingSourceMatchingModel
		diags.Append(plan.ForwardingSourceMatching.As(ctx, &m, basetypes.ObjectAsOptions{})...)
		body["forwardingSourceMatching"] = map[string]any{
			"sessionType": m.SessionType.ValueString(),
		}
	}
	if changed(plan.LocalNetworkSettings, s.LocalNetworkSettings) {
		var m LocalNetworkSettingsModel
		diags.Append(plan.LocalNetworkSettings.As(ctx, &m, basetypes.ObjectAsOptions{})...)
		body["localNetworkSettings"] = map[string]any{
			"allowLocalSubnetBypass": m.AllowLocalSubnetBypass.ValueBool(),
		}
	}
	if changed(plan.AssignmentRules, s.AssignmentRules) {
		rules, d := assignmentRulesBody(ctx, plan.AssignmentRules)
		diags.Append(d...)
		body["assignmentRules"] = rules
	}
	return diags
}

func assignmentRulesBody(ctx context.Context, obj types.Object) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics
	var m AssignmentRulesModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}
	ids := func(set types.Set) []map[string]string {
		out := []map[string]string{}
		if !known(set) {
			return out
		}
		var values []string
		diags.Append(set.ElementsAs(ctx, &values, false)...)
		slices.Sort(values)
		for _, v := range values {
			out = append(out, map[string]string{"id": v})
		}
		return out
	}
	assignments := map[string]any{}
	switch m.Membership.ValueString() {
	case MembershipAll:
		assignments["@odata.type"] = odataAssignmentsAll
	case MembershipEnumerated:
		// Graph rejects enumerated membership unless all three lists are present.
		assignments["@odata.type"] = odataAssignmentsEnumerated
		assignments["users"] = ids(m.UserIDs)
		assignments["groups"] = ids(m.GroupIDs)
		assignments["devices"] = ids(m.DeviceIDs)
	default:
		assignments["@odata.type"] = odataAssignmentsNone
	}
	var platforms []string
	if known(m.Platforms) {
		diags.Append(m.Platforms.ElementsAs(ctx, &platforms, false)...)
	}
	return map[string]any{
		"assignments":           assignments,
		"devicePropertyFilters": map[string]any{"platforms": joinPlatforms(platforms)},
	}, diags
}

// joinPlatforms orders platforms as Graph returns them so request logs are stable.
func joinPlatforms(platforms []string) string {
	out := make([]string, 0, len(platforms))
	for _, p := range ForwardingProfilePlatforms {
		if slices.Contains(platforms, p) {
			out = append(out, p)
		}
	}
	return strings.Join(out, ",")
}

// SettingsFromWire maps the writable v2 settings returned by Graph.
func SettingsFromWire(
	ctx context.Context,
	w *ForwardingProfileWire,
) (ForwardingProfileSettings, diag.Diagnostics) {
	var diags diag.Diagnostics
	out := ForwardingProfileSettings{
		State:                    types.StringPointerValue(w.State),
		ClientFallbackAction:     types.StringPointerValue(w.ClientFallbackAction),
		ForwardingSourceMatching: types.ObjectNull(ForwardingSourceMatchingAttrTypes),
		LocalNetworkSettings:     types.ObjectNull(LocalNetworkSettingsAttrTypes),
		AssignmentRules:          types.ObjectNull(AssignmentRulesAttrTypes),
	}
	if w.ForwardingSourceMatching != nil && w.ForwardingSourceMatching.SessionType != nil {
		v, d := types.ObjectValueFrom(
			ctx,
			ForwardingSourceMatchingAttrTypes,
			ForwardingSourceMatchingModel{
				SessionType: types.StringPointerValue(w.ForwardingSourceMatching.SessionType),
			},
		)
		diags.Append(d...)
		out.ForwardingSourceMatching = v
	}
	if w.LocalNetworkSettings != nil && w.LocalNetworkSettings.AllowLocalSubnetBypass != nil {
		v, d := types.ObjectValueFrom(ctx, LocalNetworkSettingsAttrTypes, LocalNetworkSettingsModel{
			AllowLocalSubnetBypass: types.BoolPointerValue(
				w.LocalNetworkSettings.AllowLocalSubnetBypass,
			),
		})
		diags.Append(d...)
		out.LocalNetworkSettings = v
	}
	if w.AssignmentRules != nil && w.AssignmentRules.Assignments != nil {
		v, d := assignmentRulesFromWire(ctx, w.AssignmentRules)
		diags.Append(d...)
		out.AssignmentRules = v
	}
	return out, diags
}

func assignmentRulesFromWire(
	ctx context.Context,
	w *ForwardingProfileAssignmentsRW,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	nullSet := types.SetNull(types.StringType)
	m := AssignmentRulesModel{UserIDs: nullSet, GroupIDs: nullSet, DeviceIDs: nullSet}
	// Graph returns empty lists for unused member kinds; the schema keeps them null.
	setOf := func(ids []ForwardingProfileDirectoryID) types.Set {
		if len(ids) == 0 {
			return nullSet
		}
		values := make([]attr.Value, 0, len(ids))
		for _, id := range ids {
			values = append(values, types.StringValue(id.ID))
		}
		v, d := types.SetValue(types.StringType, values)
		diags.Append(d...)
		return v
	}
	switch w.Assignments.ODataType {
	case odataAssignmentsAll:
		m.Membership = types.StringValue(MembershipAll)
	case odataAssignmentsEnumerated:
		m.Membership = types.StringValue(MembershipEnumerated)
		m.UserIDs = setOf(w.Assignments.Users)
		m.GroupIDs = setOf(w.Assignments.Groups)
		m.DeviceIDs = setOf(w.Assignments.Devices)
	case odataAssignmentsNone:
		m.Membership = types.StringValue(MembershipNone)
	default:
		diags.AddError(
			"Unsupported forwarding profile assignment type",
			fmt.Sprintf(
				"Graph returned assignment type %q, which this provider version does not support.",
				w.Assignments.ODataType,
			),
		)
		return types.ObjectNull(AssignmentRulesAttrTypes), diags
	}
	platforms := []attr.Value{}
	if w.DevicePropertyFilters != nil && w.DevicePropertyFilters.Platforms != nil {
		for p := range strings.SplitSeq(*w.DevicePropertyFilters.Platforms, ",") {
			p = strings.TrimSpace(p)
			if p == "" || strings.EqualFold(p, platformsNone) {
				continue
			}
			platforms = append(platforms, types.StringValue(canonicalPlatform(p)))
		}
	}
	v, d := types.SetValue(types.StringType, platforms)
	diags.Append(d...)
	m.Platforms = v
	obj, d := types.ObjectValueFrom(ctx, AssignmentRulesAttrTypes, m)
	diags.Append(d...)
	return obj, diags
}

func canonicalPlatform(p string) string {
	for _, c := range ForwardingProfilePlatforms {
		if strings.EqualFold(c, p) {
			return c
		}
	}
	return p
}
