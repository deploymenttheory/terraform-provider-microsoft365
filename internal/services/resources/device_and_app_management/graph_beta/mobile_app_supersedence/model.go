// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappsupersedence?view=graph-rest-beta

package graphBetaMobileAppSupersedence

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileAppSupersedenceResourceModel represents the root Terraform resource model for Mobile App Supersedence
type MobileAppSupersedenceResourceModel struct {
	ID                         types.String   `tfsdk:"id"`
	TargetID                   types.String   `tfsdk:"target_id"`
	TargetDisplayName          types.String   `tfsdk:"target_display_name"`
	TargetDisplayVersion       types.String   `tfsdk:"target_display_version"`
	TargetPublisher            types.String   `tfsdk:"target_publisher"`
	TargetPublisherDisplayName types.String   `tfsdk:"target_publisher_display_name"`
	SourceID                   types.String   `tfsdk:"source_id"`
	SourceDisplayName          types.String   `tfsdk:"source_display_name"`
	SourceDisplayVersion       types.String   `tfsdk:"source_display_version"`
	SourcePublisherDisplayName types.String   `tfsdk:"source_publisher_display_name"`
	TargetType                 types.String   `tfsdk:"target_type"`
	AppRelationshipType        types.String   `tfsdk:"app_relationship_type"`
	SupersedenceType           types.String   `tfsdk:"supersedence_type"`
	SupersededAppCount         types.Int32    `tfsdk:"superseded_app_count"`
	SupersedingAppCount        types.Int32    `tfsdk:"superseding_app_count"`
	Timeouts                   timeouts.Value `tfsdk:"timeouts"`
}
