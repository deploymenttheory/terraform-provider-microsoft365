package graphBetaNetworkCloudFirewallPolicyRule_test

import (
	"context"
	"testing"

	cloud "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cloud_firewall_policy_rule"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	"github.com/stretchr/testify/require"
)

func testRemoteRule() *graph.CloudFirewallRule {
	r := graph.NewCloudFirewallRule()
	id, name := "00000000-0000-0000-0000-000000000001", "test"
	priority := int64(100)
	action := graph.BLOCK_CLOUDFIREWALLACTION
	status := graph.DISABLED_SECURITYRULESTATUS
	r.SetId(&id)
	r.SetName(&name)
	r.SetPriority(&priority)
	r.SetAction(&action)
	settings := graph.NewCloudFirewallRuleSettings()
	settings.SetStatus(&status)
	r.SetSettings(settings)
	mc := graph.NewCloudFirewallMatchingConditions()
	r.SetMatchingConditions(mc)
	return r
}
func TestUnitResourceNetworkCloudFirewallPolicyRule_02_ReadSets(t *testing.T) {
	ctx := context.Background()
	r := testRemoteRule()
	source := graph.NewCloudFirewallSourceMatching()
	source.SetPorts([]string{"443", "80", "443"})
	address := graph.NewCloudFirewallSourceIpAddress()
	address.SetValues([]string{"192.0.2.129/24", "192.0.2.1", "192.0.2.129/24"})
	source.SetAddresses([]graph.CloudFirewallSourceAddressable{address})
	r.GetMatchingConditions().SetSources(source)
	var state cloud.NetworkCloudFirewallPolicyRuleResourceModel
	require.NoError(t, cloud.MapRemoteStateToTerraform(ctx, &state, r))
	ports := state.Sources.Attributes()["ports"].(types.Set)
	require.Len(t, ports.Elements(), 2)
	first := state.Sources
	source.SetPorts([]string{"80", "443"})
	address.SetValues([]string{"192.0.2.1", "192.0.2.129/24"})
	require.NoError(t, cloud.MapRemoteStateToTerraform(ctx, &state, r))
	require.True(t, first.Equal(state.Sources))
	source.SetPorts([]string{})
	require.NoError(t, cloud.MapRemoteStateToTerraform(ctx, &state, r))
	require.False(t, state.Sources.Attributes()["ports"].IsNull())
}
func TestUnitResourceNetworkCloudFirewallPolicyRule_03_UnsupportedReadPreservesState(t *testing.T) {
	for _, tc := range []string{"unknown_status", "branches", "unknown_address", "missing_protocols"} {
		t.Run(tc, func(t *testing.T) {
			ctx := context.Background()
			r := testRemoteRule()
			state := cloud.NetworkCloudFirewallPolicyRuleResourceModel{Name: types.StringValue("prior")}
			switch tc {
			case "unknown_status":
				status := graph.UNKNOWNFUTUREVALUE_SECURITYRULESTATUS
				r.GetSettings().SetStatus(&status)
			case "branches":
				s := graph.NewCloudFirewallSourceMatching()
				s.SetAdditionalData(map[string]any{"branchIds": []string{"existing-branch"}})
				r.GetMatchingConditions().SetSources(s)
			case "unknown_address":
				s := graph.NewCloudFirewallSourceMatching()
				a := graph.NewCloudFirewallSourceIpAddress()
				typ := "#microsoft.graph.networkaccess.futureAddress"
				a.SetOdataType(&typ)
				s.SetAddresses([]graph.CloudFirewallSourceAddressable{a})
				r.GetMatchingConditions().SetSources(s)
			case "missing_protocols":
				r.GetMatchingConditions().SetDestinations(graph.NewCloudFirewallDestinationMatching())
			}
			require.Error(t, cloud.MapRemoteStateToTerraform(ctx, &state, r))
			require.Equal(t, "prior", state.Name.ValueString())
		})
	}
}
