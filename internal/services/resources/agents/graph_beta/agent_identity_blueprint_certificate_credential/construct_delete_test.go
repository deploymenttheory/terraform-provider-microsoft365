package graphBetaAgentIdentityBlueprintCertificateCredential

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

func TestUnitAgentIdentityBlueprintCertificateDeletePayloads(t *testing.T) {
	removedID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	siblingID := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	for _, test := range []struct {
		name    string
		removed bool
		sibling bool
	}{
		{"last_certificate", true, false},
		{"already_empty", false, false},
		{"preserve_sibling", true, true},
		{"missing_certificate_preserves_sibling", false, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			existing := []graphmodels.KeyCredentialable{}
			if test.removed {
				credential := graphmodels.NewKeyCredential()
				credential.SetKeyId(&removedID)
				existing = append(existing, credential)
			}
			if test.sibling {
				credential := graphmodels.NewKeyCredential()
				credential.SetKeyId(&siblingID)
				credential.SetKey([]byte("sibling"))
				existing = append(existing, credential)
			}
			ctx := context.Background()
			body, err := constructDeleteResource(ctx, removedID.String(), existing)
			require.NoError(t, err)
			client := providerclient.NewMockGraphClients(&http.Client{}).MockBetaClient
			request, err := client.Applications().ByApplicationId("33333333-3333-4333-8333-333333333333").ToPatchRequestInformation(ctx, body, nil)
			require.NoError(t, err)
			var payload map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(request.Content, &payload))
			if !test.sibling {
				require.JSONEq(t, "[]", string(payload["keyCredentials"]))
				return
			}
			var credentials []map[string]any
			require.NoError(t, json.Unmarshal(payload["keyCredentials"], &credentials))
			require.Len(t, credentials, 1)
			require.Equal(t, siblingID.String(), credentials[0]["keyId"])
			require.Equal(t, "c2libGluZw==", credentials[0]["key"])
		})
	}
}
