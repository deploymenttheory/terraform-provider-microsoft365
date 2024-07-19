package assignmentFilter

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// platformValidator is the custom validator type
type platformValidator struct{}

// Validate performs the validation.
func (v platformValidator) Validate(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	_, err := models.ParseDevicePlatformType(req.ConfigValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device Platform Type",
			fmt.Sprintf("The platform type '%s' is not valid. Supported types: %v", req.ConfigValue.ValueString(), getAllPlatformStrings()),
		)
	}
}

// Description describes the validation in plain text.
func (v platformValidator) Description(ctx context.Context) string {
	return "must be a valid device platform type"
}

// MarkdownDescription describes the validation in Markdown.
func (v platformValidator) MarkdownDescription(ctx context.Context) string {
	return "must be a valid device platform type"
}

// getAllPlatformStrings returns all the valid platform strings
func getAllPlatformStrings() []string {
	platformTypes := []models.DevicePlatformType{
		models.ANDROID_DEVICEPLATFORMTYPE,
		models.ANDROIDFORWORK_DEVICEPLATFORMTYPE,
		models.IOS_DEVICEPLATFORMTYPE,
		models.MACOS_DEVICEPLATFORMTYPE,
		models.WINDOWSPHONE81_DEVICEPLATFORMTYPE,
		models.WINDOWS81ANDLATER_DEVICEPLATFORMTYPE,
		models.WINDOWS10ANDLATER_DEVICEPLATFORMTYPE,
		models.ANDROIDWORKPROFILE_DEVICEPLATFORMTYPE,
		models.UNKNOWN_DEVICEPLATFORMTYPE,
		models.ANDROIDAOSP_DEVICEPLATFORMTYPE,
		models.ANDROIDMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE,
		models.IOSMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE,
		models.UNKNOWNFUTUREVALUE_DEVICEPLATFORMTYPE,
		models.WINDOWSMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE,
	}

	var platformStrings []string
	for _, platform := range platformTypes {
		platformStrings = append(platformStrings, platform.String())
	}
	return platformStrings
}
