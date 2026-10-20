package graphBetaMobileAppSupersedence

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	serialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitSupersedenceConstructWritableFields(t *testing.T) {
	for _, kind := range []string{"update", "replace"} {
		t.Run(kind, func(t *testing.T) {
			model := MobileAppSupersedenceResourceModel{
				SourceID:            types.StringValue("11111111-1111-1111-1111-111111111111"),
				TargetID:            types.StringValue("22222222-2222-2222-2222-222222222222"),
				SupersedenceType:    types.StringValue(kind),
				SourceDisplayName:   types.StringValue("computed source"),
				TargetDisplayName:   types.StringValue("computed target"),
				AppRelationshipType: types.StringUnknown(),
			}
			body, err := constructResource(context.Background(), &model)
			require.NoError(t, err)
			require.Equal(t, kind, body.GetSupersedenceType().String())
			require.Equal(t, graphmodels.CHILD_MOBILEAPPRELATIONSHIPTYPE, *body.GetTargetType())
			writer := serialization.NewJsonSerializationWriter()
			require.NoError(t, writer.WriteObjectValue("", body))
			data, err := writer.GetSerializedContent()
			require.NoError(t, err)
			var payload map[string]any
			require.NoError(t, json.Unmarshal(data, &payload))
			require.Equal(t, "#microsoft.graph.mobileAppSupersedence", payload["@odata.type"])
			require.Equal(t, kind, payload["supersedenceType"])
			require.Equal(t, "child", payload["targetType"])
			require.NotContains(t, payload, "targetDisplayName")
			require.NotContains(t, payload, "sourceDisplayName")
		})
	}
}

func TestUnitSupersedenceConstructRejectsInvalidTypes(t *testing.T) {
	for _, value := range []types.String{types.StringValue("Update"), types.StringValue("invalid"), types.StringNull(), types.StringUnknown()} {
		_, err := constructResource(
			context.Background(),
			&MobileAppSupersedenceResourceModel{SupersedenceType: value},
		)
		require.Error(t, err)
	}
}
