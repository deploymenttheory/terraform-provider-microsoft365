package graphBetaCrossTenantAccessPartnerSettings

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

func TestUnitCrossTenantPartnerHardDeletePolicyRecycleBin(t *testing.T) {
	for _, tc := range []struct {
		name                                                                                  string
		alreadySoftDeleted, alreadyGone, softOnly, readDenied, hardDenied, concurrentlyPurged bool
		propagationDelay                                                                      int
	}{
		{name: "normal"},
		{name: "already_soft_deleted", alreadySoftDeleted: true},
		{name: "already_permanently_deleted", alreadySoftDeleted: true, alreadyGone: true},
		{name: "soft_delete_only", softOnly: true},
		{name: "purge_permission_denied", hardDenied: true},
		{name: "verification_permission_denied", readDenied: true},
		{name: "concurrent_purge_404", concurrentlyPurged: true},
		{name: "recycle_bin_propagation", propagationDelay: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tenantID := "12345678-1234-1234-1234-123456789012"
			var mu sync.Mutex
			active := !tc.alreadySoftDeleted
			recycled := tc.alreadySoftDeleted && !tc.alreadyGone
			directoryDeletes, policyDeletes := 0, 0
			beforeDelay, afterDelay := tc.propagationDelay, tc.propagationDelay
			purging := false
			falseAbsence := tc.propagationDelay > 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				notFound := func() {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprint(w, "{\"error\":{\"code\":\"Request_ResourceNotFound\",\"message\":\"Not found in this store\"}}")
				}
				denied := func() {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprint(w, "{\"error\":{\"code\":\"ResourceNotFound\",\"message\":\"Permission denied despite a misleading code\"}}")
				}
				switch {
				case req.Method == http.MethodDelete && req.URL.Path == "/policies/crossTenantAccessPolicy/partners/"+tenantID:
					if !active {
						notFound()
						return
					}
					active = false
					recycled = true
					w.WriteHeader(http.StatusNoContent)
				case req.Method == http.MethodGet && req.URL.Path == "/policies/crossTenantAccessPolicy/partners/"+tenantID:
					if tc.readDenied {
						denied()
						return
					}
					notFound()
				case req.Method == http.MethodDelete && req.URL.Path == "/directory/deletedItems/"+tenantID:
					directoryDeletes++
					notFound()
				case req.Method == http.MethodGet && req.URL.Path == "/policies/deletedItems/crossTenantPartners/"+tenantID:
					if !recycled {
						notFound()
						return
					}
					if beforeDelay > 0 {
						beforeDelay--
						notFound()
						return
					}
					if purging {
						if falseAbsence {
							falseAbsence = false
							notFound()
							return
						}
						if afterDelay == 0 {
							recycled = false
							notFound()
							return
						}
						afterDelay--
					}
					fmt.Fprintf(w, "{\"tenantId\":%q}", tenantID)
				case req.Method == http.MethodDelete && req.URL.Path == "/policies/deletedItems/crossTenantPartners/"+tenantID:
					policyDeletes++
					if tc.hardDenied {
						denied()
						return
					}
					if tc.concurrentlyPurged {
						recycled = false
						notFound()
						return
					}
					purging = true
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected request %s %s", req.Method, req.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer server.Close()
			adapter, err := msgraphbetasdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
				&authentication.AnonymousAuthenticationProvider{},
				nil,
				nil,
				server.Client(),
			)
			require.NoError(t, err)
			adapter.SetBaseUrl(server.URL)
			r := NewCrossTenantAccessPartnerSettingsResource().(*CrossTenantAccessPartnerSettingsResource)
			r.client = msgraphbetasdk.NewGraphServiceClient(adapter)
			ctx := context.Background()
			var schema resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schema)
			objectType := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
			values := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
			for name, typ := range objectType.AttributeTypes {
				values[name] = tftypes.NewValue(typ, nil)
			}
			state := tfsdk.State{Schema: schema.Schema, Raw: tftypes.NewValue(objectType, values)}
			require.False(t, state.SetAttribute(ctx, path.Root("id"), types.StringValue(tenantID)).HasError())
			require.False(t, state.SetAttribute(ctx, path.Root("tenant_id"), types.StringValue(tenantID)).HasError())
			require.False(t, state.SetAttribute(ctx, path.Root("hard_delete"), types.BoolValue(!tc.softOnly)).HasError())
			resp := resource.DeleteResponse{State: state}
			r.Delete(ctx, resource.DeleteRequest{State: state}, &resp)
			mu.Lock()
			defer mu.Unlock()
			require.Zero(t, directoryDeletes, "cross-tenant policies are not directory deleted items")
			if tc.hardDenied || tc.readDenied {
				require.True(t, resp.Diagnostics.HasError(), resp.Diagnostics)
				require.False(t, resp.State.Raw.IsNull(), "A failed purge must preserve Terraform state for retry")
				require.True(t, recycled)
				return
			}
			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			require.True(t, resp.State.Raw.IsNull())
			if tc.softOnly {
				require.True(t, recycled)
				require.Zero(t, policyDeletes)
			} else {
				require.False(t, recycled, "Terraform reported successful hard deletion but the partner remains in policy deleted items")
				if tc.alreadyGone {
					require.Zero(t, policyDeletes)
				} else {
					require.Equal(t, 1, policyDeletes)
				}
			}
		})
	}
}
