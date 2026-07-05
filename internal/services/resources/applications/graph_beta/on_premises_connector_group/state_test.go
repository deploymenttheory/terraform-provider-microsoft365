package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapRemoteResourceStateToTerraformPreservesRawRegion(t *testing.T) {
	id := "00000000-0000-0000-0000-000000000000"
	name := "unit-test-connector-group"
	connectorGroupType := "applicationProxy"
	isDefault := false
	region := "japan"
	remote := &connectorGroupResponse{
		id:                 &id,
		name:               &name,
		connectorGroupType: &connectorGroupType,
		isDefault:          &isDefault,
		region:             &region,
	}

	data := &OnPremisesConnectorGroupResourceModel{}
	MapRemoteResourceStateToTerraform(context.Background(), data, remote)

	if data.ID.ValueString() != id {
		t.Fatalf("ID = %q, expected %q", data.ID.ValueString(), id)
	}
	if data.Name.ValueString() != name {
		t.Fatalf("Name = %q, expected %q", data.Name.ValueString(), name)
	}
	if data.ConnectorGroupType.ValueString() != connectorGroupType {
		t.Fatalf("ConnectorGroupType = %q, expected %q", data.ConnectorGroupType.ValueString(), connectorGroupType)
	}
	if data.IsDefault != types.BoolValue(false) {
		t.Fatalf("IsDefault = %#v, expected false", data.IsDefault)
	}
	if data.Region.ValueString() != "japan" {
		t.Fatalf("Region = %q, expected japan", data.Region.ValueString())
	}
}
