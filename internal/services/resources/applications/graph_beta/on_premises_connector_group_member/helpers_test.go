package graphBetaApplicationsOnPremisesConnectorGroupMember

import (
	"testing"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func newConnectorGroup(id string, isDefault bool) graphmodels.ConnectorGroupable {
	connectorGroup := graphmodels.NewConnectorGroup()
	connectorGroup.SetId(&id)
	connectorGroup.SetIsDefault(&isDefault)
	return connectorGroup
}

func TestUnitSelectDefaultConnectorGroup_PicksDefault(t *testing.T) {
	groups := []graphmodels.ConnectorGroupable{
		newConnectorGroup("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", false),
		nil,
		newConnectorGroup("dddddddd-dddd-dddd-dddd-dddddddddddd", true),
	}

	if got := connectorGroupIDOf(selectDefaultConnectorGroup(groups)); got != "dddddddd-dddd-dddd-dddd-dddddddddddd" {
		t.Fatalf("default connector group = %q, expected dddddddd-dddd-dddd-dddd-dddddddddddd", got)
	}
}

func TestUnitSelectDefaultConnectorGroup_NoDefault(t *testing.T) {
	groups := []graphmodels.ConnectorGroupable{
		newConnectorGroup("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", false),
		graphmodels.NewConnectorGroup(),
	}

	if got := selectDefaultConnectorGroup(groups); got != nil {
		t.Fatalf("expected no default connector group, got %q", connectorGroupIDOf(got))
	}
}

func TestUnitSameID_IsCaseInsensitiveAndRejectsEmpty(t *testing.T) {
	if !sameID("AAAAAAAA-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "aaaaaaaa-AAAA-aaaa-aaaa-aaaaaaaaaaaa") {
		t.Fatal("expected IDs differing only in case to match")
	}
	if sameID("", "") {
		t.Fatal("expected empty IDs not to match")
	}
}
