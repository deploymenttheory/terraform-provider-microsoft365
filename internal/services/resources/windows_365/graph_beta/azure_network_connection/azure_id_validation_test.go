package graphBetaAzureNetworkConnection

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceAzureNetworkConnection_AzureIDValidation(t *testing.T) {
	ctx := context.Background()
	var response resource.SchemaResponse
	NewCloudPcOnPremisesConnectionResource().Schema(ctx, resource.SchemaRequest{}, &response)

	const subscription = "abcdefab-1234-5678-9abc-abcdefabcdef"
	const group = "/subscriptions/" + subscription + "/resourceGroups/ExampleGroup"
	const network = group + "/providers/Microsoft.Network/virtualNetworks/ExampleNetwork"
	lowercaseNames := strings.NewReplacer("ExampleGroup", "examplegroup", "ExampleNetwork", "examplenetwork", "CloudPCs", "cloudpcs")
	uppercaseNames := strings.NewReplacer("ExampleGroup", "EXAMPLEGROUP", "ExampleNetwork", "EXAMPLENETWORK", "CloudPCs", "CLOUDPCS")

	for _, attribute := range []struct {
		name  string
		value string
	}{
		{"resource_group_id", group},
		{"virtual_network_id", network},
		{"subnet_id", network + "/subnets/CloudPCs"},
	} {
		t.Run(attribute.name, func(t *testing.T) {
			validators := response.Schema.Attributes[attribute.name].(schema.StringAttribute).Validators
			require.NotEmpty(t, validators)
			for _, tc := range []struct {
				name      string
				value     types.String
				wantError bool
			}{
				{"lowercase_names", types.StringValue(lowercaseNames.Replace(attribute.value)), false},
				{"mixed_case_names", types.StringValue(attribute.value), false},
				{"uppercase_names", types.StringValue(uppercaseNames.Replace(attribute.value)), false},
				{"uppercase_subscription", types.StringValue(strings.Replace(attribute.value, subscription, strings.ToUpper(subscription), 1)), false},
				{"null", types.StringNull(), false},
				{"unknown", types.StringUnknown(), false},
				{"empty", types.StringValue(""), true},
				{"not_an_id", types.StringValue("ExampleGroup"), true},
				{"missing_subscription", types.StringValue(strings.Replace(attribute.value, subscription, "", 1)), true},
				{"invalid_subscription", types.StringValue(strings.Replace(attribute.value, subscription, "not-a-guid", 1)), true},
				{"missing_group", types.StringValue(strings.Replace(attribute.value, "ExampleGroup", "", 1)), true},
				{"space_in_name", types.StringValue(strings.Replace(attribute.value, "ExampleGroup", "Example Group", 1)), true},
				{"wrong_path", types.StringValue(strings.Replace(attribute.value, "/resourceGroups/", "/resourceGroup/", 1)), true},
				{"extra_segment", types.StringValue(attribute.value + "/unexpected"), true},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var validation validator.StringResponse
					for _, v := range validators {
						v.ValidateString(ctx, validator.StringRequest{Path: path.Root(attribute.name), ConfigValue: tc.value}, &validation)
					}
					require.Equal(t, tc.wantError, validation.Diagnostics.HasError(), validation.Diagnostics)
				})
			}
		})
	}
}

func TestUnitResourceAzureNetworkConnection_PreservesAzureIDCasing(t *testing.T) {
	ctx := context.Background()
	const subscription = "ABCDEFAB-1234-5678-9ABC-ABCDEFABCDEF"
	const group = "/subscriptions/" + subscription + "/resourceGroups/ExampleGroup"
	const network = group + "/providers/Microsoft.Network/virtualNetworks/ExampleNetwork"
	const subnet = network + "/subnets/CloudPCs"
	data := CloudPcOnPremisesConnectionResourceModel{
		DisplayName:      types.StringValue("Example connection"),
		ConnectionType:   types.StringValue("azureADJoin"),
		ResourceGroupId:  types.StringValue(group),
		VirtualNetworkId: types.StringValue(network),
		SubnetId:         types.StringValue(subnet),
		SubscriptionId:   types.StringValue(subscription),
	}
	body, err := constructResource(ctx, &data)
	require.NoError(t, err)
	require.Equal(t, group, *body.GetResourceGroupId())
	require.Equal(t, network, *body.GetVirtualNetworkId())
	require.Equal(t, subnet, *body.GetSubnetId())
	require.Equal(t, subscription, *body.GetSubscriptionId())

	var state CloudPcOnPremisesConnectionResourceModel
	MapRemoteStateToTerraform(ctx, &state, body, types.StringNull())
	require.Equal(t, data.ResourceGroupId, state.ResourceGroupId)
	require.Equal(t, data.VirtualNetworkId, state.VirtualNetworkId)
	require.Equal(t, data.SubnetId, state.SubnetId)
	require.Equal(t, data.SubscriptionId, state.SubscriptionId)
}
