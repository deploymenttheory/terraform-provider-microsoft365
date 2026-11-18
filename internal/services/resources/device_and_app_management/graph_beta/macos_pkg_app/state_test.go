package graphBetaMacOSPKGApp

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceMacOSPkgApp_06_ImportedOptionalURLs(t *testing.T) {
	for _, test := range []struct {
		name                         string
		privacy, information         any
		wantPrivacy, wantInformation string
	}{
		{name: "null"},
		{name: "empty privacy", privacy: ""},
		{name: "valid", privacy: "https://example.com/privacy?source=intune", information: "http://example.com/app",
			wantPrivacy: "https://example.com/privacy?source=intune", wantInformation: "http://example.com/app"},
		{name: "empty privacy with valid information", privacy: "", information: "https://example.com/app", wantInformation: "https://example.com/app"},
		{name: "valid privacy with null information", privacy: "https://example.com/privacy", wantPrivacy: "https://example.com/privacy"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			transport := httpmock.NewMockTransport()
			body := map[string]any{
				"@odata.type": "#microsoft.graph.macOSPkgApp",
				"id":          "11111111-1111-1111-1111-111111111111",
				"displayName": "unit-test-macos-pkg-url-import", "description": "URL import regression", "publisher": "MacAdmins",
				"isFeatured": false, "privacyInformationUrl": test.privacy, "informationUrl": test.information,
				"primaryBundleId": "com.github.macadmins.Nudge", "primaryBundleVersion": "2.1.4.81863",
				"ignoreVersionDetection": false, "minimumSupportedOperatingSystem": map[string]bool{"v12_0": true},
			}
			responder, err := httpmock.NewJsonResponder(http.StatusOK, body)
			require.NoError(t, err)
			transport.RegisterRegexpResponder("GET", regexp.MustCompile("^https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/11111111-1111-1111-1111-111111111111"), responder)
			r := NewMacOSPKGAppResource().(*MacOSPKGAppResource)
			r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
			var schema resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schema)
			imported := resource.ImportStateResponse{State: tfsdk.State{
				Schema: schema.Schema, Raw: tftypes.NewValue(schema.Schema.Type().TerraformType(ctx), nil),
			}}
			r.ImportState(ctx, resource.ImportStateRequest{ID: "11111111-1111-1111-1111-111111111111"}, &imported)
			require.False(t, imported.Diagnostics.HasError(), "%v", imported.Diagnostics)
			response := resource.ReadResponse{State: imported.State}
			r.Read(ctx, resource.ReadRequest{State: imported.State}, &response)
			require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
			var actual MacOSPKGAppResourceModel
			require.False(t, response.State.Get(ctx, &actual).HasError())
			require.Equal(t, test.wantPrivacy == "", actual.PrivacyInformationUrl.IsNull())
			require.Equal(t, test.wantInformation == "", actual.InformationUrl.IsNull())
			require.Equal(t, test.wantPrivacy, actual.PrivacyInformationUrl.ValueString())
			require.Equal(t, test.wantInformation, actual.InformationUrl.ValueString())
			require.True(t, actual.AppInstaller.IsNull(), "import must not invent an installer source")
		})
	}
}
