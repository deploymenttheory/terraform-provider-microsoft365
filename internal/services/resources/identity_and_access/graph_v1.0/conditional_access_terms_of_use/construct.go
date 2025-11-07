package graphConditionalAccessTermsOfUse

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructResource constructs a map representing the Agreement resource for API calls
func constructResource(ctx context.Context, data *ConditionalAccessTermsOfUseResourceModel) (map[string]any, error) {
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

	if !data.UserReacceptRequiredFrequency.IsNull() && !data.UserReacceptRequiredFrequency.IsUnknown() {
		requestBody["userReacceptRequiredFrequency"] = data.UserReacceptRequiredFrequency.ValueString()
	}

	if !data.TermsExpiration.IsNull() && !data.TermsExpiration.IsUnknown() {
		var termsExpiration TermsExpirationModel
		diags := data.TermsExpiration.As(ctx, &termsExpiration, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract terms expiration")
		}

		termsExpirationMap := make(map[string]any)
		if !termsExpiration.StartDateTime.IsNull() && !termsExpiration.StartDateTime.IsUnknown() {
			startDateTime := termsExpiration.StartDateTime.ValueString() + "T00:00:00.000Z"
			termsExpirationMap["startDateTime"] = startDateTime
		}
		if !termsExpiration.Frequency.IsNull() && !termsExpiration.Frequency.IsUnknown() {
			termsExpirationMap["frequency"] = termsExpiration.Frequency.ValueString()
		}

		if len(termsExpirationMap) > 0 {
			requestBody["termsExpiration"] = termsExpirationMap
		}
	}

	if !data.File.IsNull() && !data.File.IsUnknown() {
		var file AgreementFileModel
		diags := data.File.As(ctx, &file, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract file configuration")
		}

		if !file.Localizations.IsNull() && !file.Localizations.IsUnknown() {
			var localizations []AgreementFileLocalizationModel
			diags := file.Localizations.ElementsAs(ctx, &localizations, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract localizations")
			}

			localizationsList := make([]map[string]any, 0, len(localizations))
			for _, loc := range localizations {
				locMap := make(map[string]any)

				if !loc.FileName.IsNull() && !loc.FileName.IsUnknown() {
					locMap["fileName"] = loc.FileName.ValueString()
				}

				if !loc.DisplayName.IsNull() && !loc.DisplayName.IsUnknown() {
					locMap["displayName"] = loc.DisplayName.ValueString()
				}

				if !loc.Language.IsNull() && !loc.Language.IsUnknown() {
					locMap["language"] = loc.Language.ValueString()
				}

				if !loc.IsDefault.IsNull() && !loc.IsDefault.IsUnknown() {
					locMap["isDefault"] = loc.IsDefault.ValueBool()
				}

				if !loc.IsMajorVersion.IsNull() && !loc.IsMajorVersion.IsUnknown() {
					locMap["isMajorVersion"] = loc.IsMajorVersion.ValueBool()
				}

				if !loc.FileData.IsNull() && !loc.FileData.IsUnknown() {
					var fileData AgreementFileDataModel
					diags := loc.FileData.As(ctx, &fileData, basetypes.ObjectAsOptions{})
					if diags.HasError() {
						return nil, fmt.Errorf("failed to extract file data")
					}

					if !fileData.Data.IsNull() && !fileData.Data.IsUnknown() {
						base64Data := fileData.Data.ValueString()
						decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
						if err != nil {
							return nil, fmt.Errorf("failed to decode base64 file data: %w", err)
						}

						locMap["fileData"] = map[string]any{
							"data": decodedBytes,
						}
					}
				}

				localizationsList = append(localizationsList, locMap)
			}

			if len(localizationsList) > 0 {
				requestBody["file"] = map[string]any{
					"localizations": localizationsList,
				}
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
