package graphBetaApplicationsOnPremisesConnectorGroupMember

import (
	"context"
	"testing"
)

func TestUnitConstructResource_BuildsConnectorGroupReference(t *testing.T) {
	body, err := constructResource(context.Background(), "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "https://graph.microsoft.com/beta")
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	expected := "https://graph.microsoft.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	if body.GetOdataId() == nil || *body.GetOdataId() != expected {
		t.Fatalf("@odata.id = %v, expected %q", body.GetOdataId(), expected)
	}
}

func TestUnitConstructResource_UsesConfiguredGraphBaseURL(t *testing.T) {
	body, err := constructResource(context.Background(), "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "https://graph.microsoft.us/beta/")
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	expected := "https://graph.microsoft.us/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	if body.GetOdataId() == nil || *body.GetOdataId() != expected {
		t.Fatalf("@odata.id = %v, expected %q", body.GetOdataId(), expected)
	}
}

func TestUnitConstructResource_DefaultsGraphBaseURL(t *testing.T) {
	body, err := constructResource(context.Background(), "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "")
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	expected := "https://graph.microsoft.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	if body.GetOdataId() == nil || *body.GetOdataId() != expected {
		t.Fatalf("@odata.id = %v, expected %q", body.GetOdataId(), expected)
	}
}

func TestUnitConstructResource_RequiresConnectorGroupID(t *testing.T) {
	if _, err := constructResource(context.Background(), "", "https://graph.microsoft.com/beta"); err == nil {
		t.Fatal("expected an error for an empty connector group ID")
	}
}
