package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// evaluateIfContentVersionUpdateRequired determines if a new content version needs to be created and uploaded
// based on source path changes and metadata differences
func evaluateIfContentVersionUpdateRequired(ctx context.Context, object *MacOSPKGAppResourceModel, state *MacOSPKGAppResourceModel, client *msgraphbetasdk.GraphServiceClient) (bool, string, error) {

	pathChanged := false
	metadataChanged := false
	existingContentVersion := ""

	// Check installer source path changes
	var objectMetadata, stateMetadata sharedmodels.MobileAppMetaDataResourceModel

	if !object.AppInstaller.IsNull() {
		diags := object.AppInstaller.As(ctx, &objectMetadata, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return false, "", fmt.Errorf("failed to parse object metadata: %v", diags)
		}
	}

	if !state.AppInstaller.IsNull() {
		diags := state.AppInstaller.As(ctx, &stateMetadata, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return false, "", fmt.Errorf("failed to parse state metadata: %v", diags)
		}
	}

	// Path changes
	if !objectMetadata.InstallerFilePathSource.Equal(stateMetadata.InstallerFilePathSource) {
		tflog.Debug(ctx, "Installer file path has changed", map[string]interface{}{
			"old_path": stateMetadata.InstallerFilePathSource.ValueString(),
			"new_path": objectMetadata.InstallerFilePathSource.ValueString(),
		})
		pathChanged = true
	}

	if !objectMetadata.InstallerURLSource.Equal(stateMetadata.InstallerURLSource) {
		tflog.Debug(ctx, "Installer URL has changed", map[string]interface{}{
			"old_url": stateMetadata.InstallerURLSource.ValueString(),
			"new_url": objectMetadata.InstallerURLSource.ValueString(),
		})
		pathChanged = true
	}

	// If no installer source is provided at all → skip content evaluation
	if (objectMetadata.InstallerFilePathSource.IsNull() || objectMetadata.InstallerFilePathSource.ValueString() == "") &&
		(objectMetadata.InstallerURLSource.IsNull() || objectMetadata.InstallerURLSource.ValueString() == "") {
		tflog.Debug(ctx, "No installer source provided, skipping content evaluation")
		return false, "", nil
	}

	// Check metadata changes
	if !object.AppInstaller.IsNull() && !state.AppInstaller.IsNull() {
		if objectMetadata.InstallerSizeInBytes.ValueInt64() != stateMetadata.InstallerSizeInBytes.ValueInt64() {
			tflog.Debug(ctx, "Installer size has changed", map[string]interface{}{
				"previous_size": stateMetadata.InstallerSizeInBytes.ValueInt64(),
				"current_size":  objectMetadata.InstallerSizeInBytes.ValueInt64(),
			})
			metadataChanged = true
		}

		if objectMetadata.InstallerSHA256Checksum.ValueString() != stateMetadata.InstallerSHA256Checksum.ValueString() {
			tflog.Debug(ctx, "Installer SHA256 checksum has changed", map[string]interface{}{
				"previous_checksum": stateMetadata.InstallerSHA256Checksum.ValueString(),
				"current_checksum":  objectMetadata.InstallerSHA256Checksum.ValueString(),
			})
			metadataChanged = true
		}
	} else {
		tflog.Debug(ctx, "Metadata is missing from one of the objects, assuming content needs update")
		metadataChanged = true
	}

	// Retrieve existing committed content version (from API)
	resource, err := client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Get(ctx, nil)

	if err == nil {
		macOSPkgApp, ok := resource.(graphmodels.MacOSPkgAppable)
		if ok && macOSPkgApp.GetCommittedContentVersion() != nil &&
			*macOSPkgApp.GetCommittedContentVersion() != "" {
			existingContentVersion = *macOSPkgApp.GetCommittedContentVersion()
			tflog.Debug(ctx, fmt.Sprintf("Found existing committed content version from API: %s", existingContentVersion))
		}
	}

	updateNeeded := pathChanged || metadataChanged

	tflog.Debug(ctx, "Content version update evaluation result", map[string]interface{}{
		"update_needed":            updateNeeded,
		"path_changed":             pathChanged,
		"metadata_changed":         metadataChanged,
		"existing_content_version": existingContentVersion,
	})

	return updateNeeded, existingContentVersion, nil
}
