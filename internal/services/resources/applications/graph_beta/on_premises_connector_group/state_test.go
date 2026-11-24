package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapRemoteStateToTerraformPreservesRawEnumValues(t *testing.T) {
	tests := []struct {
		name           string
		fixture        string
		expectedName   string
		expectedRegion string
	}{
		{name: "create", fixture: "tests/responses/validate_create/post_connector_group_success.json", expectedName: "unit-test-connector-group", expectedRegion: "japan"},
		{name: "read", fixture: "tests/responses/validate_get/get_connector_group_success.json", expectedName: "unit-test-connector-group", expectedRegion: "nam"},
		{name: "update", fixture: "tests/responses/validate_update/patch_connector_group_success.json", expectedName: "unit-test-connector-group-renamed", expectedRegion: "eur"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseJSON, err := os.ReadFile(tt.fixture)
			if err != nil {
				t.Fatalf("failed to read response fixture: %v", err)
			}

			var remote connectorGroupResponse
			if err := json.Unmarshal(responseJSON, &remote); err != nil {
				t.Fatalf("failed to decode response fixture: %v", err)
			}

			data := &OnPremisesConnectorGroupResourceModel{}
			MapRemoteStateToTerraform(context.Background(), data, &remote)

			if data.ID.ValueString() != "00000000-0000-0000-0000-000000000000" {
				t.Fatalf("ID = %q, expected fixture id", data.ID.ValueString())
			}
			if data.Name.ValueString() != tt.expectedName {
				t.Fatalf("Name = %q, expected %q", data.Name.ValueString(), tt.expectedName)
			}
			if data.ConnectorGroupType.ValueString() != "applicationProxy" {
				t.Fatalf("ConnectorGroupType = %q, expected applicationProxy", data.ConnectorGroupType.ValueString())
			}
			if data.IsDefault != types.BoolValue(false) {
				t.Fatalf("IsDefault = %#v, expected false", data.IsDefault)
			}
			if data.Region.ValueString() != tt.expectedRegion {
				t.Fatalf("Region = %q, expected %q", data.Region.ValueString(), tt.expectedRegion)
			}
		})
	}
}
