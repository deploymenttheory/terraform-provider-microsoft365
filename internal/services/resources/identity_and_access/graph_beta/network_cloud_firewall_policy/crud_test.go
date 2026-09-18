package graphBetaNetworkCloudFirewallPolicy_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	cloud "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cloud_firewall_policy"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceNetworkCloudFirewallPolicy_02_CreateReadbackRetainsID(t *testing.T) {
	for _, body := range []string{`{"error":{"code":"NotFound","message":"readback not found"}}`, `{"id":"created","name":"incomplete"}`} {
		t.Run(body, func(t *testing.T) {
			ctx := context.Background()
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			clients := mocks.NewMocks()
			clients.AuthMocks.RegisterMocks()
			r := cloud.NewNetworkCloudFirewallPolicyResource().(*cloud.NetworkCloudFirewallPolicyResource)
			configured := resource.ConfigureResponse{}
			r.Configure(ctx, resource.ConfigureRequest{ProviderData: clients.Clients}, &configured)
			require.False(t, configured.Diagnostics.HasError())
			schema := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schema)
			plan := tfsdk.Plan{Schema: schema.Schema, Raw: tftypes.NewValue(schema.Schema.Type().TerraformType(ctx), nil)}
			for key, value := range map[string]any{"id": types.StringUnknown(), "name": "probe", "default_action": "allow", "version": types.StringUnknown(), "last_modified_date_time": types.StringUnknown()} {
				require.False(t, plan.SetAttribute(ctx, path.Root(key), value).HasError())
			}
			post, err := httpmock.NewJsonResponder(201, map[string]any{"id": "created"})
			require.NoError(t, err)
			httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkAccess/cloudFirewallPolicies", post)
			status := 200
			if body[2:7] == "error" {
				status = 404
			}
			var value map[string]any
			require.NoError(t, json.Unmarshal([]byte(body), &value))
			get, err := httpmock.NewJsonResponder(status, value)
			require.NoError(t, err)
			httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/networkAccess/cloudFirewallPolicies/created", get)
			response := resource.CreateResponse{State: tfsdk.State{Schema: schema.Schema}}
			r.Create(ctx, resource.CreateRequest{Plan: plan}, &response)
			require.True(t, response.Diagnostics.HasError(), response.Diagnostics)
			var id string
			require.False(t, response.State.GetAttribute(ctx, path.Root("id"), &id).HasError())
			require.Equal(t, "created", id)
		})
	}
}
