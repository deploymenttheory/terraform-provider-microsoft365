// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileapprelationship?view=graph-rest-beta

package graphBetaMobileAppRelationship

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileAppRelationshipDataSourceModel defines the data source model.
type MobileAppRelationshipDataSourceModel struct {
	ID                         types.String   `tfsdk:"id"`
	SourceID                   types.String   `tfsdk:"source_id"`
	RelationshipID             types.String   `tfsdk:"relationship_id"`
	TargetID                   types.String   `tfsdk:"target_id"`
	TargetDisplayName          types.String   `tfsdk:"target_display_name"`
	ODataQuery                 types.String   `tfsdk:"odata_query"`
	TargetType                 types.String   `tfsdk:"target_type"`
	RelationshipType           types.String   `tfsdk:"relationship_type"`
	TargetDisplayVersion       types.String   `tfsdk:"target_display_version"`
	TargetPublisher            types.String   `tfsdk:"target_publisher"`
	TargetPublisherDisplayName types.String   `tfsdk:"target_publisher_display_name"`
	SourceDisplayName          types.String   `tfsdk:"source_display_name"`
	SourceDisplayVersion       types.String   `tfsdk:"source_display_version"`
	SourcePublisherDisplayName types.String   `tfsdk:"source_publisher_display_name"`
	DependencyType             types.String   `tfsdk:"dependency_type"`
	DependentAppCount          types.Int32    `tfsdk:"dependent_app_count"`
	DependsOnAppCount          types.Int32    `tfsdk:"depends_on_app_count"`
	SupersedenceType           types.String   `tfsdk:"supersedence_type"`
	SupersededAppCount         types.Int32    `tfsdk:"superseded_app_count"`
	SupersedingAppCount        types.Int32    `tfsdk:"superseding_app_count"`
	Timeouts                   timeouts.Value `tfsdk:"timeouts"`
}
