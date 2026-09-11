package graphBetaNetworkCloudFirewallPolicy_test

import (
	"context"
	"os"
	"testing"

	cloud "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cloud_firewall_policy"
	kiotajson "github.com/microsoft/kiota-serialization-json-go"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceNetworkCloudFirewallPolicy_10_LiveResponseFixture(t *testing.T) {
	raw, err := os.ReadFile("tests/responses/validate_get/get_cloud_firewall_success.json")
	require.NoError(t, err)
	node, err := kiotajson.NewJsonParseNodeFactory().GetRootParseNode("application/json", raw)
	require.NoError(t, err)
	remote, err := node.GetObjectValue(graph.CreateCloudFirewallPolicyFromDiscriminatorValue)
	require.NoError(t, err)
	var state cloud.NetworkCloudFirewallPolicyResourceModel
	require.NoError(t, cloud.MapRemoteStateToTerraform(context.Background(), &state, remote.(graph.CloudFirewallPolicyable)))
	require.Equal(t, "00000000-0000-0000-0000-000000000001", state.ID.ValueString())
}
