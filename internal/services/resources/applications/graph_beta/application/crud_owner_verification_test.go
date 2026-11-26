package graphBetaApplication

import (
	"context"
	"net/http"
	"testing"
	"testing/synctest"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUnitApplicationMissingOwnerVerification(t *testing.T) {
	for _, test := range []struct {
		name                       string
		present, paginated, denied bool
		wantDeletes, wantReads     int
	}{
		{"absent", false, false, false, 1, 1},
		{"still_present", true, false, false, 2, 1},
		{"present_on_next_page", true, true, false, 2, 2},
		{"verification_denied", false, false, true, 1, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				const appID = "11111111-1111-4111-8111-111111111111"
				const ownerID = "22222222-2222-4222-8222-222222222222"
				const ownersURL = "https://graph.microsoft.com/beta/applications/" + appID + "/owners"
				deletes, reads := 0, 0
				httpmock.RegisterResponder("DELETE", ownersURL+"/"+ownerID+"/$ref", func(_ *http.Request) (*http.Response, error) {
					deletes++
					if test.present && deletes > 1 {
						return httpmock.NewStringResponse(204, ""), nil
					}
					return httpmock.NewJsonResponse(404, map[string]any{"error": map[string]any{"code": "Request_ResourceNotFound", "message": "Resource '" + appID + "' does not exist or one of its queried reference-property objects are not present."}})
				})
				httpmock.RegisterResponder("GET", ownersURL, func(_ *http.Request) (*http.Response, error) {
					reads++
					if test.denied {
						return httpmock.NewJsonResponse(403, map[string]any{"error": map[string]any{"code": "Authorization_RequestDenied", "message": "Owner verification denied"}})
					}
					body := map[string]any{"value": []any{}}
					if test.paginated {
						body["@odata.nextLink"] = ownersURL + "?page=2"
					} else if test.present {
						body["value"] = []any{map[string]any{"id": ownerID}}
					}
					return httpmock.NewJsonResponse(200, body)
				})
				httpmock.RegisterResponder("GET", ownersURL+"?page=2", func(_ *http.Request) (*http.Response, error) {
					reads++
					return httpmock.NewJsonResponse(200, map[string]any{"value": []any{map[string]any{"id": ownerID}}})
				})
				client := providerclient.NewMockGraphClients(&http.Client{}).MockBetaClient
				err := RemoveOwner(context.Background(), client.GetAdapter(), appID, ownerID)
				if test.denied {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, test.wantDeletes, deletes)
				require.Equal(t, test.wantReads, reads)
			})
		})
	}
}
