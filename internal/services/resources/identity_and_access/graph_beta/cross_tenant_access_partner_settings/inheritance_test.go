package graphBetaCrossTenantAccessPartnerSettings

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func inheritedModel() CrossTenantAccessPartnerSettingsResourceModel {
	return CrossTenantAccessPartnerSettingsResourceModel{
		ID:                           types.StringValue("12345678-1234-1234-1234-123456789012"),
		TenantID:                     types.StringValue("12345678-1234-1234-1234-123456789012"),
		B2bCollaborationInbound:      types.ObjectNull(getB2BSettingAttrTypes()),
		B2bCollaborationOutbound:     types.ObjectNull(getB2BSettingAttrTypes()),
		B2bDirectConnectInbound:      types.ObjectNull(getB2BSettingAttrTypes()),
		B2bDirectConnectOutbound:     types.ObjectNull(getB2BSettingAttrTypes()),
		InboundTrust:                 types.ObjectNull(getInboundTrustAttrTypes()),
		AutomaticUserConsentSettings: types.ObjectNull(getAutomaticUserConsentAttrTypes()),
		TenantRestrictions:           types.ObjectNull(getTenantRestrictionsAttrTypes()),
		EnforceDefaultInheritance:    types.BoolValue(true),
		Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType, "read": types.StringType,
			"update": types.StringType, "delete": types.StringType,
		})},
	}
}

func TestUnitResourceCrossTenantAccessPartnerSettings_RefreshInheritanceFlag(t *testing.T) {
	for _, flag := range []types.Bool{types.BoolNull(), types.BoolValue(false), types.BoolValue(true)} {
		data := inheritedModel()
		data.EnforceDefaultInheritance = flag
		MapRemoteResourceStateToTerraform(
			t.Context(),
			&data,
			graphmodels.NewCrossTenantAccessPolicyConfigurationPartner(),
		)
		if data.EnforceDefaultInheritance.IsNull() ||
			data.EnforceDefaultInheritance.ValueBool() != flag.ValueBool() {
			t.Fatal(
				"refresh must initialise the default and preserve explicitly enabled enforcement",
			)
		}
	}
}

func TestUnitResourceCrossTenantAccessPartnerSettings_InheritancePayload(t *testing.T) {
	for _, enforce := range []bool{false, true} {
		model := inheritedModel()
		model.EnforceDefaultInheritance = types.BoolValue(enforce)
		body, err := constructResource(t.Context(), &model)
		if err != nil {
			t.Fatal(err)
		}
		writer := jsonserialization.NewJsonSerializationWriter()
		if err := writer.WriteObjectValue("", body); err != nil {
			t.Fatal(err)
		}
		content, err := writer.GetSerializedContent()
		if err != nil {
			t.Fatal(err)
		}
		var payload map[string]json.RawMessage
		if err := json.Unmarshal(content, &payload); err != nil {
			t.Fatal(err)
		}
		for _, setting := range inheritablePartnerSettings {
			value, exists := payload[setting.graphName]
			if enforce && (!exists || string(value) != "null") {
				t.Errorf("%s must serialize as explicit null: %s", setting.graphName, content)
			} else if !enforce && exists {
				t.Errorf("legacy request must omit %s: %s", setting.graphName, content)
			}
		}
	}
}

func TestUnitResourceCrossTenantAccessPartnerSettings_InheritanceConsistency(t *testing.T) {
	ctx := t.Context()
	r := &CrossTenantAccessPartnerSettingsResource{}
	var schemaResponse resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	expected := inheritedModel()
	actual := expected
	state := tfsdk.State{Schema: schemaResponse.Schema}
	// A null expected value must reject a concrete stale override while the
	// write propagates across replicas.
	actual.InboundTrust = types.ObjectValueMust(getInboundTrustAttrTypes(), map[string]attr.Value{
		"is_mfa_accepted":                           types.BoolValue(true),
		"is_compliant_device_accepted":              types.BoolValue(false),
		"is_hybrid_azure_ad_joined_device_accepted": types.BoolValue(false),
	})
	if diags := state.Set(ctx, &actual); diags.HasError() {
		t.Fatal(diags)
	}
	if crossTenantAccessPartnerSettingsConsistencyPredicate(&expected)(ctx, state) {
		t.Fatal("stale override must not satisfy enforced inheritance")
	}
	expected.EnforceDefaultInheritance = types.BoolValue(false)
	if !crossTenantAccessPartnerSettingsConsistencyPredicate(&expected)(ctx, state) {
		t.Fatal("legacy null settings should accept the remote override")
	}
}

func TestUnitResourceCrossTenantAccessPartnerSettings_InheritancePlanning(t *testing.T) {
	ctx := t.Context()
	r := &CrossTenantAccessPartnerSettingsResource{}
	var schemaResponse resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	for _, flag := range []types.Bool{types.BoolValue(false), types.BoolValue(true), types.BoolUnknown()} {
		configured := inheritedModel()
		configured.EnforceDefaultInheritance = flag
		configData := tfsdk.Plan{Schema: schemaResponse.Schema}
		if diags := configData.Set(ctx, &configured); diags.HasError() {
			t.Fatal(diags)
		}
		planned := configured
		planned.InboundTrust = types.ObjectValueMust(
			getInboundTrustAttrTypes(),
			map[string]attr.Value{
				"is_mfa_accepted":                           types.BoolValue(true),
				"is_compliant_device_accepted":              types.BoolValue(false),
				"is_hybrid_azure_ad_joined_device_accepted": types.BoolValue(false),
			},
		)
		plan := tfsdk.Plan{Schema: schemaResponse.Schema}
		if diags := plan.Set(ctx, &planned); diags.HasError() {
			t.Fatal(diags)
		}
		response := resource.ModifyPlanResponse{Plan: plan}
		r.ModifyPlan(ctx, resource.ModifyPlanRequest{
			Plan: plan, Config: tfsdk.Config{Schema: schemaResponse.Schema, Raw: configData.Raw},
		}, &response)
		if response.Diagnostics.HasError() {
			t.Fatal(response.Diagnostics)
		}
		var actual types.Object
		if diags := response.Plan.GetAttribute(
			ctx,
			path.Root("inbound_trust"),
			&actual,
		); diags.HasError() {
			t.Fatal(diags)
		}
		switch {
		case flag.IsUnknown():
			if !actual.IsUnknown() {
				t.Fatal("unresolved enforcement must leave the inherited setting unknown")
			}
		case flag.ValueBool():
			if !actual.IsNull() {
				t.Fatal("enforced inheritance must plan an explicit null")
			}
		default:
			if !actual.Equal(planned.InboundTrust) {
				t.Fatal("legacy planning must preserve the remote computed value")
			}
		}
	}
}
