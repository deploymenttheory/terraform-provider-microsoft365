package graphConditionalAccessTermsOfUse

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructResource constructs a map representing the Agreement resource for API calls
func constructResource(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, data *ConditionalAccessTermsOfUseResourceModel, isCreate bool) (map[string]any, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, data); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := make(map[string]any)

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		requestBody["displayName"] = data.DisplayName.ValueString()
	}

	if !data.IsViewingBeforeAcceptanceRequired.IsNull() && !data.IsViewingBeforeAcceptanceRequired.IsUnknown() {
		requestBody["isViewingBeforeAcceptanceRequired"] = data.IsViewingBeforeAcceptanceRequired.ValueBool()
	}

	if !data.IsPerDeviceAcceptanceRequired.IsNull() && !data.IsPerDeviceAcceptanceRequired.IsUnknown() {
		requestBody["isPerDeviceAcceptanceRequired"] = data.IsPerDeviceAcceptanceRequired.ValueBool()
	}

	// Handle UserReacceptRequiredFrequency
	if !data.UserReacceptRequiredFrequency.IsNull() && !data.UserReacceptRequiredFrequency.IsUnknown() {
		frequencyStr := data.UserReacceptRequiredFrequency.ValueString()
		if frequencyStr != "" {
			requestBody["userReacceptRequiredFrequency"] = frequencyStr
		}
	}

	// Handle terms expiration
	if !data.TermsExpiration.IsNull() && !data.TermsExpiration.IsUnknown() {
		var termsExpiration TermsExpirationModel
		diags := data.TermsExpiration.As(ctx, &termsExpiration, basetypes.ObjectAsOptions{})
		if !diags.HasError() {
			termsExpirationObj := make(map[string]any)

			if !termsExpiration.StartDateTime.IsNull() && !termsExpiration.StartDateTime.IsUnknown() {
				startDateTime := termsExpiration.StartDateTime.ValueString()
				// If the datetime doesn't contain 'T', append T00:00:00Z for ISO 8601 format
				if !strings.Contains(startDateTime, "T") {
					startDateTime = startDateTime + "T00:00:00Z"
				}
				termsExpirationObj["startDateTime"] = startDateTime
			}

			if !termsExpiration.Frequency.IsNull() && !termsExpiration.Frequency.IsUnknown() {
				frequencyStr := termsExpiration.Frequency.ValueString()
				if frequencyStr != "" {
					termsExpirationObj["frequency"] = frequencyStr
				}
			}

			if len(termsExpirationObj) > 0 {
				requestBody["termsExpiration"] = termsExpirationObj
			}
		}
	}

	// Handle file configuration - only for create operations
	if isCreate && !data.File.IsNull() && !data.File.IsUnknown() {
		var file AgreementFileModel
		diags := data.File.As(ctx, &file, basetypes.ObjectAsOptions{})
		if !diags.HasError() {
			agreementFile := make(map[string]any)

			// Handle localizations
			if !file.Localizations.IsNull() && !file.Localizations.IsUnknown() {
				var localizations []AgreementFileLocalizationModel
				diags := file.Localizations.ElementsAs(ctx, &localizations, false)
				if !diags.HasError() {
					agreementLocalizations := make([]map[string]any, len(localizations))
					for i, loc := range localizations {
						localization := make(map[string]any)

						if !loc.FileName.IsNull() && !loc.FileName.IsUnknown() {
							localization["fileName"] = loc.FileName.ValueString()
						}

						if !loc.DisplayName.IsNull() && !loc.DisplayName.IsUnknown() {
							localization["displayName"] = loc.DisplayName.ValueString()
						}

						if !loc.Language.IsNull() && !loc.Language.IsUnknown() {
							localization["language"] = loc.Language.ValueString()
						}

						if !loc.IsDefault.IsNull() && !loc.IsDefault.IsUnknown() {
							localization["isDefault"] = loc.IsDefault.ValueBool()
						}

						if !loc.IsMajorVersion.IsNull() && !loc.IsMajorVersion.IsUnknown() {
							localization["isMajorVersion"] = loc.IsMajorVersion.ValueBool()
						}

						// Handle file data
						if !loc.FileData.IsNull() && !loc.FileData.IsUnknown() {
							var fileData AgreementFileDataModel
							diags := loc.FileData.As(ctx, &fileData, basetypes.ObjectAsOptions{})
							if !diags.HasError() {
								fileDataObj := make(map[string]any)

								if !fileData.Data.IsNull() && !fileData.Data.IsUnknown() {
									// Data should already be base64 encoded
									dataStr := fileData.Data.ValueString()
									// Validate it's proper base64
									if _, err := base64.StdEncoding.DecodeString(dataStr); err != nil {
										tflog.Warn(ctx, "File data is not valid base64", map[string]any{
											"error": err.Error(),
										})
									}
									fileDataObj["data"] = dataStr
								}

								if len(fileDataObj) > 0 {
									localization["fileData"] = fileDataObj
								}
							}
						}

						agreementLocalizations[i] = localization
					}
					agreementFile["localizations"] = agreementLocalizations
				}
			}

			if len(agreementFile) > 0 {
				requestBody["file"] = agreementFile
			}
		}
	}

	// Debug log the final JSON that will be sent to the API
	if debugJSON, err := json.MarshalIndent(requestBody, "", "    "); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), map[string]any{
			"json": "\n" + string(debugJSON),
		})
	} else {
		tflog.Error(ctx, "Failed to marshal request body for debug logging", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
