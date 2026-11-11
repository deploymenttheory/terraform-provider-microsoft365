package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"testing"

	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/types"
	serialization "github.com/microsoft/kiota-serialization-json-go"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceCloudPcProvisioningPolicy_UserExperienceSerialization(t *testing.T) {
	ctx := context.Background()
	for _, value := range []string{"cloudPc", "cloudApp", ""} {
		t.Run(value, func(t *testing.T) {
			experience := types.StringNull()
			expected := "cloudPc"
			if value != "" {
				experience = types.StringValue(value)
				expected = value
			}
			body, err := constructResource(ctx, &CloudPcProvisioningPolicyResourceModel{ID: types.StringValue("existing-policy"), ScopeIds: types.SetNull(types.StringType), UserExperienceType: experience, ProvisioningType: types.StringValue("sharedByEntraGroup")}, &msgraph.GraphServiceClient{})
			require.NoError(t, err)
			writer := serialization.NewJsonSerializationWriter()
			require.NoError(t, writer.WriteObjectValue("", body))
			encoded, err := writer.GetSerializedContent()
			require.NoError(t, err)
			var payload map[string]any
			require.NoError(t, json.Unmarshal(encoded, &payload))
			require.Equal(t, expected, payload["userExperienceType"])
			require.NotContains(t, body.GetAdditionalData(), "userExperienceType")
			var state CloudPcProvisioningPolicyResourceModel
			MapRemoteStateToTerraform(ctx, &state, body)
			require.Equal(t, expected, state.UserExperienceType.ValueString())
		})
	}
	_, err := constructResource(ctx, &CloudPcProvisioningPolicyResourceModel{ID: types.StringValue("existing-policy"), ScopeIds: types.SetNull(types.StringType), UserExperienceType: types.StringValue("cloudApp"), ProvisioningType: types.StringValue("dedicated")}, &msgraph.GraphServiceClient{})
	require.ErrorContains(t, err, "sharedByEntraGroup")
	var legacy CloudPcProvisioningPolicyResourceModel
	MapRemoteStateToTerraform(ctx, &legacy, models.NewCloudPcProvisioningPolicy())
	require.Equal(t, "cloudPc", legacy.UserExperienceType.ValueString())
}
