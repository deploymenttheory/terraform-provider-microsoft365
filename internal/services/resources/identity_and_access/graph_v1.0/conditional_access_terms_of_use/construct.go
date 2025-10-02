package graphConditionalAccessTermsOfUse

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// constructResource uses the Microsoft Graph SDK models directly instead of raw HTTP calls
func constructResource(ctx context.Context, client *msgraphsdk.GraphServiceClient, data *ConditionalAccessTermsOfUseResourceModel, isCreate bool) (models.Agreementable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, client, data, isCreate); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := models.NewAgreement()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphBool(data.IsViewingBeforeAcceptanceRequired, requestBody.SetIsViewingBeforeAcceptanceRequired)
	convert.FrameworkToGraphBool(data.IsPerDeviceAcceptanceRequired, requestBody.SetIsPerDeviceAcceptanceRequired)

	// Handle UserReacceptRequiredFrequency - use months to prevent day->week normalization
	if !data.UserReacceptRequiredFrequency.IsNull() && !data.UserReacceptRequiredFrequency.IsUnknown() {
		frequencyStr := data.UserReacceptRequiredFrequency.ValueString()
		if frequencyStr != "" {
			var isoDuration *serialization.ISODuration
			// Use months instead of days to prevent normalization (days get converted to weeks when months=0)
			// The normalize() function only converts days to weeks when both months=0 AND years=0
			switch frequencyStr {
			case "P365D":
				// 365 days = 12 months + 5 days (approximate, but close enough for the API)
				isoDuration = serialization.NewDuration(0, 12, 5, 0, 0, 0, 0)
			case "P180D":
				// 180 days = 6 months (30 days per month average)
				isoDuration = serialization.NewDuration(0, 6, 0, 0, 0, 0, 0)
			case "P90D":
				// 90 days = 3 months (30 days per month average)
				isoDuration = serialization.NewDuration(0, 3, 0, 0, 0, 0, 0)
			case "P30D":
				// 30 days = 1 month
				isoDuration = serialization.NewDuration(0, 1, 0, 0, 0, 0, 0)
			default:
				// Fallback to parsing if it's not one of our standard values
				if parsed, err := serialization.ParseISODuration(frequencyStr); err == nil {
					isoDuration = parsed
				} else {
					tflog.Warn(ctx, "Failed to parse UserReacceptRequiredFrequency", map[string]any{
						"error": err.Error(),
					})
				}
			}
			if isoDuration != nil {
				requestBody.SetUserReacceptRequiredFrequency(isoDuration)
			}
		}
	}

	if !data.TermsExpiration.IsNull() && !data.TermsExpiration.IsUnknown() {
		var termsExpiration TermsExpirationModel
		diags := data.TermsExpiration.As(ctx, &termsExpiration, basetypes.ObjectAsOptions{})
		if !diags.HasError() {
			termsExpirationObj := models.NewTermsExpiration()

			if err := convert.FrameworkToGraphTimeFromDateOnly(termsExpiration.StartDateTime, termsExpirationObj.SetStartDateTime); err != nil {
				tflog.Warn(ctx, "Failed to convert terms expiration start date time", map[string]any{
					"error": err.Error(),
				})
			}

			// Handle Frequency - use months to prevent day->week normalization
			if !termsExpiration.Frequency.IsNull() && !termsExpiration.Frequency.IsUnknown() {
				frequencyStr := termsExpiration.Frequency.ValueString()
				if frequencyStr != "" {
					var isoDuration *serialization.ISODuration
					// Use months instead of days to prevent normalization (days get converted to weeks when months=0)
					// The normalize() function only converts days to weeks when both months=0 AND years=0
					switch frequencyStr {
					case "P365D":
						// 365 days = 12 months + 5 days (approximate, but close enough for the API)
						isoDuration = serialization.NewDuration(0, 12, 5, 0, 0, 0, 0)
					case "P180D":
						// 180 days = 6 months (30 days per month average)
						isoDuration = serialization.NewDuration(0, 6, 0, 0, 0, 0, 0)
					case "P90D":
						// 90 days = 3 months (30 days per month average)
						isoDuration = serialization.NewDuration(0, 3, 0, 0, 0, 0, 0)
					case "P30D":
						// 30 days = 1 month
						isoDuration = serialization.NewDuration(0, 1, 0, 0, 0, 0, 0)
					default:
						// Fallback to parsing if it's not one of our standard values
						if parsed, err := serialization.ParseISODuration(frequencyStr); err == nil {
							isoDuration = parsed
						} else {
							tflog.Warn(ctx, "Failed to parse terms expiration frequency", map[string]any{
								"error": err.Error(),
							})
						}
					}
					if isoDuration != nil {
						termsExpirationObj.SetFrequency(isoDuration)
					}
				}
			}

			requestBody.SetTermsExpiration(termsExpirationObj)
		}
	}

	// Handle file configuration - only for create operations
	if isCreate && !data.File.IsNull() && !data.File.IsUnknown() {
		var file AgreementFileModel
		diags := data.File.As(ctx, &file, basetypes.ObjectAsOptions{})
		if !diags.HasError() {
			agreementFile := models.NewAgreementFile()

			// Handle localizations
			if !file.Localizations.IsNull() && !file.Localizations.IsUnknown() {
				var localizations []AgreementFileLocalizationModel
				diags := file.Localizations.ElementsAs(ctx, &localizations, false)
				if !diags.HasError() {
					agreementLocalizations := make([]models.AgreementFileLocalizationable, len(localizations))
					for i, loc := range localizations {
						localization := models.NewAgreementFileLocalization()
						convert.FrameworkToGraphString(loc.FileName, localization.SetFileName)
						convert.FrameworkToGraphString(loc.DisplayName, localization.SetDisplayName)
						convert.FrameworkToGraphString(loc.Language, localization.SetLanguage)
						convert.FrameworkToGraphBool(loc.IsDefault, localization.SetIsDefault)
						convert.FrameworkToGraphBool(loc.IsMajorVersion, localization.SetIsMajorVersion)

						// Handle file data
						if !loc.FileData.IsNull() && !loc.FileData.IsUnknown() {
							var fileData AgreementFileDataModel
							diags := loc.FileData.As(ctx, &fileData, basetypes.ObjectAsOptions{})
							if !diags.HasError() {
								fileDataObj := models.NewAgreementFileData()
								convert.FrameworkToGraphBytes(fileData.Data, fileDataObj.SetData)
								localization.SetFileData(fileDataObj)
							}
						}

						agreementLocalizations[i] = localization
					}
					agreementFile.SetLocalizations(agreementLocalizations)
				}
			}

			requestBody.SetFile(agreementFile)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource using Graph SDK", ResourceName))

	return requestBody, nil
}
