// REF: https://learn.microsoft.com/en-us/graph/api/resources/agreement?view=graph-rest-1.0
package graphConditionalAccessTermsOfUse

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConditionalAccessTermsOfUseResourceModel represents the terms of use agreement resource model
type ConditionalAccessTermsOfUseResourceModel struct {
	ID                                types.String   `tfsdk:"id"`
	DisplayName                       types.String   `tfsdk:"display_name"`
	IsViewingBeforeAcceptanceRequired types.Bool     `tfsdk:"is_viewing_before_acceptance_required"`
	IsPerDeviceAcceptanceRequired     types.Bool     `tfsdk:"is_per_device_acceptance_required"`
	UserReacceptRequiredFrequency     types.String   `tfsdk:"user_reaccept_required_frequency"`
	TermsExpiration                   types.Object   `tfsdk:"terms_expiration"`
	File                              types.Object   `tfsdk:"file"`
	Timeouts                          timeouts.Value `tfsdk:"timeouts"`
}

// TermsExpirationModel represents the terms expiration configuration
type TermsExpirationModel struct {
	StartDateTime types.String `tfsdk:"start_date_time"`
	Frequency     types.String `tfsdk:"frequency"`
}

// AgreementFileModel represents the agreement file configuration
type AgreementFileModel struct {
	Localizations types.Set `tfsdk:"localizations"`
}

// AgreementFileLocalizationModel represents a file localization
type AgreementFileLocalizationModel struct {
	FileName       types.String `tfsdk:"file_name"`
	DisplayName    types.String `tfsdk:"display_name"`
	Language       types.String `tfsdk:"language"`
	IsDefault      types.Bool   `tfsdk:"is_default"`
	IsMajorVersion types.Bool   `tfsdk:"is_major_version"`
	FileData       types.Object `tfsdk:"file_data"`
}

// AgreementFileDataModel represents the file data
type AgreementFileDataModel struct {
	Data types.String `tfsdk:"data"`
}
