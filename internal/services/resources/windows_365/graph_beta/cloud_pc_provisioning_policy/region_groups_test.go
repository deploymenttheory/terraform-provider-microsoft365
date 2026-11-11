package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	serialization "github.com/microsoft/kiota-serialization-json-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnitResourceCloudPcProvisioningPolicy_RegionGroups(t *testing.T) {
	ctx := context.Background()
	var response resource.SchemaResponse
	NewCloudPcProvisioningPolicyResource().Schema(ctx, resource.SchemaRequest{}, &response)
	attribute := response.Schema.Attributes["domain_join_configurations"].(schema.ListNestedAttribute).NestedObject.Attributes["region_group"].(schema.StringAttribute)
	for _, value := range []string{"default", "australia", "canada", "usCentral", "usEast", "usWest", "france", "germany", "europeUnion", "unitedKingdom", "japan", "asia", "india", "southAmerica", "euap", "usGovernment", "usGovernmentDOD", "unknownFutureValue", "norway", "switzerland", "southKorea", "middleEast", "mexico", "australasia", "europe", "singapore", "hongKong", "ireland", "sweden", "poland", "italy", "spain", "netherlands", "brazil", "israel", "automatic", "indonesia", "taiwan", "malaysia", "newZealand", "austria", "denmark", "belgium", "kenya"} {
		t.Run(value, func(t *testing.T) {
			var validation validator.StringResponse
			for _, v := range attribute.Validators {
				v.ValidateString(ctx, validator.StringRequest{Path: path.Root("domain_join_configurations").AtListIndex(0).AtName("region_group"), ConfigValue: types.StringValue(value)}, &validation)
			}
			require.False(t, validation.Diagnostics.HasError(), validation.Diagnostics)
			body, err := constructResource(ctx, &CloudPcProvisioningPolicyResourceModel{ID: types.StringValue("existing-policy"), DisplayName: types.StringValue("Region test"), ScopeIds: types.SetNull(types.StringType), DomainJoinConfigurations: []DomainJoinConfigurationModel{{DomainJoinType: types.StringValue("azureADJoin"), RegionGroup: types.StringValue(value), RegionName: types.StringValue("automatic")}}}, &msgraphbetasdk.GraphServiceClient{})
			require.NoError(t, err)
			writer := serialization.NewJsonSerializationWriter()
			require.NoError(t, writer.WriteObjectValue("", body))
			encoded, err := writer.GetSerializedContent()
			require.NoError(t, err)
			var payload map[string]any
			require.NoError(t, json.Unmarshal(encoded, &payload))
			join := payload["domainJoinConfigurations"].([]any)[0].(map[string]any)
			require.Equal(t, value, join["regionGroup"])
			var state CloudPcProvisioningPolicyResourceModel
			MapRemoteStateToTerraform(ctx, &state, body)
			require.Equal(t, value, state.DomainJoinConfigurations[0].RegionGroup.ValueString())
		})
	}
	for _, value := range []string{"uae", "southAfrica", "notARegion"} {
		var validation validator.StringResponse
		for _, v := range attribute.Validators {
			v.ValidateString(ctx, validator.StringRequest{ConfigValue: types.StringValue(value)}, &validation)
		}
		require.True(t, validation.Diagnostics.HasError(), value)
	}
}
