// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileapprelationship?view=graph-rest-beta

package graphBetaMobileAppRelationship

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileAppRelationshipDataSourceModel defines the data source model
type MobileAppRelationshipDataSourceModel struct {
	FilterType   types.String                 `tfsdk:"filter_type"`   // Required field to specify how to filter
	FilterValue  types.String                 `tfsdk:"filter_value"`  // Value to filter by (not used for "all" or "odata")
	ODataFilter  types.String                 `tfsdk:"odata_filter"`  // OData filter parameter
	ODataTop     types.Int32                  `tfsdk:"odata_top"`     // OData top parameter for limiting results
	ODataSkip    types.Int32                  `tfsdk:"odata_skip"`    // OData skip parameter for pagination
	ODataSelect  types.String                 `tfsdk:"odata_select"`  // OData select parameter for field selection
	ODataOrderBy types.String                 `tfsdk:"odata_orderby"` // OData orderby parameter for sorting
	Items        []MobileAppRelationshipModel `tfsdk:"items"`         // List of mobile app relationships that match the filters
	Timeouts     timeouts.Value               `tfsdk:"timeouts"`
}

// MobileAppRelationshipModel represents a single mobile app relationship
type MobileAppRelationshipModel struct {
	ID                         types.String `tfsdk:"id"`                            // Key of the entity
	TargetID                   types.String `tfsdk:"target_id"`                     // The unique app identifier of the target of the mobile app relationship entity
	TargetDisplayName          types.String `tfsdk:"target_display_name"`           // The display name of the app that is the target of the mobile app relationship entity
	TargetDisplayVersion       types.String `tfsdk:"target_display_version"`        // The display version of the app that is the target of the mobile app relationship entity
	TargetPublisher            types.String `tfsdk:"target_publisher"`              // The publisher of the app that is the target of the mobile app relationship entity
	TargetPublisherDisplayName types.String `tfsdk:"target_publisher_display_name"` // The publisher display name of the app that is the target of the mobile app relationship entity
	SourceID                   types.String `tfsdk:"source_id"`                     // The unique app identifier of the source of the mobile app relationship entity
	SourceDisplayName          types.String `tfsdk:"source_display_name"`           // The display name of the app that is the source of the mobile app relationship entity
	SourceDisplayVersion       types.String `tfsdk:"source_display_version"`        // The display version of the app that is the source of the mobile app relationship entity
	SourcePublisherDisplayName types.String `tfsdk:"source_publisher_display_name"` // The publisher display name of the app that is the source of the mobile app relationship entity
	TargetType                 types.String `tfsdk:"target_type"`                   // The type of relationship indicating whether the target application of a relationship is a parent or child
}
