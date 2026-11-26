package graphBetaApplication

import (
	"context"
	"net/http"
	"testing"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUnitApplicationRemoveOwnerAlreadyAbsent(t *testing.T) {
	const message = "One or more removed object references do not exist for the following modified properties: 'owners'."
	for _, test := range []struct {
		name      string
		status    int
		message   string
		wantError bool
	}{
		{"already_absent", 400, message, false},
		{"other_validation", 400, "The supplied identifier is invalid.", true},
		{"permission_error", 403, message, true},
		{"different_relationship", 400, "One or more removed object references do not exist for the following modified properties: 'groups'.", true},
		{"success", 204, "", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/beta/applications/11111111-1111-4111-8111-111111111111/owners/22222222-2222-4222-8222-222222222222/$ref", func(_ *http.Request) (*http.Response, error) {
				if test.status == http.StatusNoContent {
					return httpmock.NewStringResponse(test.status, ""), nil
				}
				return httpmock.NewJsonResponse(test.status, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": test.message}})
			})
			client := providerclient.NewMockGraphClients(&http.Client{}).MockBetaClient
			err := RemoveOwner(context.Background(), client.GetAdapter(), "11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222")
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUnitApplicationRemoveOwnerReplication(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	requests := 0
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/applications/11111111-1111-4111-8111-111111111111/owners", httpmock.NewJsonResponderOrPanic(404, map[string]any{"error": map[string]any{"code": "Request_ResourceNotFound", "message": "Application not visible on this replica"}}))
	httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/beta/applications/11111111-1111-4111-8111-111111111111/owners/22222222-2222-4222-8222-222222222222/$ref", func(_ *http.Request) (*http.Response, error) {
		requests++
		if requests == 1 {
			return httpmock.NewJsonResponse(404, map[string]any{"error": map[string]any{"code": "Request_ResourceNotFound", "message": "Resource '11111111-1111-4111-8111-111111111111' does not exist or one of its queried reference-property objects are not present."}})
		}
		return httpmock.NewJsonResponse(400, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": "One or more removed object references do not exist for the following modified properties: 'owners'."}})
	})
	client := providerclient.NewMockGraphClients(&http.Client{}).MockBetaClient
	err := RemoveOwner(context.Background(), client.GetAdapter(), "11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222")
	require.NoError(t, err)
	require.Equal(t, 2, requests)
}
