package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestValidPort(t *testing.T) {
	cases := map[string]bool{
		"443":         true,
		"0":           true,
		"65535":       true,
		"1000-2000":   true,
		"80-80":       true,
		"70000":       false,
		"2000-1000":   false,
		"1000-70000":  false,
		"abc":         false,
		"1000 - 2000": false,
		"":            false,
		"-1":          false,
		"443-":        false,
		"123456":      false,
		"00443":       false,
		"0080-0090":   false,
	}
	for value, want := range cases {
		if got := validPort(value); got != want {
			t.Errorf("validPort(%q) = %v, expected %v", value, got, want)
		}
	}
}

func TestPortValidatorSkipsUnknown(t *testing.T) {
	resp := &validator.StringResponse{}
	portValidator{}.ValidateString(context.Background(), validator.StringRequest{
		Path:        path.Root("ports"),
		ConfigValue: types.StringUnknown(),
	}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error for unknown port: %v", resp.Diagnostics)
	}
}

func TestDestinationsValidator(t *testing.T) {
	cases := map[string]struct {
		ruleType     string
		destinations []map[string]string
		wantErrs     []string
	}{
		"fqdn ok": {
			ruleType:     ruleTypeFQDN,
			destinations: []map[string]string{{"type": ruleTypeFQDN, "value": "example.com"}},
		},
		"ipv4 address ok": {
			ruleType:     ruleTypeIPAddress,
			destinations: []map[string]string{{"type": ruleTypeIPAddress, "value": "198.51.100.5"}},
		},
		"ipv4 subnet ok": {
			ruleType:     ruleTypeIPSubnet,
			destinations: []map[string]string{{"type": ruleTypeIPSubnet, "value": "198.51.100.0/24"}},
		},
		"ipv4 range ok": {
			ruleType:     ruleTypeIPRange,
			destinations: []map[string]string{{"type": ruleTypeIPRange, "begin_address": "198.51.100.10", "end_address": "198.51.100.20"}},
		},
		"type mismatch": {
			ruleType:     ruleTypeFQDN,
			destinations: []map[string]string{{"type": ruleTypeIPAddress, "value": "198.51.100.5"}},
			wantErrs:     []string{"Destination type does not match rule_type"},
		},
		"ipv6 address": {
			ruleType:     ruleTypeIPAddress,
			destinations: []map[string]string{{"type": ruleTypeIPAddress, "value": "2001:db8::5"}},
			wantErrs:     []string{"Invalid IP address"},
		},
		"ipv6 subnet": {
			ruleType:     ruleTypeIPSubnet,
			destinations: []map[string]string{{"type": ruleTypeIPSubnet, "value": "2001:db8::/48"}},
			wantErrs:     []string{"Invalid IP subnet"},
		},
		"ipv4-mapped ipv6 address": {
			ruleType:     ruleTypeIPAddress,
			destinations: []map[string]string{{"type": ruleTypeIPAddress, "value": "::ffff:198.51.100.5"}},
			wantErrs:     []string{"Invalid IP address"},
		},
		"subnet without prefix": {
			ruleType:     ruleTypeIPSubnet,
			destinations: []map[string]string{{"type": ruleTypeIPSubnet, "value": "198.51.100.0"}},
			wantErrs:     []string{"Invalid IP subnet"},
		},
		"reversed range": {
			ruleType:     ruleTypeIPRange,
			destinations: []map[string]string{{"type": ruleTypeIPRange, "begin_address": "198.51.100.20", "end_address": "198.51.100.10"}},
			wantErrs:     []string{"Invalid IP range"},
		},
		"range missing end and with value": {
			ruleType: ruleTypeIPRange,
			destinations: []map[string]string{
				{"type": ruleTypeIPRange, "value": "198.51.100.1", "begin_address": "198.51.100.10"},
			},
			wantErrs: []string{"Unexpected destination field", "Missing destination field"},
		},
		"fqdn missing value and with range fields": {
			ruleType: ruleTypeFQDN,
			destinations: []map[string]string{
				{"type": ruleTypeFQDN, "begin_address": "198.51.100.10"},
			},
			wantErrs: []string{"Missing destination field", "Unexpected destination field"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			resp := validateConfigForTest(t, &NetworkInternetAccessForwardingPolicyRuleResourceModel{
				RuleType:     types.StringValue(tc.ruleType),
				Destinations: destinationsForTest(t, tc.destinations),
			})
			var got []string
			for _, d := range resp.Diagnostics.Errors() {
				got = append(got, d.Summary())
			}
			if strings.Join(got, "|") != strings.Join(tc.wantErrs, "|") {
				t.Fatalf("errors = %v, expected %v", got, tc.wantErrs)
			}
		})
	}
}

