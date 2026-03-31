package graphBetaWindowsUpdatesAutopatchUpdatePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var GradualRolloutAttrTypes = map[string]attr.Type{
	"duration_between_offers": types.StringType,
	"devices_per_offer":       types.Int32Type,
}

var ScheduleSettingsAttrTypes = map[string]attr.Type{
	"start_date_time": types.StringType,
	"gradual_rollout": types.ObjectType{AttrTypes: GradualRolloutAttrTypes},
}

var DeploymentSettingsAttrTypes = map[string]attr.Type{
	"schedule": types.ObjectType{AttrTypes: ScheduleSettingsAttrTypes},
}
