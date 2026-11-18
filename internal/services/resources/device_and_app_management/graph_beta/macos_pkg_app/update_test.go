package graphBetaMacOSPKGApp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"sync/atomic"
	"testing"
	"testing/synctest"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func nullAttributes(objectType tftypes.Object) map[string]tftypes.Value {
	values := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attributeType := range objectType.AttributeTypes {
		values[name] = tftypes.NewValue(attributeType, nil)
	}
	return values
}

func metadataUpdateState(t *testing.T, r *MacOSPKGAppResource, sourceField, source string, updated bool) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	var response resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &response)
	rootType := response.Schema.Type().TerraformType(ctx).(tftypes.Object)
	values := nullAttributes(rootType)
	values["id"] = tftypes.NewValue(tftypes.String, "11111111-1111-1111-1111-111111111111")
	values["display_name"] = tftypes.NewValue(tftypes.String, "unit-test-macos-pkg-metadata")
	description := "Original app metadata"
	if updated {
		description = "Updated app metadata"
	}
	values["description"] = tftypes.NewValue(tftypes.String, description)
	values["publisher"] = tftypes.NewValue(tftypes.String, "MacAdmins")
	values["is_featured"] = tftypes.NewValue(tftypes.Bool, false)
	values["role_scope_tag_ids"] = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "0")})
	if sourceField != "" {
		installerType := rootType.AttributeTypes["app_installer"].(tftypes.Object)
		installer := nullAttributes(installerType)
		installer[sourceField] = tftypes.NewValue(tftypes.String, source)
		values["app_installer"] = tftypes.NewValue(installerType, installer)
	}
	pkgType := rootType.AttributeTypes["macos_pkg_app"].(tftypes.Object)
	pkg := nullAttributes(pkgType)
	pkg["ignore_version_detection"] = tftypes.NewValue(tftypes.Bool, updated)
	pkg["primary_bundle_id"] = tftypes.NewValue(tftypes.String, "com.github.macadmins.Nudge")
	pkg["primary_bundle_version"] = tftypes.NewValue(tftypes.String, "2.1.4.81863")
	osType := pkgType.AttributeTypes["minimum_supported_operating_system"].(tftypes.Object)
	minimumOS := nullAttributes(osType)
	for name := range minimumOS {
		minimumOS[name] = tftypes.NewValue(tftypes.Bool, name == "v12_0")
	}
	pkg["minimum_supported_operating_system"] = tftypes.NewValue(osType, minimumOS)
	values["macos_pkg_app"] = tftypes.NewValue(pkgType, pkg)
	return tfsdk.State{Schema: response.Schema, Raw: tftypes.NewValue(rootType, values)}
}

func TestUnitResourceMacOSPkgApp_02_UpdateWithoutInstaller(t *testing.T) {
	for _, sourceField := range []string{"installer_file_path_source", "installer_url_source", ""} {
		name := sourceField
		if name == "" {
			name = "no installer in imported state"
		}
		t.Run(name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				source := filepath.Join(t.TempDir(), "missing.pkg")
				var sourceRequests atomic.Int64
				var server *httptest.Server
				if sourceField == "installer_url_source" {
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						sourceRequests.Add(1)
						w.WriteHeader(http.StatusGone)
					}))
					defer server.Close()
					source = server.URL + "/Nudge.pkg"
				}
				transport := httpmock.NewMockTransport()
				appURI := "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/11111111-1111-1111-1111-111111111111"
				patches := 0
				remote := map[string]any{
					"@odata.type": "#microsoft.graph.macOSPkgApp",
					"id":          "11111111-1111-1111-1111-111111111111",
					"displayName": "unit-test-macos-pkg-metadata", "description": "Updated app metadata",
					"publisher": "MacAdmins", "isFeatured": false, "ignoreVersionDetection": true,
					"primaryBundleId": "com.github.macadmins.Nudge", "primaryBundleVersion": "2.1.4.81863",
					"roleScopeTagIds": []string{"0"}, "minimumSupportedOperatingSystem": map[string]bool{"v12_0": true},
				}
				transport.RegisterResponder("PATCH", appURI, func(req *http.Request) (*http.Response, error) {
					patches++
					body, err := io.ReadAll(req.Body)
					require.NoError(t, err)
					var payload map[string]any
					require.NoError(t, json.Unmarshal(body, &payload))
					require.Equal(t, "Updated app metadata", payload["description"])
					require.Equal(t, true, payload["ignoreVersionDetection"])
					for _, field := range []string{"fileName", "primaryBundleId", "primaryBundleVersion", "committedContentVersion"} {
						require.NotContains(t, payload, field, "metadata update must preserve uploaded content")
					}
					return httpmock.NewJsonResponse(http.StatusOK, remote)
				})
				responder, err := httpmock.NewJsonResponder(http.StatusOK, remote)
				require.NoError(t, err)
				transport.RegisterRegexpResponder("GET", regexp.MustCompile("^"+regexp.QuoteMeta(appURI)), responder)
				r := NewMacOSPKGAppResource().(*MacOSPKGAppResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				state := metadataUpdateState(t, r, sourceField, source, false)
				planned := metadataUpdateState(t, r, sourceField, source, true)
				response := resource.UpdateResponse{State: state}
				r.Update(context.Background(), resource.UpdateRequest{
					State: state, Plan: tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw},
				}, &response)
				require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
				require.Equal(t, 1, patches)
				require.Zero(t, sourceRequests.Load(), "Update must not request the unchanged installer URL")
				var actual MacOSPKGAppResourceModel
				require.False(t, response.State.Get(context.Background(), &actual).HasError())
				require.True(t, actual.MacOSPkgApp.IgnoreVersionDetection.ValueBool())
				require.Equal(t, "Updated app metadata", actual.Description.ValueString())
				require.True(t, actual.AppInstaller.Equal(plannedInstaller(t, planned)))
			})
		})
	}
}

func plannedInstaller(t *testing.T, state tfsdk.State) types.Object {
	t.Helper()
	var model MacOSPKGAppResourceModel
	require.False(t, state.Get(context.Background(), &model).HasError())
	return model.AppInstaller
}
