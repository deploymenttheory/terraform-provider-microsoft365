package graphBetaApplicationsAgentIdentityBlueprintServicePrincipal

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

func TestUnitAgentIdentityBlueprintPrincipalTagPayload(t *testing.T) {
	for _, create := range []bool{true, false} {
		operation := "update"
		if create {
			operation = "create"
		}
		for _, form := range []string{"null", "unknown", "empty", "populated"} {
			t.Run(operation+"/"+form, func(t *testing.T) {
				tags := types.SetNull(types.StringType)
				switch form {
				case "unknown":
					tags = types.SetUnknown(types.StringType)
				case "empty":
					tags = types.SetValueMust(types.StringType, []attr.Value{})
				case "populated":
					tags = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("acceptance")})
				}
				model := AgentIdentityBlueprintServicePrincipalResourceModel{
					AppId: types.StringValue("11111111-1111-4111-8111-111111111111"), Tags: tags,
				}
				ctx := context.Background()
				body, err := constructResource(ctx, &model)
				require.NoError(t, err)
				client := providerclient.NewMockGraphClients(&http.Client{}).MockBetaClient
				request, err := client.ServicePrincipals().ToPostRequestInformation(ctx, body, nil)
				if !create {
					request, err = client.ServicePrincipals().ByServicePrincipalId("11111111-1111-4111-8111-111111111111").ToPatchRequestInformation(ctx, body, nil)
				}
				require.NoError(t, err)
				var payload map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(request.Content, &payload))
				if form == "null" || form == "unknown" {
					require.NotContains(t, payload, "tags")
				} else if form == "empty" {
					require.JSONEq(t, "[]", string(payload["tags"]))
				} else {
					require.JSONEq(t, `["acceptance"]`, string(payload["tags"]))
				}
				require.JSONEq(t, `"#microsoft.graph.agentIdentityBlueprintPrincipal"`, string(payload["@odata.type"]))
			})
		}
	}
}
