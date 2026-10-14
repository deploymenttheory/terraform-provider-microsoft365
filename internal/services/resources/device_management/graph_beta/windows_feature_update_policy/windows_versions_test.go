package graphBetaWindowsFeatureUpdatePolicy

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

func TestUnitResourceWindowsFeatureUpdatePolicy_WindowsVersions(t *testing.T) {
	ctx := context.Background()
	var response resource.SchemaResponse
	NewWindowsFeatureUpdatePolicyResource().Schema(ctx, resource.SchemaRequest{}, &response)
	attribute := response.Schema.Attributes["feature_update_version"].(schema.StringAttribute)

	for _, version := range []string{
		"Windows 11, version 26H2",
		"Windows 11, version 25H2",
		"Windows 11, version 24H2",
		"Windows 11, version 23H2",
	} {
		t.Run(version, func(t *testing.T) {
			var validation validator.StringResponse
			for _, v := range attribute.Validators {
				v.ValidateString(ctx, validator.StringRequest{Path: path.Root("feature_update_version"), ConfigValue: types.StringValue(version)}, &validation)
			}
			require.False(t, validation.Diagnostics.HasError(), validation.Diagnostics)

			for _, forUpdate := range []bool{false, true} {
				body, err := constructResource(ctx, &WindowsFeatureUpdatePolicyResourceModel{
					FeatureUpdateVersion: types.StringValue(version),
					RoleScopeTagIds:      types.SetNull(types.StringType),
				}, forUpdate)
				require.NoError(t, err)
				require.NotNil(t, body.GetFeatureUpdateVersion())
				require.Equal(t, version, *body.GetFeatureUpdateVersion())

				var state WindowsFeatureUpdatePolicyResourceModel
				MapRemoteResourceStateToTerraform(ctx, &state, body)
				require.Equal(t, version, state.FeatureUpdateVersion.ValueString())
			}
		})
	}

	for _, version := range []string{"Windows 11, version 27H2", "Windows11_26H2", "26H2", ""} {
		t.Run("reject_"+version, func(t *testing.T) {
			var validation validator.StringResponse
			for _, v := range attribute.Validators {
				v.ValidateString(ctx, validator.StringRequest{Path: path.Root("feature_update_version"), ConfigValue: types.StringValue(version)}, &validation)
			}
			require.True(t, validation.Diagnostics.HasError(), "unsupported version %q was accepted", version)
		})
	}
}
