package graphBetaApplicationsOnPremisesConnectorGroupAssignment

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func TestMapRemoteResourceStateToTerraformSetsCompositeIDAndName(t *testing.T) {
	id := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	name := "Unit Test Connector Group"
	remote := graphmodels.NewConnectorGroup()
	remote.SetId(&id)
	remote.SetName(&name)

	data := &OnPremisesConnectorGroupAssignmentResourceModel{
		ApplicationID:    types.StringValue("11111111-1111-1111-1111-111111111111"),
		ConnectorGroupID: types.StringValue(id),
	}

	MapRemoteResourceStateToTerraform(context.Background(), data, remote)

	expectedID := "11111111-1111-1111-1111-111111111111/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	if data.ID.ValueString() != expectedID {
		t.Fatalf("ID = %q, expected %q", data.ID.ValueString(), expectedID)
	}
	if data.ConnectorGroupName.ValueString() != name {
		t.Fatalf("ConnectorGroupName = %q, expected %q", data.ConnectorGroupName.ValueString(), name)
	}
}
