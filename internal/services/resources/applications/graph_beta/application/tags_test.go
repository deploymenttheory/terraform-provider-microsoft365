package graphBetaApplication

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

func TestUnitResourceApplication_01_TagsRequestPayload(t *testing.T) {
	for _, test := range []struct {
		name string
		tags types.Set
		want []string
	}{
		{name: "omitted", tags: types.SetNull(types.StringType)},
		{name: "unknown", tags: types.SetUnknown(types.StringType)},
		{name: "empty", tags: types.SetValueMust(types.StringType, []attr.Value{}), want: []string{}},
		{name: "populated", tags: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("acceptance")}), want: []string{"acceptance"}},
	} {
		for _, create := range []bool{true, false} {
			operation := "update"
			if create {
				operation = "create"
			}
			t.Run(operation+"/"+test.name, func(t *testing.T) {
				ctx := context.Background()
				client := providerclient.NewMockGraphClients(&http.Client{}).MockBetaClient
				model := ApplicationResourceModel{DisplayName: types.StringValue("Application tags regression"), Tags: test.tags}
				body, err := constructResource(ctx, client, &model, "11111111-1111-1111-1111-111111111111", create)
				require.NoError(t, err)
				var request *abstractions.RequestInformation
				if create {
					request, err = client.Applications().ToPostRequestInformation(ctx, body, nil)
				} else {
					request, err = client.Applications().ByApplicationId("11111111-1111-1111-1111-111111111111").ToPatchRequestInformation(ctx, body, nil)
				}
				require.NoError(t, err)
				var payload map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(request.Content, &payload))
				if test.want == nil {
					require.NotContains(t, payload, "tags", "unset tags must be omitted rather than serialized as null")
				} else {
					var actual []string
					require.NoError(t, json.Unmarshal(payload["tags"], &actual))
					require.Equal(t, test.want, actual)
				}
			})
		}
	}
}
