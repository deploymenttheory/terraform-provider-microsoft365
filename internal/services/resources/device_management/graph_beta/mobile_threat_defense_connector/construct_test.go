package graphBetaMobileThreatDefenseConnector

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func serializeConnector(t *testing.T, body graphmodels.MobileThreatDefenseConnectorable) map[string]any {
	t.Helper()
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", body); err != nil {
		t.Fatalf("serialize connector: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized connector: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("parse serialized connector %s: %v", content, err)
	}
	return result
}

// The partner ID is the connector ID, and Graph needs it in the POST body to know which partner to create.
func TestUnitConstructRequestBodySendsPartnerIdAndToggles(t *testing.T) {
	data := &MobileThreatDefenseConnectorResourceModel{
		ID:                                     types.StringValue("8d0ed095-8191-4bd3-8a41-953b22d51ff7"),
		AndroidEnabled:                         types.BoolValue(true),
		IosEnabled:                             types.BoolValue(false),
		MacEnabled:                             types.BoolValue(true),
		WindowsEnabled:                         types.BoolValue(false),
		PartnerUnresponsivenessThresholdInDays: types.Int32Value(14),
		AllowPartnerToCollectIosApplicationMetadata:         types.BoolValue(true),
		AllowPartnerToCollectIosPersonalApplicationMetadata: types.BoolValue(false),
		AllowPartnerToCollectIosCertificateMetadata:         types.BoolValue(true),
		AllowPartnerToCollectIosPersonalCertificateMetadata: types.BoolValue(false),
		MicrosoftDefenderForEndpointAttachEnabled:           types.BoolValue(true),
	}

	body := serializeConnector(t, constructRequestBody(data))

	expected := map[string]any{
		"id":                                     "8d0ed095-8191-4bd3-8a41-953b22d51ff7",
		"androidEnabled":                         true,
		"iosEnabled":                             false,
		"macEnabled":                             true,
		"windowsEnabled":                         false,
		"partnerUnresponsivenessThresholdInDays": float64(14),
		"allowPartnerToCollectIOSApplicationMetadata":         true,
		"allowPartnerToCollectIOSPersonalApplicationMetadata": false,
		"allowPartnerToCollectIosCertificateMetadata":         true,
		"allowPartnerToCollectIosPersonalCertificateMetadata": false,
		"microsoftDefenderForEndpointAttachEnabled":           true,
	}
	for key, value := range expected {
		if body[key] != value {
			t.Errorf("expected %s = %v, got %v (body %v)", key, value, body[key], body)
		}
	}
}

// Graph derives the *DeviceBlockedOnMissingPartnerData fields and owns the partner* fields, so the body must not carry them.
func TestUnitConstructRequestBodyOmitsReadOnlyFields(t *testing.T) {
	data := &MobileThreatDefenseConnectorResourceModel{
		ID:                                       types.StringValue("8d0ed095-8191-4bd3-8a41-953b22d51ff7"),
		AndroidDeviceBlockedOnMissingPartnerData: types.BoolValue(true),
		IosDeviceBlockedOnMissingPartnerData:     types.BoolValue(true),
		MacDeviceBlockedOnMissingPartnerData:     types.BoolValue(true),
		WindowsDeviceBlockedOnMissingPartnerData: types.BoolValue(true),
		PartnerState:                             types.StringValue("enabled"),
		LastHeartbeatDateTime:                    types.StringValue("2026-11-15T03:34:40Z"),
		PartnerDisplayName:                       types.StringValue("Partner"),
		PartnerSupportsAndroidMam:                types.BoolValue(true),
		PartnerMinimumSupportedIosVersion:        types.StringValue("16.0"),
	}

	body := serializeConnector(t, constructRequestBody(data))

	for _, field := range []string{
		"androidDeviceBlockedOnMissingPartnerData", "iosDeviceBlockedOnMissingPartnerData", "macDeviceBlockedOnMissingPartnerData", "windowsDeviceBlockedOnMissingPartnerData",
		"partnerState", "lastHeartbeatDateTime", "partnerDisplayName", "partnerSupportsAndroidMam", "partnerMinimumSupportedIosVersion",
	} {
		if _, ok := body[field]; ok {
			t.Errorf("expected %s to be omitted, got %v", field, body)
		}
	}
}

// Unknown or null toggles are left out instead of being sent as null.
func TestUnitConstructRequestBodyOmitsUnsetValues(t *testing.T) {
	data := &MobileThreatDefenseConnectorResourceModel{
		ID:                                     types.StringValue("8d0ed095-8191-4bd3-8a41-953b22d51ff7"),
		AndroidEnabled:                         types.BoolNull(),
		IosEnabled:                             types.BoolUnknown(),
		PartnerUnresponsivenessThresholdInDays: types.Int32Null(),
	}

	body := serializeConnector(t, constructRequestBody(data))

	for _, field := range []string{"androidEnabled", "iosEnabled", "partnerUnresponsivenessThresholdInDays"} {
		if _, ok := body[field]; ok {
			t.Errorf("expected %s to be omitted, got %v", field, body)
		}
	}
}
