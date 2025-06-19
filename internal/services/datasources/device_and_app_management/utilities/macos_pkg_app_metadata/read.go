package utilityMacOSPKGAppMetadata

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers/macos_pkg/xar"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read fetches the data from the PKG file and sets it in the data source state
func (d *MacOSPKGAppMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MacOSPKGAppMetadataDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	// Get the configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filePathProvided := !config.InstallerFilePathSource.IsNull() && config.InstallerFilePathSource.ValueString() != ""
	urlProvided := !config.InstallerURLSource.IsNull() && config.InstallerURLSource.ValueString() != ""

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with file path provided: %t, URL provided: %t",
		d.ProviderTypeName, d.TypeName, filePathProvided, urlProvided))

	// Validate inputs - must have either a file path or URL, but not both
	if !filePathProvided && !urlProvided {
		resp.Diagnostics.AddError(
			"Missing Input Source",
			"Either installer_file_path_source or installer_url_source must be provided",
		)
		return
	}

	if filePathProvided && urlProvided {
		resp.Diagnostics.AddError(
			"Multiple Input Sources",
			"Only one of installer_file_path_source or installer_url_source can be provided",
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, config.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var metadata *xar.InstallerMetadata
	var md5Checksum, sha256Checksum []byte
	var err error

	// Extract metadata from file path or URL
	if filePathProvided {
		filePath := config.InstallerFilePathSource.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Extracting metadata from PKG file: %s", filePath))
		metadata, md5Checksum, sha256Checksum, err = extractMetadataFromFile(ctx, filePath)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Extracting Metadata from File",
				fmt.Sprintf("Unable to extract metadata from PKG file: %s", err),
			)
			return
		}
	} else {
		url := config.InstallerURLSource.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Extracting metadata from PKG file at URL: %s", url))
		metadata, md5Checksum, sha256Checksum, err = extractMetadataFromURL(ctx, url)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Extracting Metadata from URL",
				fmt.Sprintf("Unable to extract metadata from PKG file at URL: %s", err),
			)
			return
		}
	}

	// Create a new state model with the results
	var state MacOSPKGAppMetadataDataSourceModel

	// Copy the configuration values
	state.InstallerFilePathSource = config.InstallerFilePathSource
	state.InstallerURLSource = config.InstallerURLSource
	state.Timeouts = config.Timeouts

	// Calculate size in MB
	sizeMB := bytesToMB(metadata.Size)

	// Create the metadata object directly
	metadataObj := MetadataDataSourceModel{
		CFBundleIdentifier:         types.StringValue(metadata.BundleIdentifier),
		CFBundleShortVersionString: types.StringValue(metadata.Version),
		Name:                       types.StringValue(metadata.Name),
		InstallLocation:            types.StringValue(metadata.InstallLocation),
		MinOSVersion:               types.StringValue(metadata.MinOSVersion),
		SizeMB:                     types.Int32Value(int32(bytesToMB(metadata.Size))),
		SHA256Checksum:             types.StringValue(hex.EncodeToString(sha256Checksum)),
		MD5Checksum:                types.StringValue(hex.EncodeToString(md5Checksum)),
	}

	// Set package IDs list
	packageIDs, diags := types.ListValueFrom(ctx, types.StringType, metadata.PackageIDs)
	if !diags.HasError() {
		metadataObj.PackageIDs = packageIDs
	} else {
		metadataObj.PackageIDs = types.ListValueMust(types.StringType, []attr.Value{})
	}

	// Set app paths list
	appPaths, diags := types.ListValueFrom(ctx, types.StringType, metadata.AppPaths)
	if !diags.HasError() {
		metadataObj.AppPaths = appPaths
	} else {
		metadataObj.AppPaths = types.ListValueMust(types.StringType, []attr.Value{})
	}

	// Set included bundles
	if len(metadata.IncludedBundles) > 0 {
		bundleModels := make([]BundleInfoModel, 0, len(metadata.IncludedBundles))
		for _, bundle := range metadata.IncludedBundles {
			bundleModel := BundleInfoModel{
				BundleID:        types.StringValue(bundle.BundleID),
				Version:         types.StringValue(bundle.Version),
				Path:            types.StringValue(bundle.Path),
				CFBundleVersion: types.StringValue(bundle.CFBundleVersion),
			}
			bundleModels = append(bundleModels, bundleModel)
		}

		// Convert to list value
		bundlesList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"bundle_id":         types.StringType,
				"version":           types.StringType,
				"path":              types.StringType,
				"cf_bundle_version": types.StringType,
			},
		}, bundleModels)

		if !diags.HasError() {
			metadataObj.IncludedBundles = bundlesList
		} else {
			metadataObj.IncludedBundles = types.ListValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"bundle_id":         types.StringType,
						"version":           types.StringType,
						"path":              types.StringType,
						"cf_bundle_version": types.StringType,
					},
				},
				[]attr.Value{},
			)
		}
	} else {
		metadataObj.IncludedBundles = types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"bundle_id":         types.StringType,
					"version":           types.StringType,
					"path":              types.StringType,
					"cf_bundle_version": types.StringType,
				},
			},
			[]attr.Value{},
		)
	}

	// Assign the metadata object to the state
	state.Metadata = &metadataObj

	tflog.Debug(ctx, fmt.Sprintf("Successfully created state model with extracted metadata - Size: %.2f MB, SHA256: %s, MD5: %s",
		sizeMB,
		hex.EncodeToString(sha256Checksum)[:16]+"...", // Log just the first 16 chars of the checksum
		hex.EncodeToString(md5Checksum)[:16]+"..."))

	// Set the state from our built state model
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
