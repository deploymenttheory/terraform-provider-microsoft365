package graphBetaIdentityAndAccessMobilityManagementPolicy

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func serializePolicy(t *testing.T, policy graphmodels.MobilityManagementPolicyable) string {
	t.Helper()
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", policy); err != nil {
		t.Fatalf("serialize policy: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized policy: %v", err)
	}
	return string(content)
}

func testModel(policyType string) *MobilityManagementPolicyResourceModel {
	return &MobilityManagementPolicyResourceModel{
		PolicyType:    types.StringValue(policyType),
		PolicyID:      types.StringValue("0000000a-0000-0000-c000-000000000000"),
		AppliesTo:     types.StringValue(appliesToAll),
		ComplianceUrl: types.StringValue("https://compliance.example.com"),
		DiscoveryUrl:  types.StringValue("https://discovery.example.com"),
		TermsOfUseUrl: types.StringNull(),
		IsMdmEnrollmentDuringRegistrationDisabled: types.BoolValue(true),
	}
}

func TestConstructSettingsRequest_CreateSendsConfiguredSettingsWithoutAppliesTo(t *testing.T) {
	body := constructSettingsRequest(context.Background(), testModel(policyTypeDevice), nil)
	if body == nil {
		t.Fatal("expected a settings request body")
	}
	if _, ok := body.(graphmodels.MobileDeviceManagementPolicyable); !ok {
		t.Fatalf("expected a device policy body, got %T", body)
	}
	json := serializePolicy(t, body)
	expectedFields := []string{
		`"complianceUrl":"https://compliance.example.com"`,
		`"discoveryUrl":"https://discovery.example.com"`,
		`"isMdmEnrollmentDuringRegistrationDisabled":true`,
	}
	for _, expected := range expectedFields {
		if !strings.Contains(json, expected) {
			t.Errorf("expected %s in %s", expected, json)
		}
	}
	for _, unexpected := range []string{"appliesTo", "termsOfUseUrl", "includedGroups"} {
		if strings.Contains(json, unexpected) {
			t.Errorf("did not expect %s in %s", unexpected, json)
		}
	}
}

func TestConstructSettingsRequest_UpdateSendsOnlyChangedSettings(t *testing.T) {
	state := testModel(policyTypeDevice)
	plan := testModel(policyTypeDevice)
	plan.DiscoveryUrl = types.StringValue("https://discovery-new.example.com")

	json := serializePolicy(t, constructSettingsRequest(context.Background(), plan, state))
	if !strings.Contains(json, `"discoveryUrl":"https://discovery-new.example.com"`) {
		t.Errorf("expected the changed discoveryUrl in %s", json)
	}
	for _, unexpected := range []string{"complianceUrl", "isMdmEnrollmentDuringRegistrationDisabled", "appliesTo"} {
		if strings.Contains(json, unexpected) {
			t.Errorf("did not expect unchanged %s in %s", unexpected, json)
		}
	}
}

func TestConstructSettingsRequest_NoChangesReturnsNil(t *testing.T) {
	if body := constructSettingsRequest(context.Background(), testModel(policyTypeDevice), testModel(policyTypeDevice)); body != nil {
		t.Fatalf("expected no request body, got %s", serializePolicy(t, body))
	}
}

func TestConstructSettingsRequest_AppPolicyIgnoresMdmSetting(t *testing.T) {
	body := constructSettingsRequest(context.Background(), testModel(policyTypeApp), nil)
	if _, ok := body.(graphmodels.MobileAppManagementPolicyable); !ok {
		t.Fatalf("expected an app policy body, got %T", body)
	}
	if json := serializePolicy(t, body); strings.Contains(json, "isMdmEnrollmentDuringRegistrationDisabled") {
		t.Errorf("did not expect the MDM-only setting in %s", json)
	}
}

func TestConstructAppliesToRequest_SendsOnlyAppliesTo(t *testing.T) {
	body, err := constructAppliesToRequest(context.Background(), policyTypeDevice, appliesToNone)
	if err != nil {
		t.Fatalf("constructAppliesToRequest: %v", err)
	}
	json := serializePolicy(t, body)
	if !strings.Contains(json, `"appliesTo":"none"`) {
		t.Errorf("expected appliesTo none in %s", json)
	}
	if strings.Count(json, ":") != 1 {
		t.Errorf("expected appliesTo to be the only property in %s", json)
	}
}

func TestConstructGroupReference(t *testing.T) {
	reference := constructGroupReference("11111111-1111-1111-1111-111111111111")
	if got := *reference.GetOdataId(); got != "https://graph.microsoft.com/beta/groups/11111111-1111-1111-1111-111111111111" {
		t.Errorf("unexpected @odata.id %q", got)
	}
}
