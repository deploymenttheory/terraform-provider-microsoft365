package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"os"
	"testing"

	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestConnectorGroupResponseParsesGraphResponseFixtures(t *testing.T) {
	tests := []struct {
		name           string
		fixture        string
		expectedName   string
		expectedRegion string
	}{
		{
			name:           "create",
			fixture:        "tests/responses/validate_create/post_connector_group_success.json",
			expectedName:   "unit-test-connector-group",
			expectedRegion: "japan",
		},
		{
			name:           "read",
			fixture:        "tests/responses/validate_get/get_connector_group_success.json",
			expectedName:   "unit-test-connector-group",
			expectedRegion: "nam",
		},
		{
			name:           "update",
			fixture:        "tests/responses/validate_update/patch_connector_group_success.json",
			expectedName:   "unit-test-connector-group-renamed",
			expectedRegion: "eur",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseJSON, err := os.ReadFile(tt.fixture)
			if err != nil {
				t.Fatalf("failed to read response fixture: %v", err)
			}

			parseNode, err := jsonserialization.NewJsonParseNodeFactory().GetRootParseNode("application/json", responseJSON)
			if err != nil {
				t.Fatalf("GetRootParseNode returned error: %v", err)
			}

			parsed, err := parseNode.GetObjectValue(createConnectorGroupResponseFromDiscriminatorValue)
			if err != nil {
				t.Fatalf("GetObjectValue returned error: %v", err)
			}

			response, ok := parsed.(*connectorGroupResponse)
			if !ok {
				t.Fatalf("parsed response is %T, expected *connectorGroupResponse", parsed)
			}

			if response.id == nil || *response.id != "00000000-0000-0000-0000-000000000000" {
				t.Fatalf("id = %#v, expected fixture id", response.id)
			}
			if response.name == nil || *response.name != tt.expectedName {
				t.Fatalf("name = %#v, expected %s", response.name, tt.expectedName)
			}
			if response.connectorGroupType == nil || *response.connectorGroupType != "applicationProxy" {
				t.Fatalf("connectorGroupType = %#v, expected applicationProxy", response.connectorGroupType)
			}
			if response.isDefault == nil || *response.isDefault {
				t.Fatalf("isDefault = %#v, expected false", response.isDefault)
			}
			if response.region == nil || *response.region != tt.expectedRegion {
				t.Fatalf("region = %#v, expected %s", response.region, tt.expectedRegion)
			}
		})
	}
}
