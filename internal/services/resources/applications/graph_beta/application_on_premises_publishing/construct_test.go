package graphBetaApplicationsOnPremisesPublishing

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUnit_ConstructOnPremisesPublishingPatchPayload(t *testing.T) {
	payload := constructOnPremisesPublishingPatchPayload(OnPremisesPublishingResourceModel{
		ApplicationType:               types.StringValue("nonwebapp"),
		ExternalUrl:                   types.StringUnknown(),
		InternalUrl:                   types.StringNull(),
		IsAccessibleViaZTNAClient:     types.BoolValue(true),
		IsDnsResolutionEnabled:        types.BoolValue(false),
		IsOnPremPublishingEnabled:     types.BoolValue(true),
		TrafficRoutingMethod:          types.StringValue("none"),
		WafProvider:                   types.StringNull(),
		IsHttpOnlyCookieEnabled:       types.BoolUnknown(),
		IsPersistentCookieEnabled:     types.BoolNull(),
		IsTranslateLinksInBodyEnabled: types.BoolValue(false),
	})

	if _, ok := payload["@odata.type"]; ok {
		t.Fatal("payload must not include top-level @odata.type")
	}

	onPremisesPublishing, ok := payload["onPremisesPublishing"].(map[string]any)
	if !ok {
		t.Fatalf("expected onPremisesPublishing object, got %T", payload["onPremisesPublishing"])
	}

	expected := map[string]any{
		"applicationType":               "nonwebapp",
		"isAccessibleViaZTNAClient":     true,
		"isDnsResolutionEnabled":        false,
		"isOnPremPublishingEnabled":     true,
		"trafficRoutingMethod":          "none",
		"isTranslateLinksInBodyEnabled": false,
	}

	for key, expectedValue := range expected {
		if actualValue, ok := onPremisesPublishing[key]; !ok || actualValue != expectedValue {
			t.Fatalf("expected %s to be %v, got %v", key, expectedValue, actualValue)
		}
	}

	for _, omittedKey := range []string{"externalUrl", "internalUrl", "wafProvider", "isHttpOnlyCookieEnabled", "isPersistentCookieEnabled"} {
		if _, ok := onPremisesPublishing[omittedKey]; ok {
			t.Fatalf("expected %s to be omitted", omittedKey)
		}
	}
}
