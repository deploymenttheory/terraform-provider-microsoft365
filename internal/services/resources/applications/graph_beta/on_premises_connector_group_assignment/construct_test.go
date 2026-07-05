package graphBetaApplicationsOnPremisesConnectorGroupAssignment

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConstructResourceBuildsConnectorGroupReference(t *testing.T) {
	data := &OnPremisesConnectorGroupAssignmentResourceModel{
		ConnectorGroupID: types.StringValue("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
	}

	body, err := constructResource(context.Background(), data)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	expected := "https://graph.microsoft.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	if body.GetOdataId() == nil || *body.GetOdataId() != expected {
		t.Fatalf("@odata.id = %#v, expected %q", body.GetOdataId(), expected)
	}
}

func TestCompositeIDBuildsImportID(t *testing.T) {
	got := compositeID("11111111-1111-1111-1111-111111111111", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	expected := "11111111-1111-1111-1111-111111111111/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	if got != expected {
		t.Fatalf("compositeID = %q, expected %q", got, expected)
	}
}
