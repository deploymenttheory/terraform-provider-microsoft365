package graphBetaWin32App

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceWin32App_WindowsVersions(t *testing.T) {
	ctx := context.Background()
	var response resource.SchemaResponse
	NewWin32LobAppResource().Schema(ctx, resource.SchemaRequest{}, &response)
	attribute := response.Schema.Attributes["minimum_supported_windows_release"].(schema.StringAttribute)

	for _, version := range []string{
		"Windows11_26H2", "Windows11_25H2", "Windows11_24H2", "Windows11_23H2",
		"Windows11_22H2", "Windows11_21H2", "Windows10_22H2", "Windows10_21H2",
		"21H1", "2H20", "2004", "1909", "1903", "1809", "1803", "1709", "1703", "1607",
	} {
		t.Run(version, func(t *testing.T) {
			var validation validator.StringResponse
			for _, v := range attribute.Validators {
				v.ValidateString(ctx, validator.StringRequest{Path: path.Root("minimum_supported_windows_release"), ConfigValue: types.StringValue(version)}, &validation)
			}
			require.False(t, validation.Diagnostics.HasError(), validation.Diagnostics)

			body, err := constructResource(ctx, &Win32LobAppResourceModel{
				MinimumSupportedWindowsRelease: types.StringValue(version),
				RoleScopeTagIds:                types.SetNull(types.StringType),
			}, "installer.intunewin")
			require.NoError(t, err)
			require.NotNil(t, body.GetMinimumSupportedWindowsRelease())
			require.Equal(t, version, *body.GetMinimumSupportedWindowsRelease())

			var state Win32LobAppResourceModel
			MapRemoteResourceStateToTerraform(ctx, &state, body)
			require.Equal(t, version, state.MinimumSupportedWindowsRelease.ValueString())
		})
	}

	for _, version := range []string{"Windows11_27H2", "Windows 11, version 26H2", "26H2", ""} {
		t.Run("reject_"+version, func(t *testing.T) {
			var validation validator.StringResponse
			for _, v := range attribute.Validators {
				v.ValidateString(ctx, validator.StringRequest{Path: path.Root("minimum_supported_windows_release"), ConfigValue: types.StringValue(version)}, &validation)
			}
			require.True(t, validation.Diagnostics.HasError(), "unsupported version %q was accepted", version)
		})
	}
}
