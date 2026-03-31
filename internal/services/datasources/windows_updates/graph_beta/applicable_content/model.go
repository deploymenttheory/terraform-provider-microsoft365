package graphBetaWindowsUpdatesApplicableContent

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ApplicableContentDataSourceModel struct {
	AudienceId         types.String        `tfsdk:"audience_id"`
	CatalogEntryType   types.String        `tfsdk:"catalog_entry_type"`
	DriverClass        types.String        `tfsdk:"driver_class"`
	Manufacturer       types.String        `tfsdk:"manufacturer"`
	DeviceId           types.String        `tfsdk:"device_id"`
	IncludeNoMatches   types.Bool          `tfsdk:"include_no_matches"`
	ODataFilter        types.String        `tfsdk:"odata_filter"`
	ApplicableContent  []ApplicableContent `tfsdk:"applicable_content"`
	Timeouts           timeouts.Value      `tfsdk:"timeouts"`
}

type ApplicableContent struct {
	CatalogEntryId types.String      `tfsdk:"catalog_entry_id"`
	CatalogEntry   *CatalogEntry     `tfsdk:"catalog_entry"`
	MatchedDevices []MatchedDevice   `tfsdk:"matched_devices"`
}

type CatalogEntry struct {
	ID                      types.String `tfsdk:"id"`
	DisplayName             types.String `tfsdk:"display_name"`
	ReleaseDateTime         types.String `tfsdk:"release_date_time"`
	DeployableUntilDateTime types.String `tfsdk:"deployable_until_date_time"`
	Description             types.String `tfsdk:"description"`
	DriverClass             types.String `tfsdk:"driver_class"`
	Provider                types.String `tfsdk:"provider"`
	Manufacturer            types.String `tfsdk:"manufacturer"`
	Version                 types.String `tfsdk:"version"`
	VersionDateTime         types.String `tfsdk:"version_date_time"`
}

type MatchedDevice struct {
	DeviceId      types.String   `tfsdk:"device_id"`
	RecommendedBy []types.String `tfsdk:"recommended_by"`
}
