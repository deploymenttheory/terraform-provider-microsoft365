package sharedStater

import (
	"context"
	"testing"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestMapAppMetadataStateToTerraform tests the MapAppMetadataStateToTerraform function
func TestMapAppMetadataStateToTerraform(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		metadata *sharedmodels.MobileAppMetaDataResourceModel
		validate func(t *testing.T, result types.Object)
	}{
		{
			name:     "Nil metadata",
			metadata: nil,
			validate: func(t *testing.T, result types.Object) {
				assert.True(t, result.IsNull())
			},
		},
		{
			name: "Metadata with both fields populated",
			metadata: &sharedmodels.MobileAppMetaDataResourceModel{
				InstallerFilePathSource: types.StringValue("/path/to/installer.msi"),
				InstallerURLSource:      types.StringValue("https://example.com/installer.msi"),
			},
			validate: func(t *testing.T, result types.Object) {
				assert.False(t, result.IsNull())
				assert.False(t, result.IsUnknown())
				
				attrs := result.Attributes()
				assert.Contains(t, attrs, "installer_file_path_source")
				assert.Contains(t, attrs, "installer_url_source")
			},
		},
		{
			name: "Metadata with only file path",
			metadata: &sharedmodels.MobileAppMetaDataResourceModel{
				InstallerFilePathSource: types.StringValue("/path/to/installer.msi"),
				InstallerURLSource:      types.StringNull(),
			},
			validate: func(t *testing.T, result types.Object) {
				assert.False(t, result.IsNull())
				attrs := result.Attributes()
				filePathVal := attrs["installer_file_path_source"].(types.String)
				assert.Equal(t, "/path/to/installer.msi", filePathVal.ValueString())
			},
		},
		{
			name: "Metadata with only URL",
			metadata: &sharedmodels.MobileAppMetaDataResourceModel{
				InstallerFilePathSource: types.StringNull(),
				InstallerURLSource:      types.StringValue("https://example.com/installer.msi"),
			},
			validate: func(t *testing.T, result types.Object) {
				assert.False(t, result.IsNull())
				attrs := result.Attributes()
				urlVal := attrs["installer_url_source"].(types.String)
				assert.Equal(t, "https://example.com/installer.msi", urlVal.ValueString())
			},
		},
		{
			name: "Metadata with empty strings",
			metadata: &sharedmodels.MobileAppMetaDataResourceModel{
				InstallerFilePathSource: types.StringValue(""),
				InstallerURLSource:      types.StringValue(""),
			},
			validate: func(t *testing.T, result types.Object) {
				assert.False(t, result.IsNull())
				attrs := result.Attributes()
				assert.Contains(t, attrs, "installer_file_path_source")
				assert.Contains(t, attrs, "installer_url_source")
			},
		},
		{
			name: "Metadata with null values",
			metadata: &sharedmodels.MobileAppMetaDataResourceModel{
				InstallerFilePathSource: types.StringNull(),
				InstallerURLSource:      types.StringNull(),
			},
			validate: func(t *testing.T, result types.Object) {
				assert.False(t, result.IsNull())
			},
		},
		{
			name: "Metadata with unknown values",
			metadata: &sharedmodels.MobileAppMetaDataResourceModel{
				InstallerFilePathSource: types.StringUnknown(),
				InstallerURLSource:      types.StringUnknown(),
			},
			validate: func(t *testing.T, result types.Object) {
				assert.False(t, result.IsNull())
				// Should handle unknown values without error
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapAppMetadataStateToTerraform(ctx, tt.metadata)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
