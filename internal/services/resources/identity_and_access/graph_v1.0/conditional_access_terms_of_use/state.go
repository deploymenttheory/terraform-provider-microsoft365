package graphConditionalAccessTermsOfUse

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteResourceStateToTerraform maps the Graph API response map into the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ConditionalAccessTermsOfUseResourceModel, remoteResource map[string]any) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName":       remoteResource["displayName"],
		"resourceId":         remoteResource["id"],
		"hasFilesArray":      remoteResource["files"] != nil,
		"remoteResourceKeys": fmt.Sprintf("%v", getMapKeys(remoteResource)),
	})

	// Preserve existing file_data from prior state for matching purposes
	priorFileDataMap := extractPriorFileData(ctx, data)

	// Map ID
	if id, ok := remoteResource["id"].(string); ok {
		data.ID = types.StringValue(id)
	}

	// Map DisplayName
	if displayName, ok := remoteResource["displayName"].(string); ok {
		data.DisplayName = types.StringValue(displayName)
	}

	// Map IsViewingBeforeAcceptanceRequired
	if isViewingBeforeAcceptanceRequired, ok := remoteResource["isViewingBeforeAcceptanceRequired"].(bool); ok {
		data.IsViewingBeforeAcceptanceRequired = types.BoolValue(isViewingBeforeAcceptanceRequired)
	}

	// Map IsPerDeviceAcceptanceRequired
	if isPerDeviceAcceptanceRequired, ok := remoteResource["isPerDeviceAcceptanceRequired"].(bool); ok {
		data.IsPerDeviceAcceptanceRequired = types.BoolValue(isPerDeviceAcceptanceRequired)
	}

	// Map UserReacceptRequiredFrequency
	if userReacceptRequiredFrequency, ok := remoteResource["userReacceptRequiredFrequency"].(string); ok && userReacceptRequiredFrequency != "" {
		data.UserReacceptRequiredFrequency = types.StringValue(userReacceptRequiredFrequency)
	} else {
		data.UserReacceptRequiredFrequency = types.StringNull()
	}

	// Map TermsExpiration
	if termsExpirationRaw, ok := remoteResource["termsExpiration"].(map[string]any); ok {
		var startDateTimeValue attr.Value = types.StringNull()
		var frequencyValue attr.Value = types.StringNull()

		if startDateTime, ok := termsExpirationRaw["startDateTime"].(string); ok {
			// Strip time portion to return just the date (YYYY-MM-DD)
			// API returns "2025-11-06T00:00:00Z", we store "2025-11-06"
			dateOnly := stripTimeFromISO8601DateTime(startDateTime)
			startDateTimeValue = types.StringValue(dateOnly)
		}

		if frequency, ok := termsExpirationRaw["frequency"].(string); ok && frequency != "" {
			frequencyValue = types.StringValue(frequency)
		}

		termsExpirationAttrs := map[string]attr.Value{
			"start_date_time": startDateTimeValue,
			"frequency":       frequencyValue,
		}

		termsExpirationObj, diags := types.ObjectValue(map[string]attr.Type{
			"start_date_time": types.StringType,
			"frequency":       types.StringType,
		}, termsExpirationAttrs)

		if !diags.HasError() {
			data.TermsExpiration = termsExpirationObj
		}
	} else {
		data.TermsExpiration = types.ObjectNull(map[string]attr.Type{
			"start_date_time": types.StringType,
			"frequency":       types.StringType,
		})
	}

	// Map File configuration
	// The API returns "files" array at the root level, not nested under "file"
	if filesRaw, ok := remoteResource["files"].([]any); ok && len(filesRaw) > 0 {
		tflog.Debug(ctx, "Mapping files from API response", map[string]any{
			"filesCount": len(filesRaw),
		})

		fileAttrs := map[string]attr.Value{}

		// Handle localizations (files array from API)
		localizationsRaw := filesRaw
		if len(localizationsRaw) > 0 {
			localizationElements := make([]attr.Value, len(localizationsRaw))

			for i, locRaw := range localizationsRaw {
				if loc, ok := locRaw.(map[string]any); ok {
					tflog.Debug(ctx, "Processing localization", map[string]any{
						"index":            i,
						"fileName":         loc["fileName"],
						"language":         loc["language"],
						"isDefault":        loc["isDefault"],
						"isMajorVersion":   loc["isMajorVersion"],
						"localizationKeys": fmt.Sprintf("%v", getMapKeys(loc)),
					})

					// Try to preserve file_data from prior state by matching on fileName
					fileDataValue := types.ObjectNull(map[string]attr.Type{"data": types.StringType})
					if fileName, ok := loc["fileName"].(string); ok {
						if priorData, exists := priorFileDataMap[fileName]; exists {
							fileDataValue = priorData
							tflog.Debug(ctx, "Preserved file_data from prior state", map[string]any{
								"fileName": fileName,
							})
						}
					}

					localizationAttrs := map[string]attr.Value{
						"file_name":        getStringValue(loc, "fileName"),
						"display_name":     getStringValue(loc, "displayName"),
						"language":         getStringValue(loc, "language"),
						"is_default":       getBoolValue(loc, "isDefault"),
						"is_major_version": getBoolValue(loc, "isMajorVersion"),
						"file_data":        fileDataValue,
					}

					tflog.Debug(ctx, "Mapped localization attributes", map[string]any{
						"index":            i,
						"file_name":        localizationAttrs["file_name"],
						"language":         localizationAttrs["language"],
						"is_default":       localizationAttrs["is_default"],
						"is_major_version": localizationAttrs["is_major_version"],
					})

					localizationObj, diags := types.ObjectValue(map[string]attr.Type{
						"file_name":        types.StringType,
						"display_name":     types.StringType,
						"language":         types.StringType,
						"is_default":       types.BoolType,
						"is_major_version": types.BoolType,
						"file_data":        types.ObjectType{AttrTypes: map[string]attr.Type{"data": types.StringType}},
					}, localizationAttrs)

					if !diags.HasError() {
						localizationElements[i] = localizationObj
					}
				}
			}

			localizationsSet, diags := types.SetValue(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"file_name":        types.StringType,
					"display_name":     types.StringType,
					"language":         types.StringType,
					"is_default":       types.BoolType,
					"is_major_version": types.BoolType,
					"file_data":        types.ObjectType{AttrTypes: map[string]attr.Type{"data": types.StringType}},
				},
			}, localizationElements)

			if !diags.HasError() {
				fileAttrs["localizations"] = localizationsSet
			}
		}

		if len(fileAttrs) > 0 {
			fileObj, diags := types.ObjectValue(map[string]attr.Type{
				"localizations": types.SetType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"file_name":        types.StringType,
							"display_name":     types.StringType,
							"language":         types.StringType,
							"is_default":       types.BoolType,
							"is_major_version": types.BoolType,
							"file_data":        types.ObjectType{AttrTypes: map[string]attr.Type{"data": types.StringType}},
						},
					},
				},
			}, fileAttrs)

			if !diags.HasError() {
				data.File = fileObj
				tflog.Debug(ctx, "Successfully mapped file configuration to state")
			} else {
				tflog.Error(ctx, "Error mapping file configuration", map[string]any{
					"diagnostics": diags,
				})
			}
		}
	} else {
		tflog.Warn(ctx, "No files array found in API response or files array is empty")
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// Helper functions to safely extract values from maps
func getStringValue(m map[string]any, key string) types.String {
	if val, ok := m[key].(string); ok {
		return types.StringValue(val)
	}
	return types.StringNull()
}

func getBoolValue(m map[string]any, key string) types.Bool {
	if val, ok := m[key].(bool); ok {
		return types.BoolValue(val)
	}
	return types.BoolNull()
}

// stripTimeFromISO8601DateTime extracts just the date portion from an ISO 8601 datetime
// Converts "2025-11-06T00:00:00Z" to "2025-11-06"
func stripTimeFromISO8601DateTime(datetime string) string {
	if idx := strings.Index(datetime, "T"); idx > 0 {
		return datetime[:idx]
	}
	return datetime
}

// getMapKeys returns the keys of a map as a slice
func getMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// extractPriorFileData extracts file_data values from the prior state, keyed by fileName
func extractPriorFileData(ctx context.Context, data *ConditionalAccessTermsOfUseResourceModel) map[string]types.Object {
	fileDataMap := make(map[string]types.Object)

	if data.File.IsNull() || data.File.IsUnknown() {
		return fileDataMap
	}

	var file AgreementFileModel
	diags := data.File.As(ctx, &file, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		tflog.Debug(ctx, "Failed to extract file from prior state", map[string]any{
			"diagnostics": diags,
		})
		return fileDataMap
	}

	if file.Localizations.IsNull() || file.Localizations.IsUnknown() {
		return fileDataMap
	}

	var localizations []AgreementFileLocalizationModel
	diags = file.Localizations.ElementsAs(ctx, &localizations, false)
	if diags.HasError() {
		tflog.Debug(ctx, "Failed to extract localizations from prior state", map[string]any{
			"diagnostics": diags,
		})
		return fileDataMap
	}

	for _, loc := range localizations {
		if !loc.FileName.IsNull() && !loc.FileName.IsUnknown() {
			fileName := loc.FileName.ValueString()
			if !loc.FileData.IsNull() && !loc.FileData.IsUnknown() {
				fileDataMap[fileName] = loc.FileData
				tflog.Debug(ctx, "Extracted file_data from prior state", map[string]any{
					"fileName": fileName,
				})
			}
		}
	}

	return fileDataMap
}
