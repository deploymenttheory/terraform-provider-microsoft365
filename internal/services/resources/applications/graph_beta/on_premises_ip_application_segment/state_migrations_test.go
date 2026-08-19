package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUpgradeOnPremisesIpApplicationSegmentStateV0toV1(t *testing.T) {
	t.Parallel()

	for _, protocol := range []string{"tcp", "udp"} {
		t.Run(protocol, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			priorModel := onPremisesIpApplicationSegmentResourceModelV0{
				ID:                  types.StringValue("segment-id"),
				ApplicationObjectID: types.StringValue("application-id"),
				DestinationHost:     types.StringValue("192.168.1.100"),
				DestinationType:     types.StringValue("ipAddress"),
				Ports: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("443-443"),
				}),
				Protocol: types.StringValue(protocol),
				Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
					"create": types.StringType,
					"read":   types.StringType,
					"update": types.StringType,
					"delete": types.StringType,
				})},
			}

			resourceUnderTest := &OnPremisesIpApplicationSegmentResource{}
			upgrader, ok := resourceUnderTest.UpgradeState(ctx)[0]
			if !ok {
				t.Fatal("schema version 0 state upgrader is not registered")
			}

			priorState := tfsdk.State{Schema: upgrader.PriorSchema}
			if diags := priorState.Set(ctx, &priorModel); diags.HasError() {
				t.Fatalf("failed to create v0 state: %v", diags)
			}

			var schemaResponse resource.SchemaResponse
			resourceUnderTest.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
			if schemaResponse.Diagnostics.HasError() {
				t.Fatalf("failed to obtain v1 schema: %v", schemaResponse.Diagnostics)
			}

			request := resource.UpgradeStateRequest{State: &priorState}
			response := resource.UpgradeStateResponse{
				State: tfsdk.State{Schema: schemaResponse.Schema},
			}
			upgrader.StateUpgrader(ctx, request, &response)
			if response.Diagnostics.HasError() {
				t.Fatalf("state migration returned errors: %v", response.Diagnostics)
			}

			var upgradedModel OnPremisesIpApplicationSegmentResourceModel
			if diags := response.State.Get(ctx, &upgradedModel); diags.HasError() {
				t.Fatalf("failed to read v1 state: %v", diags)
			}

			var upgradedProtocols []string
			if diags := upgradedModel.Protocol.ElementsAs(ctx, &upgradedProtocols, false); diags.HasError() {
				t.Fatalf("failed to read upgraded protocols: %v", diags)
			}
			if len(upgradedProtocols) != 1 || upgradedProtocols[0] != protocol {
				t.Fatalf("protocol = %#v, expected [%s]", upgradedProtocols, protocol)
			}

			if upgradedModel.ID.ValueString() != priorModel.ID.ValueString() ||
				upgradedModel.ApplicationObjectID.ValueString() != priorModel.ApplicationObjectID.ValueString() ||
				upgradedModel.DestinationHost.ValueString() != priorModel.DestinationHost.ValueString() ||
				upgradedModel.DestinationType.ValueString() != priorModel.DestinationType.ValueString() ||
				!upgradedModel.Ports.Equal(priorModel.Ports) ||
				!upgradedModel.Timeouts.Equal(priorModel.Timeouts) {
				t.Fatal("state migration did not preserve non-protocol attributes")
			}
		})
	}
}
