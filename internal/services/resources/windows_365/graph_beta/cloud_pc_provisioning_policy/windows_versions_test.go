package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceCloudPcProvisioningPolicy_WindowsVersions(t *testing.T) {
	ctx := context.Background()
	var response resource.SchemaResponse
	NewCloudPcProvisioningPolicyResource().Schema(ctx, resource.SchemaRequest{}, &response)
	attribute := response.Schema.Attributes["image_id"].(schema.StringAttribute)

	for _, tc := range []struct {
		imageID     string
		displayName string
	}{
		{"microsoftwindowsdesktop_windows-ent-cpc_win11-26h2-ent-cpc", "Windows 11 Enterprise 26H2"},
		{"microsoftwindowsdesktop_windows-ent-cpc_win11-26h2-ent-cpc-m365", "Windows 11 Enterprise + Microsoft 365 Apps 26H2"},
		{"microsoftwindowsdesktop_windows-ent-cpc_win11-25h2-ent-cpc", "Windows 11 Enterprise 25H2"},
		{"microsoftwindowsdesktop_windows-ent-cpc_win11-25h2-ent-cpc-m365", "Windows 11 Enterprise + Microsoft 365 Apps 25H2"},
		{"microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc", "Windows 11 Enterprise 24H2"},
		{"microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365", "Windows 11 Enterprise + Microsoft 365 Apps 24H2"},
		{"microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc", "Windows 11 Enterprise 23H2"},
		{"microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc-m365", "Windows 11 Enterprise + Microsoft 365 Apps 23H2"},
	} {
		t.Run(tc.imageID, func(t *testing.T) {
			var validation validator.StringResponse
			for _, v := range attribute.Validators {
				v.ValidateString(ctx, validator.StringRequest{Path: path.Root("image_id"), ConfigValue: types.StringValue(tc.imageID)}, &validation)
			}
			require.False(t, validation.Diagnostics.HasError(), validation.Diagnostics)

			// An existing policy skips the name uniqueness request; construction requires no network calls.
			body, err := constructResource(ctx, &CloudPcProvisioningPolicyResourceModel{
				ID:          types.StringValue("existing-policy"),
				DisplayName: types.StringValue("Windows version test"),
				ImageId:     types.StringValue(tc.imageID),
				ScopeIds:    types.SetNull(types.StringType),
			}, &msgraphbetasdk.GraphServiceClient{})
			require.NoError(t, err)
			require.NotNil(t, body.GetImageId())
			require.NotNil(t, body.GetImageDisplayName())
			require.Equal(t, tc.imageID, *body.GetImageId())
			require.Equal(t, tc.displayName, *body.GetImageDisplayName())

			var state CloudPcProvisioningPolicyResourceModel
			MapRemoteStateToTerraform(ctx, &state, body)
			require.Equal(t, tc.imageID, state.ImageId.ValueString())
			require.Equal(t, tc.displayName, state.ImageDisplayName.ValueString())
		})
	}

	for _, imageID := range []string{"microsoftwindowsdesktop_windows-ent-cpc_win11-27h2-ent-cpc", "Windows 11 Enterprise 26H2", ""} {
		t.Run("reject_"+imageID, func(t *testing.T) {
			var validation validator.StringResponse
			for _, v := range attribute.Validators {
				v.ValidateString(ctx, validator.StringRequest{Path: path.Root("image_id"), ConfigValue: types.StringValue(imageID)}, &validation)
			}
			require.True(t, validation.Diagnostics.HasError(), "unsupported image %q was accepted", imageID)
		})
	}
}
