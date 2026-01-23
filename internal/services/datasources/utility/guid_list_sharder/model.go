package utilityGuidListSharder

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GuidListSharderDataSourceModel struct {
	Id               types.String   `tfsdk:"id"`
	ResourceType     types.String   `tfsdk:"resource_type"`
	GroupId          types.String   `tfsdk:"group_id"`
	ODataQuery       types.String   `tfsdk:"odata_query"`
	ShardCount       types.Int64    `tfsdk:"shard_count"`
	ShardPercentages types.List     `tfsdk:"shard_percentages"`
	Strategy         types.String   `tfsdk:"strategy"`
	Seed             types.String   `tfsdk:"seed"`
	Shards           types.Map      `tfsdk:"shards"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}
