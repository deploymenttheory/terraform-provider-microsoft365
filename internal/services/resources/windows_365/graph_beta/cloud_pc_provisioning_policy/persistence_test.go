package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"testing"

	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	serialization "github.com/microsoft/kiota-serialization-json-go"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceCloudPcProvisioningPolicy_PersistenceSerialization(t *testing.T) {
	ctx := context.Background()
	objectTypes := map[string]attr.Type{"enabled": types.BoolType, "storage_size_category": types.StringType}
	for _, size := range []string{"fourGB", "eightGB", "sixteenGB", "thirtyTwoGB", "sixtyFourGB"} {
		for _, enabled := range []bool{true, false} {
			t.Run(size+types.BoolValue(enabled).String(), func(t *testing.T) {
				settings := types.ObjectValueMust(objectTypes, map[string]attr.Value{"enabled": types.BoolValue(enabled), "storage_size_category": types.StringValue(size)})
				body, err := constructResource(ctx, &CloudPcProvisioningPolicyResourceModel{ID: types.StringValue("existing-policy"), ScopeIds: types.SetNull(types.StringType), ProvisioningType: types.StringValue("sharedByEntraGroup"), UserSettingsPersistenceConfiguration: settings}, &msgraph.GraphServiceClient{})
				require.NoError(t, err)
				writer := serialization.NewJsonSerializationWriter()
				require.NoError(t, writer.WriteObjectValue("", body))
				encoded, err := writer.GetSerializedContent()
				require.NoError(t, err)
				var payload map[string]any
				require.NoError(t, json.Unmarshal(encoded, &payload))
				actual := payload["userSettingsPersistenceConfiguration"].(map[string]any)
				require.Equal(t, enabled, actual["userSettingsPersistenceEnabled"])
				require.Equal(t, size, actual["userSettingsPersistenceStorageSizeCategory"])
				var state CloudPcProvisioningPolicyResourceModel
				MapRemoteStateToTerraform(ctx, &state, body)
				require.True(t, settings.Equal(state.UserSettingsPersistenceConfiguration))
			})
		}
	}
	body, err := constructResource(ctx, &CloudPcProvisioningPolicyResourceModel{ID: types.StringValue("existing-policy"), ScopeIds: types.SetNull(types.StringType), UserSettingsPersistenceConfiguration: types.ObjectNull(objectTypes)}, &msgraph.GraphServiceClient{})
	require.NoError(t, err)
	require.Nil(t, body.GetUserSettingsPersistenceConfiguration())
	var state CloudPcProvisioningPolicyResourceModel
	MapRemoteStateToTerraform(ctx, &state, models.NewCloudPcProvisioningPolicy())
	require.True(t, state.UserSettingsPersistenceConfiguration.IsNull())
	for _, size := range []types.String{types.StringNull(), types.StringUnknown()} {
		settings := types.ObjectValueMust(objectTypes, map[string]attr.Value{"enabled": types.BoolValue(true), "storage_size_category": size})
		body, err := constructResource(ctx, &CloudPcProvisioningPolicyResourceModel{ID: types.StringValue("existing-policy"), ScopeIds: types.SetNull(types.StringType), ProvisioningType: types.StringValue("sharedByEntraGroup"), UserSettingsPersistenceConfiguration: settings}, &msgraph.GraphServiceClient{})
		require.NoError(t, err)
		require.Equal(t, "fourGB", body.GetUserSettingsPersistenceConfiguration().GetUserSettingsPersistenceStorageSizeCategory().String())
	}

}