func TestDestinationsValidatorSkipsUnknown(t *testing.T) {
	resp := validateConfigForTest(t, &NetworkInternetAccessForwardingPolicyRuleResourceModel{
		RuleType:     types.StringUnknown(),
		Destinations: destinationsForTest(t, []map[string]string{{"type": ruleTypeIPAddress, "value": "198.51.100.5"}}),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors with unknown rule_type: %v", resp.Diagnostics)
	}
}

func TestConstructResourceSendsClientFallbackActionOnCreateOnly(t *testing.T) {
	model := &NetworkInternetAccessForwardingPolicyRuleResourceModel{
		ID:                   types.StringValue("66f1c2ee-9f17-4681-a7bc-c324e7dff554"),
		Name:                 types.StringValue("rule"),
		Action:               types.StringValue("forward"),
		RuleType:             types.StringValue(ruleTypeFQDN),
		ClientFallbackAction: types.StringValue("block"),
		Ports:                stringSetForTest(t, []string{"1000-2000"}),
		Protocol:             types.StringValue("udp"),
		Destinations:         destinationsForTest(t, []map[string]string{{"type": ruleTypeFQDN, "value": "example.com"}}),
	}

	create := serializedRulePayload(t, context.Background(), model, false)
	if create["clientFallbackAction"] != "block" {
		t.Fatalf("create clientFallbackAction = %#v, expected block", create["clientFallbackAction"])
	}
	update := serializedRulePayload(t, context.Background(), model, true)
	if _, ok := update["clientFallbackAction"]; ok {
		t.Fatalf("update payload included clientFallbackAction, but Graph rejects updating it")
	}

	model.ClientFallbackAction = types.StringUnknown()
	if _, ok := serializedRulePayload(t, context.Background(), model, false)["clientFallbackAction"]; ok {
		t.Fatalf("create payload included clientFallbackAction although it was not configured")
	}
}

func validateConfigForTest(
	t *testing.T,
	model *NetworkInternetAccessForwardingPolicyRuleResourceModel,
) *resource.ValidateConfigResponse {
	t.Helper()
	ctx := context.Background()

	schemaResp := &resource.SchemaResponse{}
	(&NetworkInternetAccessForwardingPolicyRuleResource{}).Schema(ctx, resource.SchemaRequest{}, schemaResp)
	model.ID = types.StringNull()
	model.ForwardingPolicyID = types.StringNull()
	model.Name = types.StringNull()
	model.Action = types.StringNull()
	model.ClientFallbackAction = types.StringNull()
	model.Ports = types.SetNull(types.StringType)
	model.Protocol = types.StringNull()
	model.Timeouts = timeouts.Value{Object: types.ObjectNull(timeoutAttrTypes(ctx, schemaResp))}

	plan := tfsdk.Plan{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
	}
	if diags := plan.Set(ctx, model); diags.HasError() {
		t.Fatalf("failed to build config: %v", diags)
	}

	resp := &resource.ValidateConfigResponse{}
	destinationsValidator{}.ValidateResource(ctx, resource.ValidateConfigRequest{
		Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: plan.Raw},
	}, resp)

	return resp
}

func timeoutAttrTypes(_ context.Context, schemaResp *resource.SchemaResponse) map[string]attr.Type {
	return schemaResp.Schema.Attributes["timeouts"].GetType().(timeouts.Type).AttrTypes
}

