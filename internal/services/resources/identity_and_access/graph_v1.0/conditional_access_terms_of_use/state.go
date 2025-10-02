package graphConditionalAccessTermsOfUse

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the Graph API model into the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ConditionalAccessTermsOfUseResourceModel, remoteResource models.Agreementable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.IsViewingBeforeAcceptanceRequired = convert.GraphToFrameworkBool(remoteResource.GetIsViewingBeforeAcceptanceRequired())
	data.IsPerDeviceAcceptanceRequired = convert.GraphToFrameworkBool(remoteResource.GetIsPerDeviceAcceptanceRequired())

	if remoteResource.GetUserReacceptRequiredFrequency() != nil {
		isoValue := convert.GraphToFrameworkISODuration(remoteResource.GetUserReacceptRequiredFrequency())
		if !isoValue.IsNull() {
			// Handle server-side normalization by converting back to expected format
			normalizedValue := isoValue.ValueString()
			switch normalizedValue {
			case "P52W", "P1Y", "P365DT0.001S", "P12M5D":
				data.UserReacceptRequiredFrequency = types.StringValue("P365D")
			case "P12W", "P3M", "P2M30D", "P90DT0.001S":
				data.UserReacceptRequiredFrequency = types.StringValue("P90D")
			case "P26W", "P6M", "P180DT0.001S":
				data.UserReacceptRequiredFrequency = types.StringValue("P180D")
			case "P4W", "P1M", "P30DT0.001S":
				data.UserReacceptRequiredFrequency = types.StringValue("P30D")
			default:
				data.UserReacceptRequiredFrequency = isoValue
			}
		} else {
			data.UserReacceptRequiredFrequency = types.StringNull()
		}
	} else {
		data.UserReacceptRequiredFrequency = types.StringNull()
	}

	if termsExpiration := remoteResource.GetTermsExpiration(); termsExpiration != nil {
		var startDateTimeValue attr.Value = types.StringNull()
		var frequencyValue attr.Value = types.StringNull()

		startDateTimeValue = convert.GraphToFrameworkTimeAsDateOnly(termsExpiration.GetStartDateTime())

		// Handle frequency with server-side normalization conversion
		if termsExpiration.GetFrequency() != nil {
			isoValue := convert.GraphToFrameworkISODuration(termsExpiration.GetFrequency())
			if !isoValue.IsNull() {
				normalizedValue := isoValue.ValueString()
				switch normalizedValue {
				case "P52W", "P1Y", "P365DT0.001S", "P12M5D":
					frequencyValue = types.StringValue("P365D")
				case "P12W", "P3M", "P2M30D", "P90DT0.001S":
					frequencyValue = types.StringValue("P90D")
				case "P26W", "P6M", "P180DT0.001S":
					frequencyValue = types.StringValue("P180D")
				case "P4W", "P1M", "P30DT0.001S":
					frequencyValue = types.StringValue("P30D")
				default:
					frequencyValue = isoValue
				}
			}
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

	// Handle file configuration
	if file := remoteResource.GetFile(); file != nil {
		fileAttrs := map[string]attr.Value{}

		// Handle localizations
		if localizations := file.GetLocalizations(); len(localizations) > 0 {
			localizationElements := make([]attr.Value, len(localizations))

			for i, loc := range localizations {
				var fileDataValue attr.Value = types.StringNull()
				if fileData := loc.GetFileData(); fileData != nil {
					fileDataValue = convert.GraphToFrameworkBytes(fileData.GetData())
				}

				fileDataAttrs := map[string]attr.Value{
					"data": fileDataValue,
				}

				fileDataObj, diags := types.ObjectValue(map[string]attr.Type{
					"data": types.StringType,
				}, fileDataAttrs)

				if diags.HasError() {
					continue
				}

				localizationAttrs := map[string]attr.Value{
					"file_name":        convert.GraphToFrameworkString(loc.GetFileName()),
					"display_name":     convert.GraphToFrameworkString(loc.GetDisplayName()),
					"language":         convert.GraphToFrameworkString(loc.GetLanguage()),
					"is_default":       convert.GraphToFrameworkBool(loc.GetIsDefault()),
					"is_major_version": convert.GraphToFrameworkBool(loc.GetIsMajorVersion()),
					"file_data":        fileDataObj,
				}

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
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
