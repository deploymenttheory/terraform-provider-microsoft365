package graphBetaMobileAppDependency

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MobileAppDependencyResourceModel struct {
	ID                         types.String   `tfsdk:"id"`
	SourceID                   types.String   `tfsdk:"source_id"`
	TargetID                   types.String   `tfsdk:"target_id"`
	DependencyType             types.String   `tfsdk:"dependency_type"`
	TargetType                 types.String   `tfsdk:"target_type"`
	SourceDisplayName          types.String   `tfsdk:"source_display_name"`
	SourceDisplayVersion       types.String   `tfsdk:"source_display_version"`
	SourcePublisherDisplayName types.String   `tfsdk:"source_publisher_display_name"`
	TargetDisplayName          types.String   `tfsdk:"target_display_name"`
	TargetDisplayVersion       types.String   `tfsdk:"target_display_version"`
	TargetPublisher            types.String   `tfsdk:"target_publisher"`
	TargetPublisherDisplayName types.String   `tfsdk:"target_publisher_display_name"`
	DependentAppCount          types.Int32    `tfsdk:"dependent_app_count"`
	DependsOnAppCount          types.Int32    `tfsdk:"depends_on_app_count"`
	Timeouts                   timeouts.Value `tfsdk:"timeouts"`
}