func TestDestinationsValidatorSkipsUnknownValues(t *testing.T) {
	objectType := destinationObjectType()
	unknownFields := types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
		"type":          types.StringValue(ruleTypeIPRange),
		"value":         types.StringNull(),
		"begin_address": types.StringUnknown(),
		"end_address":   types.StringUnknown(),
	})

	for name, destinations := range map[string]types.List{
		"unknown list":    types.ListUnknown(objectType),
		"unknown element": types.ListValueMust(objectType, []attr.Value{types.ObjectUnknown(objectType.AttrTypes)}),
		"unknown fields":  types.ListValueMust(objectType, []attr.Value{unknownFields}),
	} {
		t.Run(name, func(t *testing.T) {
			resp := validateConfigForTest(t, &NetworkInternetAccessForwardingPolicyRuleResourceModel{
				RuleType:     types.StringValue(ruleTypeIPRange),
				Destinations: destinations,
			})
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected errors: %v", resp.Diagnostics)
			}
		})
	}
}

// Pins the client_fallback_action modifiers. On update, an unknown config value
// must stay unknown and force replacement, and a removed value must keep state.
// On create, nothing is replaced and an unset value is left for Graph to fill.
func TestClientFallbackActionPlanModifiers(t *testing.T) {
	ctx := context.Background()
	schemaResp := &resource.SchemaResponse{}
	(&NetworkInternetAccessForwardingPolicyRuleResource{}).Schema(ctx, resource.SchemaRequest{}, schemaResp)
	attribute := schemaResp.Schema.Attributes["client_fallback_action"].(schema.StringAttribute)
	raw := nonNullObjectForTest(ctx, schemaResp)

	cases := map[string]struct {
		config      types.String
		wantPlan    types.String
		wantReplace bool
		create      bool
	}{
		"unknown config forces replacement": {types.StringUnknown(), types.StringUnknown(), true, false},
		"removed keeps state":               {types.StringNull(), types.StringValue("block"), false, false},
		"same value":                        {types.StringValue("block"), types.StringValue("block"), false, false},
		"new value forces replacement":      {types.StringValue("bypass"), types.StringValue("bypass"), true, false},
		"create leaves unset value unknown": {types.StringNull(), types.StringUnknown(), false, true},
		"create with unknown value":         {types.StringUnknown(), types.StringUnknown(), false, true},
		"create with a value":               {types.StringValue("block"), types.StringValue("block"), false, true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			planValue := tc.config
			if tc.config.IsNull() {
				planValue = types.StringUnknown() // Optional+Computed: core plans unknown.
			}
			state, stateRaw := types.StringValue("block"), raw
			if tc.create {
				state = types.StringNull()
				stateRaw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)
			}
			replace := false
			for _, m := range attribute.PlanModifiers {
				resp := &planmodifier.StringResponse{PlanValue: planValue}
				m.PlanModifyString(ctx, planmodifier.StringRequest{
					Path:        path.Root("client_fallback_action"),
					ConfigValue: tc.config,
					StateValue:  state,
					PlanValue:   planValue,
					Config:      tfsdk.Config{Schema: schemaResp.Schema, Raw: raw},
					Plan:        tfsdk.Plan{Schema: schemaResp.Schema, Raw: raw},
					State:       tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw},
				}, resp)
				if resp.Diagnostics.HasError() {
					t.Fatalf("modifier error: %v", resp.Diagnostics)
				}
				planValue = resp.PlanValue
				replace = replace || resp.RequiresReplace
			}
			if !planValue.Equal(tc.wantPlan) || replace != tc.wantReplace {
				t.Fatalf("plan = %s, replace = %v; expected %s, %v", planValue, replace, tc.wantPlan, tc.wantReplace)
			}
		})
	}
}

func nonNullObjectForTest(ctx context.Context, schemaResp *resource.SchemaResponse) tftypes.Value {
	objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	values := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attrType := range objectType.AttributeTypes {
		values[name] = tftypes.NewValue(attrType, nil)
	}

	return tftypes.NewValue(objectType, values)
}
