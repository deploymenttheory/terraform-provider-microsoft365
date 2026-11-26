package graphBetaServicePrincipal

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

func TestUnitServicePrincipalCollectionPayloads(t *testing.T) {
	for _, form := range []string{"null", "unknown", "empty", "populated"} {
		for _, update := range []bool{false, true} {
			operation := "create"
			if update {
				operation = "update"
			}
			t.Run(operation+"/"+form, func(t *testing.T) {
				ctx := context.Background()
				set := types.SetNull(types.StringType)
				switch form {
				case "unknown":
					set = types.SetUnknown(types.StringType)
				case "empty":
					set = types.SetValueMust(types.StringType, []attr.Value{})
				case "populated":
					set = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("collection-value")})
				}
				model := ServicePrincipalResourceModel{
					AppID: types.StringValue("11111111-1111-4111-8111-111111111111"),
					Tags:  set, NotificationEmailAddresses: set, AlternativeNames: set,
					SamlSingleSignOnSettings: types.ObjectNull(map[string]attr.Type{"relay_state": types.StringType}),
					TokenEncryptionKeyID:     types.StringNull(),
				}
				body, err := constructResource(ctx, &model, update)
				require.NoError(t, err)
				client := providerclient.NewMockGraphClients(&http.Client{}).MockBetaClient
				var request *abstractions.RequestInformation
				if update {
					request, err = client.ServicePrincipals().ByServicePrincipalId("11111111-1111-4111-8111-111111111111").ToPatchRequestInformation(ctx, body, nil)
				} else {
					request, err = client.ServicePrincipals().ToPostRequestInformation(ctx, body, nil)
				}
				require.NoError(t, err)
				var payload map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(request.Content, &payload))
				for _, property := range []string{"tags", "notificationEmailAddresses", "alternativeNames"} {
					if form == "null" || form == "unknown" {
						require.NotContains(t, payload, property, "unset non-nullable collections must be omitted")
						continue
					}
					require.Contains(t, payload, property)
					if form == "empty" {
						require.JSONEq(t, "[]", string(payload[property]), "explicit empties must still clear remote values")
					} else {
						require.JSONEq(t, `["collection-value"]`, string(payload[property]))
					}
				}
			})
		}
	}
}
