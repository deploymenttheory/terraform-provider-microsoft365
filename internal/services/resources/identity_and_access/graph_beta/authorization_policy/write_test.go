package graphBetaAuthorizationPolicy

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceAuthorizationPolicy_22_WriteErrorsAndOmittedInvitations(t *testing.T) {
	for _, test := range []struct {
		name        string
		invites     bool
		permissions bool
		failAt      int
		calls       int
	}{
		{"initial_write_failure", true, true, 1, 1},
		{"invitation_write_failure", true, true, 2, 2},
		{"omitted_invitations", false, true, 0, 1},
		{"invitation_only", true, false, 0, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			remote := map[string]any{"id": singletonID}
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				if q.Method == http.MethodGet {
					require.Equal(t, readPath, q.URL.Path)
					encoded, err := json.Marshal(map[string]any{"value": []any{remote}})
					require.NoError(t, err)
					reply(w, 200, string(encoded))
					return
				}
				require.Equal(t, http.MethodPatch, q.Method)
				require.Equal(t, patchPath, q.URL.Path)
				calls++
				if calls == test.failAt {
					reply(w, http.StatusForbidden, `{"error":{"code":"Forbidden","message":"Policy write denied"}}`)
					return
				}
				var reader io.Reader = q.Body
				if q.Header.Get("Content-Encoding") == "gzip" {
					compressed, err := gzip.NewReader(q.Body)
					require.NoError(t, err)
					defer compressed.Close()
					reader = compressed
				}
				var payload map[string]any
				require.NoError(t, json.NewDecoder(reader).Decode(&payload))
				for key, value := range payload {
					remote[key] = value
				}
				reply(w, http.StatusNoContent, "")
			})
			model := testModel()
			if !test.invites {
				model.AllowInvitesFrom = types.StringUnknown()
			}
			if !test.permissions {
				model.DefaultUserRolePermissions = types.ObjectNull(defaultUserRolePermissionsAttributeTypes())
			}
			body, err := constructResource(context.Background(), &model)
			require.NoError(t, err)
			state := testState(t, r, model)
			err = r.patchSettings(context.Background(), body)
			if err == nil {
				err = r.completeInvitationWrite(context.Background(), body, state)
			}
			if test.failAt != 0 {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.calls, calls)
		})
	}
}
