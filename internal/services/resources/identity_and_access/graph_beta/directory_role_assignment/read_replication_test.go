package graphBetaDirectoryRoleAssignment

import (
	"context"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUnitDirectoryRoleAssignmentReadReplication(t *testing.T) {
	for _, test := range []struct {
		name         string
		missingReads int
		deadline     bool
		wantMissing  bool
		wantError    bool
		wantGets     int
	}{
		{"one_missing_replica", 1, false, false, false, 2},
		{"two_missing_replicas", 2, false, false, false, 3},
		{"three_missing_replicas", 3, false, false, false, 4},
		{"confirmed_absence", 100, false, true, false, 16},
		{"timeout_preserves_state", 10, true, false, true, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				transport := httpmock.NewMockTransport()
				gets := 0
				transport.RegisterResponder("GET", "https://graph.microsoft.com/beta/roleManagement/directory/roleAssignments/"+testAssignmentID, func(_ *http.Request) (*http.Response, error) {
					gets++
					if gets <= test.missingReads {
						return httpmock.NewJsonResponse(404, map[string]any{"error": map[string]any{"code": "Request_ResourceNotFound", "message": "Assignment missing on this replica"}})
					}
					return httpmock.NewJsonResponse(200, map[string]any{
						"id": testAssignmentID, "principalId": "11111111-1111-4111-8111-111111111111",
						"roleDefinitionId": "88d8e3e3-8f55-4a1e-953a-9b9898b8876b", "directoryScopeId": "/",
					})
				})
				r := NewDirectoryRoleAssignmentResource().(*DirectoryRoleAssignmentResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				state := assignmentTestState(t, r)
				ctx := context.Background()
				if test.deadline {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, time.Second)
					defer cancel()
				}
				response := resource.ReadResponse{State: state}
				r.Read(ctx, resource.ReadRequest{State: state}, &response)
				require.Equal(t, test.wantError, response.Diagnostics.HasError(), response.Diagnostics)
				require.Equal(t, test.wantMissing, response.State.Raw.IsNull())
				require.Equal(t, test.wantGets, gets)
			})
		})
	}
}
