package utilityGuidListSharder

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GuidListSharderResourceModel struct {
	Id                     types.String   `tfsdk:"id"`
	RecalculateOnNextRun   types.Bool     `tfsdk:"recalculate_on_next_run"`
	ResourceType           types.String   `tfsdk:"resource_type"`
	GroupId                types.String   `tfsdk:"group_id"`
	ODataFilter            types.String   `tfsdk:"odata_filter"`
	ShardCount             types.Int64    `tfsdk:"shard_count"`
	ShardPercentages       types.List     `tfsdk:"shard_percentages"`
	ShardSizes             types.List     `tfsdk:"shard_sizes"`
	Strategy               types.String   `tfsdk:"strategy"`
	Seed                   types.String   `tfsdk:"seed"`
	Shards                 types.Map      `tfsdk:"shards"`
	Timeouts               timeouts.Value `tfsdk:"timeouts"`
}
