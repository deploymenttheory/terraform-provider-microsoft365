// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsupdateforbusinessconfiguration?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-windowsupdateforbusinessconfiguration-extendfeatureupdatespause?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-windowsupdateforbusinessconfiguration-extendqualityupdatespause?view=graph-rest-beta
package graphBetaWindowsUpdateRingAction

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsUpdateRingActionResourceModel defines the model for Windows Update Ring Action resource
type WindowsUpdateRingActionResourceModel struct {
	ID           types.String `tfsdk:"id"`
	UpdateRingId types.String `tfsdk:"update_ring_id"` // Required: The ID of the Windows Update Ring

	// Feature Update Actions
	PauseFeatureUpdates       types.Bool `tfsdk:"pause_feature_updates"`        // Optional: Set to true to pause
	ResumeFeatureUpdates      types.Bool `tfsdk:"resume_feature_updates"`       // Optional: Set to true to resume
	ExtendFeatureUpdatesPause types.Bool `tfsdk:"extend_feature_updates_pause"` // Optional: Set to true to extend pause (max 35 days)
	RollbackFeatureUpdates    types.Bool `tfsdk:"rollback_feature_updates"`     // Optional: Set to true to rollback

	// Quality Update Actions
	PauseQualityUpdates    types.Bool `tfsdk:"pause_quality_updates"`    // Optional: Set to true to pause
	ResumeQualityUpdates   types.Bool `tfsdk:"resume_quality_updates"`   // Optional: Set to true to resume
	RollbackQualityUpdates types.Bool `tfsdk:"rollback_quality_updates"` // Optional: Set to true to rollback

	// Metadata
	Description         types.String   `tfsdk:"description"`           // Optional: Description of actions
	LastActionPerformed types.String   `tfsdk:"last_action_performed"` // Computed: Last action that was executed
	LastActionTimestamp types.String   `tfsdk:"last_action_timestamp"` // Computed: When last action occurred
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}
