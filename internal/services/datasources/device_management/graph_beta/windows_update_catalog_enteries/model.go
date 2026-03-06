// REF: https://learn.microsoft.com/en-us/graph/api/windowsupdates-catalog-list-entries?view=graph-rest-beta
package graphBetaWindowsUpdateCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsUpdateCatalogEnteriesDataSourceModel defines the data source model
type WindowsUpdateCatalogEnteriesDataSourceModel struct {
	FilterType  types.String                `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String                `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Entries     []WindowsUpdateCatalogEntry `tfsdk:"entries"`      // List of catalog entries that match the filters
	Timeouts    timeouts.Value              `tfsdk:"timeouts"`
}

// WindowsUpdateCatalogEntry represents a single catalog entry (can be feature or quality update)
type WindowsUpdateCatalogEntry struct {
	ID                      types.String `tfsdk:"id"`
	DisplayName             types.String `tfsdk:"display_name"`
	ReleaseDateTime         types.String `tfsdk:"release_date_time"`
	DeployableUntilDateTime types.String `tfsdk:"deployable_until_date_time"`
	CatalogEntryType        types.String `tfsdk:"catalog_entry_type"`

	// Feature Update specific fields
	Version types.String `tfsdk:"version"`

	// Quality Update specific fields
	CatalogName                 types.String            `tfsdk:"catalog_name"`
	ShortName                   types.String            `tfsdk:"short_name"`
	IsExpeditable               types.Bool              `tfsdk:"is_expeditable"`
	QualityUpdateClassification types.String            `tfsdk:"quality_update_classification"`
	QualityUpdateCadence        types.String            `tfsdk:"quality_update_cadence"`
	CveSeverityInformation      *CveSeverityInformation `tfsdk:"cve_severity_information"`
}

// CveSeverityInformation contains CVE severity details for quality updates
type CveSeverityInformation struct {
	MaxSeverity   types.String   `tfsdk:"max_severity"`
	MaxBaseScore  types.Float64  `tfsdk:"max_base_score"`
	ExploitedCves []ExploitedCve `tfsdk:"exploited_cves"`
}

// ExploitedCve represents a CVE that has been exploited
type ExploitedCve struct {
	Number types.String `tfsdk:"number"`
	Url    types.String `tfsdk:"url"`
}
